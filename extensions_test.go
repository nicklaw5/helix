package helix

import (
	"net/http"
	"testing"
)

func TestGetExtensionTransactions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		params         *ExtensionTransactionsParams
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ExtensionTransactionsParams{ExtensionID: "some-extension-id"},
			`{"data":[{"id":"74c52265-e214-48a6-91b9-23b6014e8041","timestamp":"2019-01-28T04:15:53.325Z","broadcaster_id":"439964613","broadcaster_login":"chikuseuma","broadcaster_name":"chikuseuma","user_id":"424596340","user_login":"quotrok","user_name":"quotrok","product_type":"BITS_IN_EXTENSION","product_data":{"sku":"testSku100","cost":{"amount":100,"type":"bits"},"displayName":"Test Sku","inDevelopment":false}}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6M319"}}`,
			"",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ExtensionTransactionsParams{ExtensionID: "some-extension-id", ID: []string{"74c52265-e214-48a6-91b9-23b6014e8041", "8d303dc6-a460-4945-9f48-59c31d6735cb"}, First: 2},
			`{"data":[{"id":"74c52265-e214-48a6-91b9-23b6014e8041","timestamp":"2019-01-28T04:15:53.325Z","broadcaster_id":"439964613","broadcaster_login":"chikuseuma","broadcaster_name":"chikuseuma","user_id":"424596340","user_login":"quotrok","user_name":"quotrok","product_type":"BITS_IN_EXTENSION","product_data":{"sku":"testSku100","cost":{"amount":100,"type":"bits"},"displayName":"Test Sku","inDevelopment":false}},{"id":"8d303dc6-a460-4945-9f48-59c31d6735cb","timestamp":"2019-01-18T09:10:13.397Z","broadcaster_id":"439964613","broadcaster_login":"chikuseuma","broadcaster_name":"chikuseuma","user_id":"439966926","user_login":"liscuit","user_name":"liscuit","product_type":"BITS_IN_EXTENSION","product_data":{"sku":"testSku100","cost":{"amount":100,"type":"bits"},"displayName":"Test Sku","inDevelopment":false}}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6M319"}}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetExtensionTransactions(testCase.params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusForbidden {
			if resp.Error != "Forbidden" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusForbidden {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusForbidden, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if testCase.params.First != 0 && testCase.params.First != len(resp.Data.ExtensionTransactions) {
			t.Errorf("expected %d transactions, got %d", testCase.params.First, len(resp.Data.ExtensionTransactions))
		}

		if testCase.params.ID != nil {
			for _, ID := range testCase.params.ID {
				found := false
				for _, txn := range resp.Data.ExtensionTransactions {
					if txn.ID == ID {
						found = true
					}
				}

				if !found {
					t.Errorf("expected response to conatin transaction id %s, but didn't", ID)
				}
			}
		}
	}

	// Test with HTTP Failure
	c, err := NewClient(&Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	})
	if err != nil {
		t.Error(err)
	}

	_, err = c.GetExtensionTransactions(&ExtensionTransactionsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestExtensionSendChatMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *ExtensionSendChatMessageParams
		respBody      string
		validationErr string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ExtensionSendChatMessageParams{
				Text:             "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore",
				BroadcasterID:    "100249558",
				ExtensionVersion: "0.0.1",
				ExtensionID:      "my-ext-id",
			},
			`{"error":"Bad Request","status":400,"message":"text exceeds 280 characters"}`,
			"error: chat message length exceeds 280 characters",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ExtensionSendChatMessageParams{
				Text:             "welcome to the stream",
				ExtensionVersion: "0.0.1",
				ExtensionID:      "my-ext-id",
			},
			`{"error":"Bad Request","status":400,"message":"missing broadcaster id"}`,
			"error: broadcaster ID must be specified",
		},
		{
			http.StatusOK,
			&Options{
				ClientID: "my-client-id",
				ExtensionOpts: ExtensionOptions{
					Secret:      "my-ext-secret",
					OwnerUserID: "ext-owner-id",
				},
			},
			&ExtensionSendChatMessageParams{
				ExtensionID:      "my-ext-id",
				Text:             "welcome to the stream!",
				ExtensionVersion: "0.0.1",
				BroadcasterID:    "100249558",
			},
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SendExtensionChatMessage(testCase.params)
		if err != nil {
			if err.Error() == testCase.validationErr {
				continue
			}

			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be \"%s\", got \"%s\"", "Unauthorized", resp.Error)
			}

			if resp.ErrorStatus != http.StatusUnauthorized {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", http.StatusUnauthorized, resp.ErrorStatus)
			}

			expectedErrMsg := "JWT token is missing"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}
	}
}
