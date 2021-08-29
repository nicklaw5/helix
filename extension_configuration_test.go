package helix

import (
	"net/http"
	"testing"
)

func TestSetExtensionConfigurationSegmentS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *ExtensionSetConfigurationParams
		respBody      string
		validationErr string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&ExtensionSetConfigurationParams{},
			`{"error":"Unauthorized","status":401,"message":"JWT token is missing"}`,
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ExtensionSetConfigurationParams{
				Segment:       ExtensionConfigurationGlobalSegment,
				ExtensionID:   "my-ext-id",
				BroadcasterID: "100249558",
				Version:       "ext-configversion",
				Content:       "{}",
			},
			`{"error":"Bad Request","status":400,"message":"error: developer or broadcaster extension configuration segment type must be provided for broadcasters"}`,
			"error: developer or broadcaster extension configuration segment type must be provided for broadcasters",
		},
		{
			http.StatusNoContent,
			&Options{
				ClientID: "my-client-id",
				ExtensionOpts: ExtensionOptions{
					Secret:      "my-ext-secret",
					OwnerUserID: "ext-owner-id",
				},
			},
			&ExtensionSetConfigurationParams{
				Segment:       ExtensionConfigrationBroadcasterSegment,
				ExtensionID:   "my-ext-id",
				BroadcasterID: "broadcasterId",
				Version:       "ext-configversion",
				Content:       "{}}",
			},
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SetExtensionSegmentConfig(testCase.params)
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

func TestGetExtensionConfigurationSegmentS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *ExtensionGetConfigurationParams
		respBody      string
		validationErr string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&ExtensionGetConfigurationParams{},
			`{"error":"Unauthorized","status":401,"message":"JWT token is missing"}`,
			"",
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ExtensionGetConfigurationParams{
				Segments:      []ExtensionSegmentType{ExtensionConfigurationGlobalSegment},
				ExtensionID:   "my-ext-id",
				BroadcasterID: "100249558",
			},
			`{"error":"Bad Request","status":400,"message":"error: developer or broadcaster extension configuration segment type must be provided"}`,
			"error: only developer or broadcaster extension configuration segment type must be provided for broadcasters",
		},
		{
			http.StatusNoContent,
			&Options{
				ClientID: "my-client-id",
				ExtensionOpts: ExtensionOptions{
					Secret:      "my-ext-secret",
					OwnerUserID: "ext-owner-id",
				},
			},
			&ExtensionGetConfigurationParams{
				Segments:    []ExtensionSegmentType{ExtensionConfigurationGlobalSegment},
				ExtensionID: "my-ext-id",
			},
			`{ "data": [{ "segment": "global", "broadcaster_id": "", "content": "{\n\"images\": {\n\"one\": \"https://i.giphy.com/media/NsEXpJpIt3lRWBcLol/source.gif\"\n},\n\"motd\": \"get gaming!\"\n}", "version": "3" }] }`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetExtensionConfigurationSegment(testCase.params)
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

		if len(resp.Data.Segments) != 1 {
			t.Errorf("expected 1 configuration segment got %d", len(resp.Data.Segments))
		}
	}
}

func TestExtensionSetConfigurationSegmentRequired(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode    int
		options       *Options
		params        *ExtensionSetRequiredConfigurationParams
		respBody      string
		validationErr string
	}{
		{
			http.StatusUnauthorized,
			&Options{ClientID: "my-client-id"},
			&ExtensionSetRequiredConfigurationParams{},
			`{"error":"Unauthorized","status":401,"message":"JWT token is missing"}`,
			"",
		},
		{
			http.StatusNoContent,
			&Options{
				ClientID: "my-client-id",
				ExtensionOpts: ExtensionOptions{
					Secret:      "my-ext-secret",
					OwnerUserID: "ext-owner-id",
				},
			},
			&ExtensionSetRequiredConfigurationParams{
				ExtensionID:           "my-ext-id",
				ExtensionVersion:      "0.0.1",
				RequiredConfiguration: "100249558",
				ConfigurationVersion:  "1",
			},
			"",
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SetExtensionRequiredConfiguration(testCase.params)
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
