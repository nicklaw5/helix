package helix

import (
	"net/http"
	"testing"
)

func TestGetUserExtensions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			`{"data":[{"id":"wi08ebtatdc7oj83wtl9uxwz807l8b","version":"1.1.8","name":"Streamlabs Leaderboard","can_activate":true,"type":["panel"]},{"id":"d4uvtfdr04uq6raoenvj7m86gdk16v","version":"2.0.2","name":"Prime Subscription and Loot Reminder","can_activate":true,"type":["overlay"]}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUserExtensions()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode == http.StatusOK && len(resp.Data.UserExtensions) == 0 {
			t.Error("failed to parse successful UserExtension response data")
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}
}

func TestGetUserActiveExtensions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id", UserAccessToken: "my-access-token"},
			`{"data":{"panel":{"1":{"active":true,"id":"rh6jq1q334hqc2rr1qlzqbvwlfl3x0","version":"1.1.0","name":"TopClip"},"2":{"active":true,"id":"wi08ebtatdc7oj83wtl9uxwz807l8b","version":"1.1.8","name":"Streamlabs Leaderboard"},"3":{"active":true,"id":"naty2zwfp7vecaivuve8ef1hohh6bo","version":"1.0.9","name":"Streamlabs Stream Schedule & Countdown"}},"overlay":{"1":{"active":true,"id":"zfh2irvx2jb4s60f02jq0ajm8vwgka","version":"1.0.19","name":"Streamlabs"}},"component":{"1":{"active":true,"id":"lqnf3zxk0rv0g7gq92mtmnirjz2cjj","version":"0.0.1","name":"Dev Experience Test","x":0,"y":0},"2":{"active":false}}}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUserActiveExtensions(nil)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode == http.StatusOK {
			data := resp.Data.UserActiveExtensions
			if data.Component == nil || data.Panel == nil || data.Overlay == nil {
				t.Error("failed to parse successful UserActiveExtension response data")
			}
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}
}
