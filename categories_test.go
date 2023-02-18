package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestSearchCategories(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		First      int
		respBody   string
		parsed     []Category
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			2,
			`{"data":[{"name":"Pocket Monster","id":"390264144","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/390264144_IGDB-52x72.jpg"},{"name":"Pokémon Black/White","id":"27602","box_art_url":"https://static-cdn.jtvnw.net/ttv-boxart/27602-52x72.jpg"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6Mn19"}}`,
			[]Category{
				{
					Name:      "Pocket Monster",
					ID:        "390264144",
					BoxArtURL: "https://static-cdn.jtvnw.net/ttv-boxart/390264144_IGDB-52x72.jpg",
				},
				{
					Name:      "Pokémon Black/White",
					ID:        "27602",
					BoxArtURL: "https://static-cdn.jtvnw.net/ttv-boxart/27602-52x72.jpg",
				},
			},
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			101,
			`{"error":"Bad Request","status":400,"message":"The parameter \"first\" was malformed: the value must be less than or equal to 100"}`,
			[]Category{},
		},
	}
	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.SearchCategories(&SearchCategoriesParams{
			First: testCase.First,
		})
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "The parameter \"first\" was malformed: the value must be less than or equal to 100"
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if len(resp.Data.Categories) != testCase.First {
			t.Errorf("expected \"%d\" streams, got \"%d\"", testCase.First, len(resp.Data.Categories))
		}

		for i, category := range resp.Data.Categories {
			if category.ID != testCase.parsed[i].ID {
				t.Errorf("Expected struct field ID = %s, was %s", testCase.parsed[i].ID, category.ID)
			}
			if category.Name != testCase.parsed[i].Name {
				t.Errorf("Expected struct field ID = %s, was %s", testCase.parsed[i].Name, category.Name)
			}
			if category.BoxArtURL != testCase.parsed[i].BoxArtURL {
				t.Errorf("Expected struct field ID = %s, was %s", testCase.parsed[i].BoxArtURL, category.BoxArtURL)
			}
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
		ctx:  context.Background(),
	}

	_, err := c.SearchCategories(&SearchCategoriesParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
