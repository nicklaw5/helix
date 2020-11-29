package helix

import (
	"net/http"
	"testing"
)

func TestClient_GetEntitlementCodeStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		count      int
		userId     string
		codes      []string
		respBody   string
	}{
		// TODO: expand with other test cases, including negative scenarios
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			"156900877",
			[]string{"KUHXV-4GXYP-AKAKK", "XZDDZ-5SIQR-RT5M3"},
			`{"data":[{"code":"KUHXV-4GXYP-AKAKK","status":"UNUSED"},{"code":"XZDDZ-5SIQR-RT5M3","status":"ALREADY_CLAIMED"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		params := &CodesParams{
			UserID: testCase.userId,
			Codes:  testCase.codes,
		}

		resp, err := c.GetEntitlementCodeStatus(params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// TODO: Test error cases

		// Test success cases
		if len(resp.Data.Codes) != testCase.count {
			t.Errorf("expected number of results to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.Codes))
		}

		for index, code := range testCase.codes {
			codes := resp.Data.Codes[index]
			if codes.Code != code {
				t.Error("expected entitlement code \"#{code}\", got \"#{codes.Code}\"")
			}
		}

	}
}

func TestClient_RedeemEntitlementCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		count      int
		userId     string
		codes      []string
		respBody   string
	}{
		// TODO: expand with other test cases, including negative scenarios
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			"156900877",
			[]string{"KUHXV-4GXYP-AKAKK", "XZDDZ-5SIQR-RT5M3"},
			`{"data":[{"code":"KUHXV-4GXYP-AKAKK","status":"SUCCESSFULLY_REDEEMED"},{"code":"XZDDZ-5SIQR-RT5M3","status":"ALREADY_CLAIMED"}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		params := &CodesParams{
			UserID: testCase.userId,
			Codes:  testCase.codes,
		}

		resp, err := c.GetEntitlementCodeStatus(params)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// TODO: Test error cases

		// Test success cases
		if len(resp.Data.Codes) != testCase.count {
			t.Errorf("expected number of results to be \"%d\", got \"%d\"", testCase.count, len(resp.Data.Codes))
		}

		for index, code := range testCase.codes {
			codes := resp.Data.Codes[index]
			if codes.Code != code {
				t.Error("expected entitlement code \"#{code}\", got \"#{codes.Code}\"")
			}
		}

	}
}
