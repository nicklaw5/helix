package helix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultAPIBaseURL is the base URL for composing API requests.
	DefaultAPIBaseURL = "https://api.twitch.tv/helix"

	// AuthBaseURL is the base URL for composing authentication requests.
	AuthBaseURL = "https://id.twitch.tv/oauth2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Logger is the interface used by the helix client for debug logging.
// The stdlib *log.Logger satisfies this interface.
type Logger interface {
	Printf(format string, v ...interface{})
}

type Client struct {
	mu           sync.RWMutex
	ctx          context.Context
	opts         *Options
	lastResponse *Response
	callbacks    struct {
		onUserAccessTokenRefreshed func(newAccessToken, newRefreshToken string)
		onAppAccessTokenRefreshed  func(newAccessToken string)
	}
}

type Options struct {
	ClientID          string
	ClientSecret      string
	AppAccessToken    string
	AppAccessScopes   []string
	DeviceAccessToken string
	UserAccessToken   string
	RefreshToken      string
	UserAgent         string
	RedirectURI       string
	HTTPClient        HTTPClient
	RateLimitFunc     RateLimitFunc
	APIBaseURL        string
	ExtensionOpts     ExtensionOptions
	// DebugMode enables debug logging of outgoing requests and incoming responses.
	// WARNING: debug logs may contain sensitive data such as tokens and API responses.
	// Only enable in non-production environments.
	DebugMode bool
	// Logger is the logger used when DebugMode is true. If nil, a default logger
	// writing to os.Stderr is used. Any type implementing Printf(string, ...interface{})
	// is accepted (e.g. *log.Logger).
	Logger Logger
}

type ExtensionOptions struct {
	OwnerUserID    string
	Secret         string
	SignedJWTToken string
}

// DateRange is a generic struct used by various responses.
type DateRange struct {
	StartedAt Time `json:"started_at"`
	EndedAt   Time `json:"ended_at"`
}

type RateLimitFunc func(*Response) error

// DefaultRateLimitFunc is a default rate limit function that sleeps until the
// rate limit reset time if the rate limit has been reached (i.e. no remaining
// requests). It can be used as the RateLimitFunc in Options.
func DefaultRateLimitFunc(lastResponse *Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	reset64 := int64(lastResponse.GetRateLimitReset())
	currentTime := time.Now().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64-currentTime) * time.Second
		time.Sleep(timeDiff)
	}

	return nil
}

type ResponseCommon struct {
	StatusCode   int
	Header       http.Header
	Error        string `json:"error"`
	ErrorStatus  int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func (rc *ResponseCommon) convertHeaderToInt(str string) int {
	i, _ := strconv.Atoi(str)

	return i
}

// GetRateLimit returns the "RateLimit-Limit" header as an int.
func (rc *ResponseCommon) GetRateLimit() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Limit"))
}

// GetRateLimitRemaining returns the "RateLimit-Remaining" header as an int.
func (rc *ResponseCommon) GetRateLimitRemaining() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Remaining"))
}

// GetRateLimitReset returns the "RateLimit-Reset" header as an int.
func (rc *ResponseCommon) GetRateLimitReset() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Reset"))
}

type Response struct {
	ResponseCommon
	Data interface{}
}

// HydrateResponseCommon copies the content of the source response's ResponseCommon to the supplied ResponseCommon argument
func (r *Response) HydrateResponseCommon(rc *ResponseCommon) {
	rc.StatusCode = r.ResponseCommon.StatusCode
	rc.Header = r.ResponseCommon.Header
	rc.Error = r.ResponseCommon.Error
	rc.ErrorStatus = r.ResponseCommon.ErrorStatus
	rc.ErrorMessage = r.ResponseCommon.ErrorMessage
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

// NewClient returns a new Twitch Helix API client. It returns an
// if clientID is an empty string. It is concurrency safe.
func NewClient(options *Options) (*Client, error) {
	return NewClientWithContext(context.Background(), options)
}

func NewClientWithContext(ctx context.Context, options *Options) (*Client, error) {
	if options.ClientID == "" {
		return nil, errors.New("A client ID was not provided but is required")
	}

	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}

	if options.APIBaseURL == "" {
		options.APIBaseURL = DefaultAPIBaseURL
	}

	if options.DebugMode && options.Logger == nil {
		options.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	client := &Client{
		ctx:  ctx,
		opts: options,
	}

	return client, nil
}

func (c *Client) logf(format string, v ...interface{}) {
	if c.opts.DebugMode && c.opts.Logger != nil {
		c.opts.Logger.Printf(format, v...)
	}
}

func (c *Client) get(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodGet, path, respData, reqData, false)
}

