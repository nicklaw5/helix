package helix

import (
	"net/http"
	"net/http/httptest"
)

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

type mockHTTPClient struct {
	mockHandler func(http.ResponseWriter, *http.Request)
}

func (mtc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mtc.mockHandler)
	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}
