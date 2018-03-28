package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	methodGet  = "GET"
	methodPost = "POST"
	queryTag   = "query"

	// APIBaseURL is the base URL for composing API requests.
	APIBaseURL = "https://api.twitch.tv/helix"

	// AuthBaseURL is the base URL for composing authentication requests.
	AuthBaseURL = "https://id.twitch.tv/oauth2"
)

// HTTPClient ...
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	clientID      string
	clientSecret  string
	accessToken   string
	userAgent     string
	redirectURI   string
	scopes        []string
	httpClient    HTTPClient
	rateLimitFunc RateLimitFunc

	baseURL      string
	lastResponse *Response
}

// Options ...
type Options struct {
	ClientID      string
	ClientSecret  string
	AccessToken   string
	UserAgent     string
	RedirectURI   string
	Scopes        []string
	HTTPClient    HTTPClient
	RateLimitFunc RateLimitFunc
}

// RateLimitFunc ...
type RateLimitFunc func(*Response) error

// ResponseCommon ...
type ResponseCommon struct {
	StatusCode   int
	Error        string `json:"error"`
	ErrorStatus  int    `json:"status"`
	ErrorMessage string `json:"message"`
	RateLimit
	StreamsMetadataRateLimit
}

// Response ...
type Response struct {
	ResponseCommon
	Data interface{}
}

// RateLimit ...
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
}

// Pagination ...
type Pagination struct {
	Cursor string `json:"cursor"`
}

// NewClient returns a new Twicth Helix API client. It panics if
// clientID is an empty string. It is concurrecy safe.
func NewClient(options *Options) (*Client, error) {
	if options.ClientID == "" {
		return nil, errors.New("A client ID was not provided but is required")
	}

	c := &Client{
		clientID:   options.ClientID,
		httpClient: http.DefaultClient,
	}

	// Set options
	if options.HTTPClient != nil {
		c.httpClient = options.HTTPClient
	}
	c.clientSecret = options.ClientSecret
	c.accessToken = options.AccessToken
	c.userAgent = options.UserAgent
	c.rateLimitFunc = options.RateLimitFunc
	c.scopes = options.Scopes
	c.redirectURI = options.RedirectURI

	// Set non-options
	c.baseURL = APIBaseURL

	return c, nil
}

func (c *Client) get(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(methodGet, path, respData, reqData)
}

func (c *Client) post(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(methodPost, path, respData, reqData)
}

func (c *Client) sendRequest(method, path string, respData, reqData interface{}) (*Response, error) {
	resp := &Response{}
	if respData != nil {
		resp.Data = respData
	}

	req, err := c.newRequest(method, path, reqData)
	if err != nil {
		return nil, err
	}

	err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func buildQueryString(req *http.Request, v interface{}) (string, error) {
	isNil, err := isZero(v)
	if err != nil {
		return "", err
	}

	if isNil {
		return "", nil
	}

	query := req.URL.Query()
	t := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		var defaultValue string

		field := t.Field(i)
		tag := field.Tag.Get(queryTag)

		// Get the default value from the struct tag
		if strings.Contains(tag, ",") {
			tagSlice := strings.Split(tag, ",")

			tag = tagSlice[0]
			defaultValue = tagSlice[1]
		}

		// Get the value assigned to the query param
		if field.Type.Kind() == reflect.Slice {
			fieldVal := val.Field(i)
			for j := 0; j < fieldVal.Len(); j++ {
				query.Add(tag, fmt.Sprintf("%v", fieldVal.Index(j)))
			}
		} else {
			value := fmt.Sprintf("%v", val.Field(i))

			// If no value was set by the user, use the default
			// value specified in the struct tag.
			if value == "" || value == "0" {
				if defaultValue == "" {
					continue
				}

				value = defaultValue
			}

			query.Add(tag, value)
		}
	}

	return query.Encode(), nil
}

func isZero(v interface{}) (bool, error) {
	t := reflect.TypeOf(v)
	if !t.Comparable() {
		return false, fmt.Errorf("type is not comparable: %v", t)
	}
	return v == reflect.Zero(t).Interface(), nil
}

func (c *Client) newRequest(method, path string, data interface{}) (*http.Request, error) {
	url := c.getBaseURL(path) + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if data != nil {
		query, err := buildQueryString(req, data)
		if err != nil {
			return nil, err
		}

		req.URL.RawQuery = query
	}

	return req, nil
}

func (c *Client) getBaseURL(path string) string {
	for _, authPath := range authPaths {
		if strings.Contains(path, authPath) {
			return AuthBaseURL
		}
	}

	return APIBaseURL
}

func (c *Client) doRequest(req *http.Request, resp *Response) error {
	c.setRequestHeaders(req)

	for {
		if c.lastResponse != nil && c.rateLimitFunc != nil {
			err := c.rateLimitFunc(c.lastResponse)
			if err != nil {
				return err
			}
		}

		response, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("Failed to execute API request: %s", err.Error())
		}
		defer response.Body.Close()

		setResponseStatusCode(resp, "StatusCode", response.StatusCode)
		setRateLimitValue(&resp.RateLimit, "Limit", response.Header.Get("RateLimit-Limit"))
		setRateLimitValue(&resp.RateLimit, "Remaining", response.Header.Get("RateLimit-Remaining"))
		setRateLimitValue(&resp.RateLimit, "Reset", response.Header.Get("RateLimit-Reset"))

		if strings.Contains(req.URL.Path, streamsMetadataPath) {
			setRateLimitValue(&resp.StreamsMetadataRateLimit, "Limit", response.Header.Get("Ratelimit-Helixstreamsmetadata-Limit"))
			setRateLimitValue(&resp.StreamsMetadataRateLimit, "Remaining", response.Header.Get("Ratelimit-Helixstreamsmetadata-Remaining"))
		}

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		// Only attempt to decode the response if we have a response we can handle
		if len(bodyBytes) > 0 && resp.StatusCode < http.StatusInternalServerError {
			if resp.Data != nil && resp.StatusCode < http.StatusBadRequest {
				// Successful request
				err = json.Unmarshal(bodyBytes, &resp.Data)
			} else {
				// Failed request
				err = json.Unmarshal(bodyBytes, &resp)
			}

			if err != nil {
				return fmt.Errorf("Failed to decode API response: %s", err.Error())
			}
		}

		if c.rateLimitFunc == nil {
			break
		} else {
			c.lastResponse = resp

			if c.rateLimitFunc != nil &&
				c.lastResponse.StatusCode == http.StatusTooManyRequests {
				// Rate limit exceeded,  retry to send request after
				// applying rate limiter callback
				continue
			}

			break
		}
	}

	return nil
}

func (c *Client) setRequestHeaders(req *http.Request) {
	req.Header.Set("Client-ID", c.clientID)
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
}

func setResponseStatusCode(v interface{}, fieldName string, code int) {
	s := reflect.ValueOf(v).Elem()
	field := s.FieldByName(fieldName)
	field.SetInt(int64(code))
}

func setRateLimitValue(v interface{}, fieldName, value string) {
	s := reflect.ValueOf(v).Elem()
	field := s.FieldByName(fieldName)
	intVal, _ := strconv.Atoi(value)
	field.SetInt(int64(intVal))
}

// SetAccessToken ...
func (c *Client) SetAccessToken(AccessToken string) {
	c.accessToken = AccessToken
}

// SetUserAgent ...
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// SetScopes ...
func (c *Client) SetScopes(scopes []string) {
	c.scopes = scopes
}

// SetRedirectURI ...
func (c *Client) SetRedirectURI(uri string) {
	c.redirectURI = uri
}
