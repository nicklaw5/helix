package helix

import (
	"context"
	"net/http"
	"testing"
)

func TestGetExtensionAnalytics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		extensionID    string
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusForbidden,
			&Options{ClientID: "my-client-id"},
			"493057",
			`{"error":"Forbidden","status":403,"message":"User Not Associated To Companies"}`,
			"User Not Associated To Companies",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			"493057",
			`{"data":[{"extension_id":"493057", "URL": "https://twitch-piper-reports.s3-us-west-2.amazonaws.com/dynamic/LoL%20ADC_overview_v2_2018-03-01_2018-06-01_8a879932-8e70-7a4c-2b97-e0eaba28c3b0.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAI7LOPgTrAIVYE6KQ%2F20180731%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20180731T202847Z&X-Amz-Expires=60&X-Amz-Security-Token=FQoDYXdzEDMaDObdwyOVISdo6feHSSK3A9R9gMFeS5frG5Dsr4k4tJemqCIazhsQJrpsehBoOufaQkCxrb8RD3oU0xC5pWrZe9kN%2BnezIoLOgTtFRAqTzdIr7J5iUOxGFyKN9XmrmUHGexFfALvoPQWUJNbxoFU6shajSmO3sPK2GnuEaGmIrAqjKrim8saLHDV%2FdSi2ZH3fFx6sBQEGv13Lx0zua7AsvaL%2BSfhIAcOazWjYLMU5N9bxXmaN7IAIF4UjNPqbg07RMWW70hm0edH0RPi%2Fw00faeeSvmreHq6c1C1Lu8a7AysMb0pEGBT7VxmuGmWsXyjLWZ6oNgbx88HXoMJpmAn5Y1hUu7VzOaa84T%2BmCF5Sbn7hbB1xIiPdzaVQ%2Bd85sy4ln09h7dgKh6GFE1VTas2v7RJU1lyD%2FZ%2FWKBwV5Ol8GEGrF1pme8mSBpPGUAJ4vxjLmrGL7ctty%2F0vXke3PyD%2B4%2FtHZ67xaw0y8EKrau23Xvt3blkcDNoQYOfcS%2FqbaK%2BHpyVq4bIBtQq%2BHYU5MuFkgEuwSe5zPDle1ysKSN11B6B6Sy7Httrq542OONS%2BfURkczMbKSPEShddN32Y9VUqKYdUo%2FsWVQQoy7uC2wU%3D&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3Bfilename%3D%22WoW%20Armory_overview_v1_2018-04-30_2018-06-01.csv%22&X-Amz-Signature=eb7721e40cbfd1d7409887dae3792cdb2add025ace953a63ba8e3545b92ae058"}]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetExtensionAnalytics(&ExtensionAnalyticsParams{
			ExtensionID: testCase.extensionID,
			First:       1,
			Type:        "overview_v1",
		})
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

		if resp.Data.ExtensionAnalytics[0].ExtensionID != testCase.extensionID {
			t.Errorf("expected extension id to be \"%s\", got \"%s\"", testCase.extensionID, resp.Data.ExtensionAnalytics[0].ExtensionID)
		}

		if len(resp.Data.ExtensionAnalytics[0].URL) < 1 {
			t.Errorf("expected extension analytics url not to be an empty string, got \"%s\"", resp.Data.ExtensionAnalytics[0].URL)
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

	_, err := c.GetExtensionAnalytics(&ExtensionAnalyticsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetGameAnalytics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		params         *GameAnalyticsParams
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusForbidden,
			&Options{ClientID: "my-client-id"},
			&GameAnalyticsParams{GameID: "493057", Type: "overview_v2"},
			`{"error":"Forbidden","status":403,"message":"User Not Associated To Companies"}`,
			"User Not Associated To Companies",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GameAnalyticsParams{GameID: "493057", Type: "overview_v2"},
			`{"data":[{"game_id":"493057","URL":"https://twitch-piper-reports.s3-us-west-2.amazonaws.com/games/66170/overview/1518307200000.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAJP7WFIAF26K7BC2Q%2F20180222%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20180222T220651Z&X-Amz-Expires=60&X-Amz-Security-Token=FQoDYXdzEE0aDLKNl9aCgfuikMKI%2ByK3A4e%2FR%2B4to%2BmRZFUuslNKs%2FOxKeySB%2BAU87PBtNGCxQaQuN2Q8KI4Vg%2Bve2x5eenZdoH0ZM7uviM94sf2GlbE9Z0%2FoJRmNGNhlU3Ua%2FupzvByCoMdefrU8Ziiz4j8EJCgg0M1j2aF9f8bTC%2BRYwcpP0kjaZooJS6RFY1TTkh659KBA%2By%2BICdpVK0fxOlrQ%2FfZ6vIYVFzvywBM05EGWX%2F3flCIW%2BuZ9ZxMAvxcY4C77cOLQ0OvY5g%2F7tuuGSO6nvm9Eb8MeMEzSYPr4emr3zIjxx%2Fu0li9wjcF4qKvdmnyk2Bnd2mepX5z%2BVejtIGBzfpk%2Fe%2FMqpMrcONynKoL6BNxxDL4ITo5yvVzs1x7OumONHcsvrTQsd6aGNQ0E3lrWxcujBAmXmx8n7Qnk4pZnHZLgcBQam1fIGba65Gf5Ern71TwfRUsolxnyIXyHsKhd2jSmXSju8jH3iohjv99a2vGaxSg8SBCrQZ06Bi0pr%2FTiSC52U1g%2BlhXYttdJB4GUdOvaxR8n6PwMS7HuAtDJUui8GKWK%2F9t4OON3qhF2cBt%2BnV%2BDg8bDMZkQ%2FAt5blvIlg6rrlCu0cYko4ojb281AU%3D&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3Bfilename%3DWarframe-overview-2018-02-11.csv&X-Amz-Signature=49cc07cbd9d753b00315b66f49b9e4788570062ff3bd956288ab4f164cf96708","type":"overview_v2"}]}`,
			"",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetGameAnalytics(testCase.params)
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

		if resp.Data.GameAnalytics[0].GameID != testCase.params.GameID {
			t.Errorf("expected game id to be \"%s\", got \"%s\"", testCase.params.GameID, resp.Data.GameAnalytics[0].GameID)
		}

		if len(resp.Data.GameAnalytics[0].URL) < 1 {
			t.Errorf("expected game analytics url not to be an empty string, got \"%s\"", resp.Data.GameAnalytics[0].URL)
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

	_, err := c.GetGameAnalytics(&GameAnalyticsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "Failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}
