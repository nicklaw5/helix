package helix

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockHTTPClient struct {
	mockHandler func(http.ResponseWriter, *http.Request)
}

func (mtc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mtc.mockHandler)
	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func newMockClient(clientID string, mockHandler func(http.ResponseWriter, *http.Request)) *Client {
	mc := &Client{}
	mc.clientID = clientID
	mc.httpClient = &mockHTTPClient{mockHandler}

	return mc
}

func newMockHandler(statusCode int, json string, headers map[string]string) func(http.ResponseWriter, *http.Request) {
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
				ClientID:      "my-client-id",
				ClientSecret:  "my-client-secret",
				HTTPClient:    &http.Client{},
				AccessToken:   "my-access-token",
				UserAgent:     "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
				RateLimitFunc: func(*Response) error { return nil },
				Scopes:        []string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
				RedirectURI:   "http://localhost/auth/callback",
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

		if client.clientID != testCase.options.ClientID {
			t.Errorf("expected clientID to be \"%s\", got \"%s\"", testCase.options.ClientID, client.clientID)
		}

		if client.clientSecret != testCase.options.ClientSecret {
			t.Errorf("expected clientSecret to be \"%s\", got \"%s\"", testCase.options.ClientSecret, client.clientSecret)
		}

		if reflect.TypeOf(client.rateLimitFunc).Kind() != reflect.Func {
			t.Errorf("expected rateLimitFunc to be a function, got %+v", reflect.TypeOf(client.rateLimitFunc).Kind())
		}

		if client.httpClient != testCase.options.HTTPClient {
			t.Errorf("expected httpClient to be \"%s\", got \"%s\"", testCase.options.HTTPClient, client.httpClient)
		}

		if client.userAgent != testCase.options.UserAgent {
			t.Errorf("expected userAgent to be \"%s\", got \"%s\"", testCase.options.UserAgent, client.userAgent)
		}

		if client.accessToken != testCase.options.AccessToken {
			t.Errorf("expected accessToken to be \"%s\", got \"%s\"", testCase.options.AccessToken, client.accessToken)
		}

		if len(client.scopes) != len(testCase.options.Scopes) {
			t.Errorf("expected \"%d\" scopes, got \"%d\"", len(testCase.options.Scopes), len(client.scopes))
		}

		if client.redirectURI != testCase.options.RedirectURI {
			t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", testCase.options.RedirectURI, client.redirectURI)
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

	if client.clientID != options.ClientID {
		t.Errorf("expected clientID to be \"%s\", got \"%s\"", options.ClientID, client.clientID)
	}

	if client.clientSecret != "" {
		t.Errorf("expected clientSecret to be \"%s\", got \"%s\"", options.ClientSecret, client.clientSecret)
	}

	if client.userAgent != "" {
		t.Errorf("expected userAgent to be \"%s\", got \"%s\"", "", client.userAgent)
	}

	if client.accessToken != "" {
		t.Errorf("expected accesstoken to be \"\", got \"%s\"", client.accessToken)
	}

	if client.httpClient != http.DefaultClient {
		t.Errorf("expected httpClient to be \"%v\", got \"%v\"", http.DefaultClient, client.httpClient)
	}

	if client.rateLimitFunc != nil {
		t.Errorf("expected rateLimitFunc to be \"%v\", got \"%v\"", nil, client.rateLimitFunc)
	}

	if len(client.scopes) != len(options.Scopes) {
		t.Errorf("expected \"%d\" scopes, got \"%d\"", len(options.Scopes), len(client.scopes))
	}

	if client.redirectURI != options.RedirectURI {
		t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", options.RedirectURI, client.redirectURI)
	}
}

func TestSetRequestHeaders(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Options{
		ClientID:  "cid",
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetAccessToken("my-access-token")

	req, _ := http.NewRequest("GET", "/blah/blah", nil)

	client.setRequestHeaders(req)

	expectedAuthHeader := "Bearer " + client.accessToken
	if req.Header.Get("Authorization") != expectedAuthHeader {
		t.Errorf("expected Authorization header to be \"%s\", got \"%s\"", expectedAuthHeader, req.Header.Get("Authorization"))
	}

	if req.Header.Get("User-Agent") != client.userAgent {
		t.Errorf("expected User-Agent header to be \"%s\", got \"%s\"", client.userAgent, req.Header.Get("User-Agent"))
	}
}

func TestSetAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-access-token"

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	client.SetAccessToken(accessToken)

	if client.accessToken != accessToken {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", accessToken, client.accessToken)
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

	if client.userAgent != userAgent {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", userAgent, client.accessToken)
	}
}

func TestSetScopes(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Options{
		ClientID: "cid",
	})
	if err != nil {
		t.Errorf("Did not expect an error, got \"%s\"", err.Error())
	}

	scopes := []string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"}
	client.SetScopes(scopes)

	if len(client.scopes) != len(scopes) {
		t.Errorf("expected \"%d\" scopes, got \"%d\"", len(scopes), len(client.scopes))
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

	if client.redirectURI != redirectURI {
		t.Errorf("expected redirectURI to be \"%s\", got \"%s\"", redirectURI, client.redirectURI)
	}
}
