package helix

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetDropsEntitlements(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		gameID         string
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusUnauthorized,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "invalid-app-access-token", // invalid app access token
			},
			"1337",
			`{"error":"Unauthorized","status":401,"message":"Must provide valid app token."}`,
			"Must provide valid app token.",
		},
		{
			http.StatusForbidden,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token", // valid app access token
			},
			"1337",
			`{"error":"Forbidden","status":403,"message":"game not managed by this organization."}`,
			"game not managed by this organization.",
		},
		{
			http.StatusOK,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token",
			},
			"33214",
			`{"data": [{ "id": "fb78259e-fb81-4d1b-8333-34a06ffc24c0", "benefit_id": "74c52265-e214-48a6-91b9-23b6014e8041", "timestamp": "2019-01-28T04:17:53.325Z", "user_id": "25009227", "game_id": "33214", "fulfillment_status": "CLAIMED", "updated_at": "2019-01-28T04:17:53.325Z" }, { "id": "862750a5-265e-4ab6-9f0a-c64df3d54dd0", "benefit_id": "74c52265-e214-48a6-91b9-23b6014e8041", "timestamp": "2019-01-28T04:16:53.325Z", "user_id": "25009227", "game_id": "33214", "fulfillment_status": "CLAIMED", "updated_at": "2019-01-28T04:17:53.325Z" }, { "id": "d8879baa-3966-4d10-8856-15fdd62cce02", "benefit_id": "cdfdc5c3-65a2-43bc-8767-fde06eb4ab2c", "timestamp": "2019-01-28T04:15:53.325Z", "user_id": "25009227", "game_id": "33214", "fulfillment_status": "FULFILLED", "updated_at": "2019-01-28T04:18:53.325Z" }], "pagination": { "cursor": "eyJiIjpudW..." } }`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetDropsEntitlements(&GetDropEntitlementsParams{GameID: testCase.gameID})
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
		if len(resp.Data.Entitlements) < 1 {
			t.Errorf("Expected the number of entitlements to be a positive number")
		}

		if !strings.EqualFold(resp.Data.Entitlements[0].GameID, testCase.gameID) {
			t.Errorf("Expected the Entitlement's GameID to be the same as the requested GameID - wanted \"%s\", got \"%s\"",
				resp.Data.Entitlements[0].GameID, testCase.gameID)
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
	}

	_, err := c.GetDropsEntitlements(&GetDropEntitlementsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateDropsEntitlements(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode        int
		options           *Options
		requestedIDs      []string
		fulfillmentStatus string
		respBody          string
		expectedErrMsg    string
	}{
		{
			http.StatusUnauthorized,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "invalid-app-access-token", // invalid app access token
			},
			[]string{},
			"FULFILLED",
			`{"error":"Unauthorized","status":401,"message":"Must provide valid app token."}`,
			"Must provide valid app token.",
		},
		{
			http.StatusOK,
			&Options{
				ClientID:       "valid-client-id",
				AppAccessToken: "valid-app-access-token",
			},
			[]string{
				"fb78259e-fb81-4d1b-8333-34a06ffc24c0",
				"862750a5-265e-4ab6-9f0a-c64df3d54dd0",
				"d8879baa-3966-4d10-8856-15fdd62cce02",
				"9a290126-7e3b-4f66-a9ae-551537893b65",
			},
			"FULFILLED",
			`{ "data": [{ "status": "SUCCESS", "ids": ["fb78259e-fb81-4d1b-8333-34a06ffc24c0", "862750a5-265e-4ab6-9f0a-c64df3d54dd0"] }, { "status": "UNAUTHORIZED", "ids": ["d8879baa-3966-4d10-8856-15fdd62cce02"] }, { "status": "UPDATE_FAILED", "ids": ["9a290126-7e3b-4f66-a9ae-551537893b65"] }] }`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateDropsEntitlements(&UpdateDropsEntitlementsParams{
			EntitlementIDs:    testCase.requestedIDs,
			FulfillmentStatus: testCase.fulfillmentStatus,
		})
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
		if len(resp.Data.EntitlementSets) < 1 {
			t.Errorf("Expected the number of entitlements sets to be a positive number")
		}

		// Check all requested ids are in the response
		totatReponseIDs := 0
		for _, es := range resp.Data.EntitlementSets {
			for _, id := range es.IDs {
				if !stringInSlice(id, testCase.requestedIDs) {
					t.Errorf("Expected entitlement ID %v to be in requested IDs", id)
				}
			}
			totatReponseIDs += len(es.IDs)
		}
		if totatReponseIDs != len(testCase.requestedIDs) {
			t.Errorf("Expected number of entitlement IDs in response to match number of entitlement IDs requested")
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
	}

	_, err := c.UpdateDropsEntitlements(&UpdateDropsEntitlementsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
