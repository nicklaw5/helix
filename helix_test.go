package helix

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type mockHTTPClient struct {
	mockHandler http.HandlerFunc
}

func (mtc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mtc.mockHandler)
	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func newMockClient(options *Options, mockHandler http.HandlerFunc) *Client {
	options.HTTPClient = &mockHTTPClient{mockHandler}
	return &Client{opts: options}
}

func newMockHandler(statusCode int, json string, headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if headers != nil && len(headers) > 0 {
			for key, value := range headers {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(json))
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		extpectErr bool
		options    *Options
	}{
		{
			true,
			&Options{}, // no client id
		},
		{
			false,
			&Options{
				ClientID:        "my-client-id",
				ClientSecret:    "my-client-secret",
				HTTPClient:      &http.Client{},
				AppAccessToken:  "my-app-access-token",
				UserAccessToken: "my-user-access-token",
				UserAgent:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
				RateLimitFunc:   func(*Response) error { return nil },
				RedirectURI:     "http://localhost/auth/callback",
				APIBaseURL:      "http://localhost/proxy",
			},
		},
	}

	for _, testCase := range testCases {
		client, err := NewClient(testCase.options)
		if err != nil && !testCase.extpectErr {
			t.Errorf("Did not expect an error, got \"%s\"", err.Error())
		}

		if testCase.extpectErr {
			continue
		}

		opts := client.opts

		if opts.ClientID != testCase.options.ClientID {
			t.Errorf("expected clientID to be \"%s\", got \"%s\"", testCase.options.ClientID, opts.ClientID)
		}

		if opts.ClientSecret != testCase.options.ClientSecret {
			t.Errorf("expected clientSecret to be \"%s\", got \"%s\"", testCase.options.ClientSecret, opts.ClientSecret)
		}

		if reflect.TypeOf(opts.RateLimitFunc).Kind() != reflect.Func {
			t.Errorf("expected rateLimitFunc to be a function, got \"%+v\"", reflect.TypeOf(opts.RateLimitFunc).Kind())
		}

		if opts.HTTPClient != testCase.options.HTTPClient {
			t.Errorf("expected httpClient to be \"%s\", got \"%s\"", testCase.options.HTTPClient, opts.HTTPClient)
		}

		if opts.UserAgent != testCase.options.UserAgent {
			t.Errorf("expected userAgent to be \"%s\", got \"%s\"", testCase.options.UserAgent, opts.UserAgent)
		}

		if opts.AppAccessToken != testCase.options.AppAccessToken {
			t.Errorf("expected accessToken to be \"%s\", got \"%s\"", testCase.options.AppAccessToken, opts.AppAccessToken)
		}

		if opts.UserAccessToken != testCase.options.UserAccessToken {
			t.Errorf("expected accessToken to be \"%s\", got \"%s\"", testCase.options.UserAccessToken, opts.UserAccessToken)
		}

		if opts.RedirectURI != testCase.options.RedirectURI {
			t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", testCase.options.RedirectURI, opts.RedirectURI)
		}

		if opts.APIBaseURL != testCase.options.APIBaseURL {
			t.Errorf("expected APIBaseURL to be \"%s\", got \"%s\"", testCase.options.APIBaseURL, opts.APIBaseURL)
		}
	}
}

func TestNewClientDefault(t *testing.T) {
	t.Parallel()

	options := &Options{ClientID: "my-client-id"}

	client, err := NewClient(options)
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	opts := client.opts

	if opts.ClientID != options.ClientID {
		t.Errorf("expected clientID to be \"%s\", got \"%s\"", options.ClientID, opts.ClientID)
	}

	if opts.ClientSecret != "" {
		t.Errorf("expected clientSecret to be \"%s\", got \"%s\"", options.ClientSecret, opts.ClientSecret)
	}

	if opts.UserAgent != "" {
		t.Errorf("expected userAgent to be \"%s\", got \"%s\"", "", opts.UserAgent)
	}

	if opts.UserAccessToken != "" {
		t.Errorf("expected accesstoken to be \"\", got \"%s\"", opts.UserAccessToken)
	}

	if opts.HTTPClient != http.DefaultClient {
		t.Errorf("expected httpClient to be \"%v\", got \"%v\"", http.DefaultClient, opts.HTTPClient)
	}

	if opts.RateLimitFunc != nil {
		t.Errorf("expected rateLimitFunc to be \"%v\", got \"%v\"", nil, opts.RateLimitFunc)
	}

	if opts.RedirectURI != options.RedirectURI {
		t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", options.RedirectURI, opts.RedirectURI)
	}
}

