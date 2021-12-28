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

func TestGetExtensionLiveChannels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode     int
		options        *Options
		params         *ExtensionLiveChannelsParams
		respBody       string
		expectedErrMsg string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ExtensionLiveChannelsParams{ExtensionID: "some-extension-id"},
			`{ "data": [{ "broadcaster_id": "121086094", "broadcaster_name": "khaizer93", "game_name": "Art", "game_id": "509660", "title": "random artstream sketching Kiryu COCO" }, { "broadcaster_id": "165951395", "broadcaster_name": "MelloTodd", "game_name": "Jackbox Party Packs", "game_id": "493174", "title": "[OPEN] WInner Choses Next Game (1-8) | !jackbox !uptime #envtuber" }, { "broadcaster_id": "21724294", "broadcaster_name": "Mahoog47", "game_name": "Escape from Tarkov", "game_id": "491931", "title": "LVL 39| The holy relic Ash-12 has been acquired" }, { "broadcaster_id": "253663808", "broadcaster_name": "MrHatcher_", "game_name": "Dota 2", "game_id": "29595", "title": "Road to 53/55 followers! Dota 2 Ranked 1k mmr (british/filipino)" }, { "broadcaster_id": "245641098", "broadcaster_name": "ChoKoii", "game_name": "Escape from Tarkov", "game_id": "491931", "title": "First Drops Enabled Stream? | Labs Main | 1 Follower=5 push-ups" }, { "broadcaster_id": "268444856", "broadcaster_name": "D4RK_5KY", "game_name": "Always On", "game_id": "499973", "title": "24/7 FULLSEND HOST RAFFLE - Need THAT #SUPPORT!? #Affiliate PUSH!? Try Your LUCK \u0026 WIN The Raffle!" }, { "broadcaster_id": "42871388", "broadcaster_name": "mieudiary", "game_name": "Stardew Valley", "game_id": "490744", "title": "I'm very sleepy but let's farm | !melaomi" }, { "broadcaster_id": "429972112", "broadcaster_name": "andy_gra", "game_name": "twitch", "game_id": "", "title": "wbijaj smialo :)" }, { "broadcaster_id": "486154226", "broadcaster_name": "mrboone521", "game_name": "Escape from Tarkov", "game_id": "491931", "title": "TARK TARK offline and scav runs" }, { "broadcaster_id": "503028811", "broadcaster_name": "Uwlsy2k", "game_name": "Fortnite", "game_id": "33214", "title": "Bot Zonewars?! Join up and chat " }, { "broadcaster_id": "520878515", "broadcaster_name": "me_fon", "game_name": "Teamfight Tactics", "game_id": "513143", "title": "Ranking up in TFT Mob" }, { "broadcaster_id": "521301629", "broadcaster_name": "acrolic_", "game_name": "Apex Legends", "game_id": "511224", "title": "Come say hi! |road to 50 followers | song choices" }, { "broadcaster_id": "54270050", "broadcaster_name": "ELIASS_1", "game_name": "SCUM", "game_id": "495811", "title": "walking simulator 2022 | !setup | !sleep  | 386/400 followers | " }, { "broadcaster_id": "611701485", "broadcaster_name": "Semmy_22", "game_name": "Overwatch", "game_id": "488552", "title": "Support slave at your service " }, { "broadcaster_id": "63501619", "broadcaster_name": "KittySinisterr", "game_name": "Fortnite", "game_id": "33214", "title": "Winterfest challenges" }, { "broadcaster_id": "625059457", "broadcaster_name": "unisclan", "game_name": "Battlefield 2042", "game_id": "514974", "title": "crazy gameplay tanks   will not live " }, { "broadcaster_id": "604281079", "broadcaster_name": "viperarishyt", "game_name": "FIFA 22", "game_id": "1869092879", "title": "Grab Your Breakfast and Join me. Lets chat :-)" }, { "broadcaster_id": "666411722", "broadcaster_name": "SarahBree", "game_name": "Apex Legends", "game_id": "511224", "title": "Winter express only before it goes away :'(" }, { "broadcaster_id": "647613771", "broadcaster_name": "batbat0508", "game_name": "Battlefield 4", "game_id": "66402", "title": "( GOVS ) ~Fr-En ~ rules (LOCKER)" }, { "broadcaster_id": "653487605", "broadcaster_name": "honka2019", "game_name": "Identity V", "game_id": "508662", "title": "JPN♡本日も23:30頃までまったりプレイ⚠️mobile play" }], "pagination": "YVc1emRHRnNiQ00yTXpVd01UWXhPVHBsT1ROalpqZzNNekJ1WkRFeGVqZG5aWEJyYkhreVozSjVOV3QyT0dzNjoz" }`,		"",
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ExtensionLiveChannelsParams{ExtensionID: "", First: 2},
			`{"error":"Bad Request","status":400,"message":"missing extension id"}`,
			"error: extension ID must be specified",
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetExtensionLiveChannels(testCase.params)
		if err != nil {
			if err.Error() == testCase.expectedErrMsg {
				continue
			}

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

		if testCase.params.First != 0 && testCase.params.First != len(resp.Data.LiveChannels) {
			t.Errorf("expected %d transactions, got %d", testCase.params.First, len(resp.Data.LiveChannels))
		}
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
