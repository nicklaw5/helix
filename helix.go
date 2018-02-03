package helix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	basePath = "https://api.twitch.tv/helix"
)

// HTTPClient ...
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	clientID    string
	accessToken string
	userAgent   string
	httpClient  HTTPClient
}

// ResponseCommon ...
type ResponseCommon struct {
	Error              string `json:"error"`
	ErrorStatus        int    `json:"status"`
	ErrorMessage       string `json:"message"`
	RatelimitLimit     int
	RatelimitRemaining int
	RatelimitReset     int
	StatusCode         int
}

// NewClient ... It is concurrecy safe.
func NewClient(clientID string, httpClient HTTPClient) (*Client, error) {
	c := &Client{}
	c.clientID = clientID

	if c.clientID == "" {
		return nil, errors.New("clientID cannot be an empty string")
	}

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}

	return c, nil
}

// Get ...
func (c *Client) Get(path string, v interface{}) error {
	req, err := newRequest("GET", path)
	if err != nil {
		return err
	}

	err = c.doRequest(req, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	c.setRequestHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute API request: %s", err.Error())
	}
	defer resp.Body.Close()

	setResponseStatusCode(v, "StatusCode", resp.StatusCode)
	setRatelimitValue(v, "RatelimitLimit", resp.Header.Get("Ratelimit-Limit"))
	setRatelimitValue(v, "RatelimitRemaining", resp.Header.Get("Ratelimit-Remaining"))
	setRatelimitValue(v, "RatelimitReset", resp.Header.Get("Ratelimit-Reset"))

	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return fmt.Errorf("Failed to decode API response: %s", err.Error())
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

func setRatelimitValue(v interface{}, fieldName, value string) {
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

func newRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, basePath+path, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = req.URL.Query().Encode()

	return req, nil
}

// concatString concatenates each of the strings provided by
// strs in the order they are presented. You may also pass in
// an optional delimiter to be appended along with the strings.
func concatString(strs []string, delimiter ...string) string {
	var buffer bytes.Buffer
	appendDelimiter := len(delimiter) > 0

	for _, str := range strs {
		s := str
		if appendDelimiter {
			s = concatString([]string{s, delimiter[0]})
		}
		buffer.Write([]byte(s))
	}

	if appendDelimiter {
		return strings.TrimSuffix(buffer.String(), delimiter[0])
	}

	return buffer.String()
}