func (c *Client) post(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodPost, path, respData, reqData, false)
}

func (c *Client) put(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodPut, path, respData, reqData, false)
}

func (c *Client) delete(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodDelete, path, respData, reqData, false)
}

func (c *Client) patchAsJSON(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodPatch, path, respData, reqData, true)
}

func (c *Client) postAsJSON(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodPost, path, respData, reqData, true)
}

func (c *Client) putAsJSON(path string, respData, reqData interface{}) (*Response, error) {
	return c.sendRequest(http.MethodPut, path, respData, reqData, true)
}

func (c *Client) sendRequest(method, path string, respData, reqData interface{}, hasJSONBody bool) (*Response, error) {
	resp := &Response{}
	if respData != nil {
		resp.Data = respData
	}

	req, err := c.newRequest(method, path, reqData, hasJSONBody)
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
	vType := reflect.TypeOf(v).Elem()
	vValue := reflect.ValueOf(v).Elem()

	for i := 0; i < vType.NumField(); i++ {
		var defaultValue string

		field := vType.Field(i)
		tag := field.Tag.Get("query")

		if tag == "" {
			continue
		}

		// Get the default value from the struct tag
		if strings.Contains(tag, ",") {
			tagSlice := strings.Split(tag, ",")

			tag = tagSlice[0]
			defaultValue = tagSlice[1]

			if defaultValue == "omitempty" {
				defaultValue = ""
			}
		}

		if field.Type.Kind() == reflect.Slice {
			// Attach any slices as query params
			fieldVal := vValue.Field(i)
			for j := 0; j < fieldVal.Len(); j++ {
				query.Add(tag, fmt.Sprintf("%v", fieldVal.Index(j)))
			}
		} else if isDatetimeTagField(tag) {
			// Get and correctly format datetime fields, and attach them query params
			dateStr := fmt.Sprintf("%v", vValue.Field(i))

			if strings.Contains(dateStr, " m=") {
				datetimeSplit := strings.Split(dateStr, " m=")
				dateStr = datetimeSplit[0]
			}

			date, err := time.Parse(requestDateTimeFormat, dateStr)
			if err != nil {
				return "", err
			}

			// Determine if the date has been set. If it has we'll add it to the query.
			if !date.IsZero() {
				query.Add(tag, date.Format(time.RFC3339))
			}
		} else {
			// Add any scalar values as query params
			fieldVal := fmt.Sprintf("%v", vValue.Field(i))

			// If no value was set by the user, use the default
			// value specified in the struct tag.
			if fieldVal == "" || fieldVal == "0" {
				if defaultValue == "" {
					continue
				}

				fieldVal = defaultValue
			}

			query.Add(tag, fieldVal)
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

func (c *Client) newRequest(method, path string, data interface{}, hasJSONBody bool) (*http.Request, error) {
	url := c.getBaseURL(path) + path

	if hasJSONBody {
		return c.newJSONRequest(method, url, data)
	}

	return c.newStandardRequest(method, url, data)
}

func (c *Client) newStandardRequest(method, url string, data interface{}) (*http.Request, error) {
	req, err := http.NewRequestWithContext(c.ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return req, nil
	}

	query, err := buildQueryString(req, data)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query

	return req, nil
}

func (c *Client) newJSONRequest(method, url string, data interface{}) (*http.Request, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)

	req, err := http.NewRequestWithContext(c.ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	query, err := buildQueryString(req, data)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) getBaseURL(path string) string {
	for _, authPath := range authPaths {
		if strings.Contains(path, authPath) {
			return AuthBaseURL
		}
	}

	return c.opts.APIBaseURL
}

func (c *Client) doRequest(req *http.Request, resp *Response) error {
	c.setRequestHeaders(req)

	rateLimitFunc := c.opts.RateLimitFunc
	tokenRefreshed := false

	for {
		if c.lastResponse != nil && rateLimitFunc != nil {
			err := rateLimitFunc(c.lastResponse)
			if err != nil {
				return err
			}
		}

		c.logf("helix: request %s %s", req.Method, req.URL.String())

		response, err := c.opts.HTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("Failed to execute API request: %s", err.Error())
		}
		defer response.Body.Close()

		resp.Header = response.Header

		setResponseStatusCode(resp, "StatusCode", response.StatusCode)

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		c.logf("helix: response %d %s", response.StatusCode, string(bodyBytes))

		contentType := response.Header.Get("Content-Type")
		mt, _, _ := mime.ParseMediaType(contentType)
		isHTML := strings.EqualFold(mt, "text/html")

		// Only attempt to decode the response if we have a response we can handle
		if len(bodyBytes) > 0 && resp.StatusCode < http.StatusInternalServerError && !isHTML {
			if resp.Data != nil && resp.StatusCode < http.StatusBadRequest {
				// Successful request
				err = json.Unmarshal(bodyBytes, &resp.Data)
			} else {
				// Failed request
				err = json.Unmarshal(bodyBytes, &resp)
				if err != nil {
					return fmt.Errorf("Failed to decode API response: %s", err.Error())
				}

				// A 401 may mean Twitch wants us to refresh our token:
				// https://dev.twitch.tv/docs/authentication/refresh-tokens/
				// However, if the error is about missing scopes, refreshing
				// won't help and would cause an infinite loop.
				if resp.StatusCode == http.StatusUnauthorized && !tokenRefreshed &&
					!strings.HasPrefix(resp.ErrorMessage, "Missing scope") {
					refreshed, refreshErr := c.tryRefreshToken()
					if refreshErr != nil {
						log.Printf("Failed to refresh helix token: %v", refreshErr)
						break
					}
					if refreshed {
						tokenRefreshed = true
						// Try again now that we have a new token
						c.setRequestHeaders(req)
						continue
					}
				}
			}

			if err != nil {
				return fmt.Errorf("Failed to decode API response: %s", err.Error())
			}
		}

		if rateLimitFunc == nil {
			break
		} else {
			c.mu.Lock()
			c.lastResponse = resp
			c.mu.Unlock()

			if rateLimitFunc != nil &&
				c.lastResponse.StatusCode == http.StatusTooManyRequests {
				// Rate limit exceeded, retry to send request after
				// applying rate limiter callback
				continue
			}

			break
		}
	}

	return nil
}

func (c *Client) canRefreshToken() bool {
	return ((c.opts.UserAccessToken != "" && c.opts.ClientSecret != "") ||
		c.opts.DeviceAccessToken != "") &&
		(c.opts.ClientID != "" && c.opts.RefreshToken != "")
}

// tryRefreshToken attempts to refresh whichever token type is currently in use.
// It returns true if a refresh was successfully performed, false if no refresh
// was applicable, and an error if a refresh was attempted but failed.
func (c *Client) tryRefreshToken() (bool, error) {
	if c.canRefreshToken() {
		return true, c.refreshToken()
	}
	if c.canRefreshAppToken() {
		return true, c.refreshAppToken()
	}
	return false, nil
}

func (c *Client) refreshToken() error {
	resp, err := c.RefreshUserAccessToken(c.opts.RefreshToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		statusCode := -1
		var errorMessage string
		if resp != nil {
			statusCode = resp.StatusCode
			errorMessage = resp.ErrorMessage
		}
		return fmt.Errorf("failed to refresh token: (%d: %s) %v", statusCode, errorMessage, err)
	}

	c.mu.Lock()
	c.opts.UserAccessToken = resp.Data.AccessToken
	c.opts.RefreshToken = resp.Data.RefreshToken
	c.mu.Unlock()

	if cb := c.callbacks.onUserAccessTokenRefreshed; cb != nil {
		go cb(resp.Data.AccessToken, resp.Data.RefreshToken)
	}

	return nil
}

func (c *Client) canRefreshAppToken() bool {
	return c.opts.AppAccessToken != "" && c.opts.ClientID != "" && c.opts.ClientSecret != ""
}

func (c *Client) refreshAppToken() error {
	resp, err := c.RequestAppAccessToken(c.opts.AppAccessScopes)
	if err != nil || resp.StatusCode != http.StatusOK {
		statusCode := -1
		var errorMessage string
		if resp != nil {
			statusCode = resp.StatusCode
			errorMessage = resp.ErrorMessage
		}
		return fmt.Errorf("failed to refresh app token: (%d: %s) %v", statusCode, errorMessage, err)
	}

	c.mu.Lock()
	c.opts.AppAccessToken = resp.Data.AccessToken
	c.mu.Unlock()

	if cb := c.callbacks.onAppAccessTokenRefreshed; cb != nil {
		go cb(resp.Data.AccessToken)
	}

	return nil
}

func (c *Client) setRequestHeaders(req *http.Request) {
	opts := c.opts

	req.Header.Set("Client-ID", opts.ClientID)

	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	var bearerToken string
	if opts.AppAccessToken != "" {
		bearerToken = opts.AppAccessToken
	}
	if opts.DeviceAccessToken != "" {
		bearerToken = opts.DeviceAccessToken
	}
	if opts.UserAccessToken != "" {
		bearerToken = opts.UserAccessToken
	}
	if opts.ExtensionOpts.SignedJWTToken != "" {
		bearerToken = opts.ExtensionOpts.SignedJWTToken
	}

	authType := "Bearer"
	// Token validation requires different type of Auth
	if req.URL.String() == AuthBaseURL+authPaths["validate"] {
		authType = "OAuth"
	}

	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, bearerToken))
	}
}