func TestRatelimitCallback(t *testing.T) {
	t.Parallel()

	respBody1 := `{"error":"Too Many Requests","status":429,"message":"Request limit exceeded"}`
	options1 := &Options{
		ClientID: "my-client-id",
		RateLimitFunc: func(resp *Response) error {
			return nil
		},
	}

	c := newMockClient(options1, newMockHandler(http.StatusTooManyRequests, respBody1, nil))
	go func() {
		_, err := c.GetStreams(&StreamsParams{})
		if err != nil {
			t.Errorf("Did not expect error, got \"%s\"", err.Error())
		}
	}()

	time.Sleep(5 * time.Millisecond)

	respBody2 := `{"data":[{"id":"EncouragingPluckySlothSSSsss","url":"https://clips.twitch.tv/EncouragingPluckySlothSSSsss","embed_url":"https://clips.twitch.tv/embed?clip=EncouragingPluckySlothSSSsss","broadcaster_id":"26490481","creator_id":"143839181","video_id":"222004532","game_id":"490377","language":"en","title":"summit and fat tim discover how to use maps","view_count":81808,"created_at":"2018-01-25T04:04:15Z","thumbnail_url":"https://clips-media-assets.twitch.tv/182509178-preview-480x272.jpg"}]}`
	options2 := &Options{
		ClientID: "my-client-id",
	}

	cOK := newMockClient(options2, newMockHandler(http.StatusOK, respBody2, nil))
	_, err := cOK.GetStreams(&StreamsParams{})
	if err != nil {
		t.Errorf("Did not expect error, got \"%s\"", err.Error())
	}
}

func TestRatelimitCallbackFailsOnError(t *testing.T) {
	t.Parallel()

	errMsg := "Oops! Your rate limiter funciton is broken :("
	respBody1 := `{"error":"Too Many Requests","status":429,"message":"Request limit exceeded"}`
	options1 := &Options{
		ClientID: "my-client-id",
		RateLimitFunc: func(resp *Response) error {
			return errors.New(errMsg)
		},
	}

	c := newMockClient(options1, newMockHandler(http.StatusTooManyRequests, respBody1, nil))
	_, err := c.GetStreams(&StreamsParams{})
	if err == nil {
		t.Error("Expected error, got none")
	}

	if err.Error() != errMsg {
		t.Errorf("Expected error to be \"%s\", got \"%s\"", errMsg, err.Error())
	}
}

