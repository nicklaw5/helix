package helix

import (
	"net/http"
	"testing"
)

func TestGetPredictions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode        int
		options           *Options
		PredictionsParams *PredictionsParams
		respBody          string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&PredictionsParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&PredictionsParams{BroadcasterID: "121445595"},
			`{"data":[{"id":"0eb6f04e-8867-4f9f-8bdd-e4038556d0e8","broadcaster_id":"145328278","broadcaster_name":"Scorfly","broadcaster_login":"scorfly","title":"test 1","winning_outcome_id":"760f8303-5a4f-420e-8649-527752447e0f","outcomes":[{"id":"760f8303-5a4f-420e-8649-527752447e0f","title":"choice blue","users":1,"channel_points":10,"top_predictors":[{"user_id":"250117050","user_login":"botvause","user_name":"BotVause","channel_points_used":10,"channel_points_won":10}],"color":"BLUE"},{"id":"6c0fe617-1309-4bb7-ad03-9c1b232f2251","title":"choice pink","users":0,"channel_points":0,"top_predictors":null,"color":"PINK"}],"prediction_window":30,"status":"RESOLVED","created_at":"2021-05-07T21:30:28.20509235Z","ended_at":"2021-05-07T21:32:12.402517544Z","locked_at":"2021-05-07T21:30:57.242055129Z"}],"pagination":{}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetPredictions(testCase.PredictionsParams)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
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

	_, err := c.GetPredictions(&PredictionsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestCreatePrediction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode             int
		options                *Options
		CreatePredictionParams *CreatePredictionParams
		respBody               string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&CreatePredictionParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&CreatePredictionParams{
				BroadcasterID: "145328278",
				Title:         "Test",
				Outcomes: []PredictionChoiceParam{
					PredictionChoiceParam{Title: "Panda"},
					PredictionChoiceParam{Title: "Tiger"},
				},
				PredictionWindow: 300,
			},
			`{"data":[{"id":"92bdcb5c-6d83-4c75-95d6-fdd34f128d43","broadcaster_id":"145328278","broadcaster_name":"Scorfly","broadcaster_login":"scorfly","title":"Test","winning_outcome_id":null,"outcomes":[{"id":"6afe5daf-e54c-48d7-9c57-d07e791c496b","title":"choix 1","users":0,"channel_points":0,"top_predictors":null,"color":"BLUE"},{"id":"d8ffec60-f87f-44ac-b7b1-c53001bf2e4b","title":"choix 2","users":0,"channel_points":0,"top_predictors":null,"color":"PINK"}],"prediction_window":300,"status":"ACTIVE","created_at":"2021-05-07T22:15:34.457301028Z","ended_at":null,"locked_at":null}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreatePrediction(testCase.CreatePredictionParams)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
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

	_, err := c.CreatePrediction(&CreatePredictionParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestEndPrediction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode          int
		options             *Options
		EndPredictionParams *EndPredictionParams
		respBody            string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&EndPredictionParams{BroadcasterID: ""},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&EndPredictionParams{
				BroadcasterID:    "145328278",
				ID:               "92bdcb5c-6d83-4c75-95d6-fdd34f128d43",
				Status:           "RESOLVED",
				WinningOutcomeID: "6afe5daf-e54c-48d7-9c57-d07e791c496b",
			},
			`{"data":[{"id":"92bdcb5c-6d83-4c75-95d6-fdd34f128d43","broadcaster_id":"145328278","broadcaster_name":"Scorfly","broadcaster_login":"scorfly","title":"Test","winning_outcome_id":"6afe5daf-e54c-48d7-9c57-d07e791c496b","outcomes":[{"id":"6afe5daf-e54c-48d7-9c57-d07e791c496b","title":"choix 1","users":0,"channel_points":0,"top_predictors":null,"color":"BLUE"},{"id":"d8ffec60-f87f-44ac-b7b1-c53001bf2e4b","title":"choix 2","users":0,"channel_points":0,"top_predictors":null,"color":"PINK"}],"prediction_window":300,"status":"RESOLVED","created_at":"2021-05-07T22:15:34.457301028Z","ended_at":"2021-05-07T22:18:14.015776526Z","locked_at":null}]}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.EndPrediction(testCase.EndPredictionParams)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
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

	_, err := c.EndPrediction(&EndPredictionParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