func setResponseStatusCode(v interface{}, fieldName string, code int) {
	s := reflect.ValueOf(v).Elem()
	field := s.FieldByName(fieldName)
	field.SetInt(int64(code))
}

// GetAppAccessToken returns the current app access token.
func (c *Client) GetAppAccessToken() string {
	return c.opts.AppAccessToken
}

func (c *Client) SetAppAccessToken(accessToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.AppAccessToken = accessToken
}

// GetDeviceAccessToken returns the current device access token.
func (c *Client) GetDeviceAccessToken() string {
	return c.opts.DeviceAccessToken
}

// SetDeviceAccessToken sets the current device access token.
func (c *Client) SetDeviceAccessToken(accessToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.DeviceAccessToken = accessToken
}

// GetUserAccessToken returns the current user access token.
func (c *Client) GetUserAccessToken() string {
	return c.opts.UserAccessToken
}

func (c *Client) SetUserAccessToken(accessToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.UserAccessToken = accessToken
}

// GetRefreshToken returns the current refresh token.
func (c *Client) GetRefreshToken() string {
	return c.opts.RefreshToken
}

func (c *Client) SetRefreshToken(refreshToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.RefreshToken = refreshToken
}

// GetAppAccessToken returns the current app access token.
func (c *Client) GetExtensionSignedJWTToken() string {
	return c.opts.ExtensionOpts.SignedJWTToken
}

func (c *Client) SetExtensionSignedJWTToken(jwt string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.ExtensionOpts.SignedJWTToken = jwt
}

func (c *Client) SetUserAgent(userAgent string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.UserAgent = userAgent
}

func (c *Client) SetRedirectURI(uri string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts.RedirectURI = uri
}

func (c *Client) OnUserAccessTokenRefreshed(f func(newAccessToken, newRefreshToken string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks.onUserAccessTokenRefreshed = f
}

func (c *Client) OnAppAccessTokenRefreshed(f func(newAccessToken string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks.onAppAccessTokenRefreshed = f
}