func TestSetRequestHeaders(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		endpoint        string
		method          string
		userBearerToken string
		appBearerToken  string
	}{
		{"/users", "GET", "my-user-access-token", "my-app-access-token"},
		{"/entitlements/upload", "POST", "", "my-app-access-token"},
		{"/streams", "GET", "", ""},
	}

	for _, testCase := range testCases {
		client, err := NewClient(&Options{
			ClientID:        "my-client-id",
			UserAgent:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
			AppAccessToken:  testCase.appBearerToken,
			UserAccessToken: testCase.userBearerToken,
		})
		if err != nil {
			t.Errorf("Did not expect an error, got \"%s\"", err.Error())
		}

		req, _ := http.NewRequest(testCase.method, testCase.endpoint, nil)
		client.setRequestHeaders(req)

		if testCase.userBearerToken != "" {
			expectedAuthHeader := "Bearer " + testCase.userBearerToken
			if req.Header.Get("Authorization") != expectedAuthHeader {
				t.Errorf("expected Authorization header to be \"%s\", got \"%s\"", expectedAuthHeader, req.Header.Get("Authorization"))
			}
		}

		if testCase.userBearerToken == "" && testCase.appBearerToken != "" {
			expectedAuthHeader := "Bearer " + testCase.appBearerToken
			if req.Header.Get("Authorization") != expectedAuthHeader {
				t.Errorf("expected Authorization header to be \"%s\", got \"%s\"", expectedAuthHeader, req.Header.Get("Authorization"))
			}
		}

		if testCase.userBearerToken == "" && testCase.appBearerToken == "" {
			if req.Header.Get("Authorization") != "" {
				t.Error("did not expect Authorization header to be set")
			}
		}

		if req.Header.Get("User-Agent") != client.opts.UserAgent {
			t.Errorf("expected User-Agent header to be \"%s\", got \"%s\"", client.opts.UserAgent, req.Header.Get("User-Agent"))
		}
	}
}

func TestGetRateLimitHeaders(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode      int
		options         *Options
		Logins          []string
		respBody        string
		headerLimit     string
		headerRemaining string
		headerReset     string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			[]string{"summit1g"},
			`{"data":[{"id":"26490481","login":"summit1g","display_name":"summit1g","type":"","broadcaster_type":"partner","description":"I'm a competitive CounterStrike player who likes to play casually now and many other games. You will mostly see me play CS, H1Z1,and single player games at night. There will be many othergames played on this stream in the future as they come out:D","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/200cea12142f2384-profile_image-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/summit1g-channel_offline_image-e2f9a1df9e695ec1-1920x1080.png","view_count":202707885}]}`,
			"30",
			"29",
			"1529183210",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			[]string{"summit1g"},
			`{"data":[{"id":"26490481","login":"summit1g","display_name":"summit1g","type":"","broadcaster_type":"partner","description":"I'm a competitive CounterStrike player who likes to play casually now and many other games. You will mostly see me play CS, H1Z1,and single player games at night. There will be many othergames played on this stream in the future as they come out:D","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/200cea12142f2384-profile_image-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/summit1g-channel_offline_image-e2f9a1df9e695ec1-1920x1080.png","view_count":202707885}]}`,
			"",
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		mockRespHeaders := make(map[string]string)

		if testCase.headerLimit != "" {
			mockRespHeaders["Ratelimit-Limit"] = testCase.headerLimit
			mockRespHeaders["Ratelimit-Remaining"] = testCase.headerRemaining
			mockRespHeaders["Ratelimit-Reset"] = testCase.headerReset
		}

		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, mockRespHeaders))

		resp, err := c.GetUsers(&UsersParams{
			Logins: testCase.Logins,
		})
		if err != nil {
			t.Error(err)
		}

		expctedHeaderLimit, _ := strconv.Atoi(testCase.headerLimit)
		if resp.GetRateLimit() != expctedHeaderLimit {
			t.Errorf("expected \"Ratelimit-Limit\" to be \"%d\", got \"%d\"", expctedHeaderLimit, resp.GetRateLimit())
		}

		expctedHeaderRemaining, _ := strconv.Atoi(testCase.headerRemaining)
		if resp.GetRateLimitRemaining() != expctedHeaderRemaining {
			t.Errorf("expected \"Ratelimit-Remaining\" to be \"%d\", got \"%d\"", expctedHeaderRemaining, resp.GetRateLimitRemaining())
		}

		expctedHeaderReset, _ := strconv.Atoi(testCase.headerReset)
		if resp.GetRateLimitReset() != expctedHeaderReset {
			t.Errorf("expected \"Ratelimit-Reset\" to be \"%d\", got \"%d\"", expctedHeaderReset, resp.GetRateLimitReset())
		}
	}
}

type badMockHTTPClient struct {
	mockHandler http.HandlerFunc
}

func (mtc *badMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("Oops, that's bad :(")
}

