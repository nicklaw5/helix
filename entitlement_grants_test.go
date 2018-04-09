package helix

import (
	"net/http"
	"strings"
	"testing"
)

func TestCreateEntitlementsUploadURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode      int
		options         *Options
		manifestID      string
		entitlementType string
		respBody        string
		expectedErrMsg  string
	}{
		{
			http.StatusUnauthorized,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "invalid-app-access-token", // invalid app access token
			},
			"my-manifest-id",
			"bulk_drops_grant",
			`{"error":"Unauthorized","status":401,"message":"Must provide valid app token."}`,
			"Must provide valid app token.",
		},
		{
			http.StatusBadRequest,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token",
			},
			"", // invalid manifest id
			"bulk_drops_grant",
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"manifest_id\""}`,
			"Missing required parameter \"manifest_id\"",
		},
		{
			http.StatusBadRequest,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token",
			},
			"my-manifest-id",
			"invalid_grant", // invalid entitlement type
			`{"error":"Bad Request","status":400,"message":"The parameter \"type\" was malformed: value must be one of \"bulk_drops_grant\""}`,
			"The parameter \"type\" was malformed: value must be one of \"bulk_drops_grant\"",
		},
		{
			http.StatusOK,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token",
			},
			"my-manifest-id",
			"bulk_drops_grant",
			`{"data":[{"url": "https://twitch-ds-vhs-drops-granted-uploads-us-west-2-prod.s3-us-west-2.amazonaws.com/valid-client-id/<time>-my-manifest-id.json?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=<credential>%2Fus-west-2%2Fs3%2Faws4_request\u0026X-Amz-Date=<date>\u0026X-Amz-Expires=900\u0026X-Amz-Security-Token=<token>\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=<signature>"}]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateEntitlementsUploadURL(testCase.manifestID, testCase.entitlementType)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		// Test error cases
		if resp.StatusCode != http.StatusOK {
			if resp.ErrorStatus != testCase.statusCode {
				t.Errorf("expected error status to be \"%d\", got \"%d\"", testCase.statusCode, resp.ErrorStatus)
			}

			if resp.ErrorMessage != testCase.expectedErrMsg {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", testCase.expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		// Test success case
		if len(resp.Data.URLs[0].URL) < 1 {
			t.Errorf("Expected URL not to be an empty string, got \"%s\"", resp.Data.URLs[0].URL)
		}

		if !strings.Contains(resp.Data.URLs[0].URL, testCase.options.ClientID) {
			t.Errorf("Expected URL to contain client id \"%s\", got \"%s\"", testCase.options.ClientID, resp.Data.URLs[0].URL)
		}

		if !strings.Contains(resp.Data.URLs[0].URL, testCase.manifestID) {
			t.Errorf("Expected URL to contain manifest id \"%s\", got \"%s\"", testCase.manifestID, resp.Data.URLs[0].URL)
		}
	}
}