func TestFailedHTTPClientDoRequest(t *testing.T) {
	t.Parallel()

	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}

	c := &Client{
		opts: options,
	}

	_, err := c.GetUsers(&UsersParams{
		Logins: []string{"summit1g"},
	})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDecodingBadJSON(t *testing.T) {
	t.Parallel()

	// Invalid JSON (missing `"` from the beginning)
	c := newMockClient(&Options{ClientID: "my-client-id"}, newMockHandler(http.StatusOK, `data":["some":"data"]}`, nil))

	_, err := c.GetUsers(&UsersParams{
		Logins: []string{"summit1g"},
	})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to decode API response: invalid character 'd' looking for beginning of value" {
		t.Error("expected error does match return error")
	}
}

func TestGetAppAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-app-access-token"

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetAppAccessToken(accessToken)

	if client.GetAppAccessToken() != accessToken {
		t.Errorf("expected GetAppAccessToken to return \"%s\", got \"%s\"", accessToken, client.GetAppAccessToken())
	}
}

func TestSetAppAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-app-access-token"

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetAppAccessToken(accessToken)

	if client.opts.AppAccessToken != accessToken {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", accessToken, client.opts.AppAccessToken)
	}
}

func TestGetUserAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-user-access-token"

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetUserAccessToken(accessToken)

	if client.GetUserAccessToken() != accessToken {
		t.Errorf("expected GetUserAccessToken to return \"%s\", got \"%s\"", accessToken, client.GetUserAccessToken())
	}
}

func TestSetUserAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-user-access-token"

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetUserAccessToken(accessToken)

	if client.opts.UserAccessToken != accessToken {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", accessToken, client.opts.UserAccessToken)
	}
}

func TestSetUserAgent(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36"
	client.SetUserAgent(userAgent)

	if client.opts.UserAgent != userAgent {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", userAgent, client.opts.UserAgent)
	}
}

func TestSetRedirectURI(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	redirectURI := "http://localhost/auth/callback"
	client.SetRedirectURI(redirectURI)

	if client.opts.RedirectURI != redirectURI {
		t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", redirectURI, client.opts.RedirectURI)
	}
}

func TestHydrateRequestCommon(t *testing.T) {
	t.Parallel()
	var sourceResponse Response
	sampleStatusCode := 200
	sampleHeaders := http.Header{}
	sampleHeaders.Set("Content-Type", "application/json")
	sampleError := "foo"
	sampleErrorStatus := 1
	sampleErrorMessage := "something done broke"
	sourceResponse.ResponseCommon.StatusCode = sampleStatusCode
	sourceResponse.ResponseCommon.Header = sampleHeaders
	sourceResponse.ResponseCommon.Error = sampleError
	sourceResponse.ResponseCommon.ErrorStatus = sampleErrorStatus
	sourceResponse.ResponseCommon.ErrorMessage = sampleErrorMessage

	var targetResponse Response
	sourceResponse.HydrateResponseCommon(&targetResponse.ResponseCommon)
	if targetResponse.StatusCode != sampleStatusCode {
		t.Errorf("expected StatusCode to be \"%d\", got \"%d\"", sampleStatusCode, targetResponse.ResponseCommon.StatusCode)
	}

	if targetResponse.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected headers to match")
	}

	if targetResponse.Error != sampleError {
		t.Errorf("expected Error to be \"%s\", got \"%s\"", sampleError, targetResponse.ResponseCommon.Error)
	}

	if targetResponse.ErrorStatus != sampleErrorStatus {
		t.Errorf("expected ErrorStatus to be \"%d\", got \"%d\"", sampleErrorStatus, targetResponse.ResponseCommon.ErrorStatus)
	}

	if targetResponse.ErrorMessage != sampleErrorMessage {
		t.Errorf("expected ErrorMessage to be \"%s\", got \"%s\"", sampleErrorMessage, targetResponse.ResponseCommon.ErrorMessage)
	}
}
