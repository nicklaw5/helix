package main

import (
	"fmt"
	"log"

	"github.com/nicklaw5/helix"
)

func main() {
	// var topic helix.WebhookTopic = 1

	// fmt.Printf("%+v\n", topic)

	// fmt.Printf("%+v\n", helix.UserFollowsRegexp.MatchString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1336&to_id=1337>; rel=\"self\""))

	// fmt.Printf("%+v\n", helix.StreamChangedRegexp.MatchString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/streams?user_id=104137656>; rel=\"self\""))

	// fmt.Printf("%+v\n", helix.UserChangedRegexp.MatchString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users?id=1234>; rel=\"self\""))

	// fmt.Printf("%+v\n", helix.GameAnalyticsRegexp.MatchString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?game_id=1234>; rel=\"self\""))

	// fmt.Printf("%+v\n", helix.ExtensionAnalyticsRegexp.MatchString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?extension_id=1234>; rel=\"self\""))

	// topic, err := helix.GetWebhookTopicFromLinkHeader("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=1336&to_id=1337>; rel=\"self\"")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", topic)

	// matches := helix.ExtensionAnalyticsRegexp.FindAllString("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/analytics?extension_id=1234>; rel=\"self\"", -1)
	// fmt.Printf("%+v\n", matches)

	// matches := helix.GetWebhookTopicValuesFromLinkHeader("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&from_id=11111&to_id=2222>; rel=\"self\"", helix.UserFollowsTopic)
	// matches := helix.GetWebhookTopicValuesFromLinkHeader("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", https://api.twitch.tv/helix/streams?user_id=104137656>; rel=\"self\"", helix.StreamChangedTopic)
	// fmt.Printf("%+v\n", matches)

	// for i, match := range matches {
	// 	fmt.Printf("%s: %s\n", i, match)
	// }

	// matches := helix.GetWebhookTopicValuesFromLinkHeader("<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/streams?user_id=104137656>; rel=\"self\"", helix.StreamChangedTopic)
	// fmt.Printf("%+v\n", matches)

	c, err := helix.NewClient(&helix.Options{
		ClientID:        "",
		ClientSecret:    "",
		RedirectURI:     "http://localhost:8888/auth/callback",
		UserAccessToken: "",
		// AppAccessToken: "",
		Scopes: []string{"user:read:email"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := c.GetAppAccessToken()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v\n", resp)

	// GET AUTH URL
	// authURL := c.GetAuthorizationURL("some-state", false)
	// fmt.Printf("%+v\n", authURL)
	// os.Exit(0)

	// GET USER ACCESS TOKEN
	// code := ""
	// resp, err := c.GetUserAccessToken(code)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v\n", resp)

	// GET USERS
	resp, _ := c.GetUsers(&helix.UsersParams{
		Logins: []string{"carpetsausage"},
		// 	// IDs: []string{"127506955"},
	})
	fmt.Printf("%+v\n", resp)

	// GET GAMES
	// resp, _ := c.GetGames(&helix.GamesParams{
	// 	Names: []string{"Sea of Thieves"},
	// })
	// fmt.Printf("%+v\n", resp)

	// GET STREAMS
	// resp, _ := c.GetStreams(&helix.StreamsParams{
	// 	First: 20,
	// 	// Language: []string{"en"},
	// 	// Type: "vodcast",
	// })
	// fmt.Printf("%+v\n", resp)

	// GET STREAM MARKERS
	// resp, _ := c.GetStreamMarkers(&helix.StreamMarkersParams{
	// 	First: 1,
	// 	// VideoID: "342339273",
	// 	UserID: "104137656",
	// })
	// fmt.Printf("%+v\n", resp)

	// CREATE CLIP
	// resp, _ := c.CreateClip(&helix.CreateClipParams{
	// 	BroadcasterID: "31557869",
	// 	HasDelay:      true,
	// })
	// fmt.Printf("%+v\n", resp)

	// // GET CLIPS
	// resp, _ := c.GetClips(&helix.ClipsParams{
	// 	// IDs: []string{"poop"},
	// 	// GameID: "1234",
	// 	BroadcasterID: "26490481",
	// 	First:         2,
	// })
	// fmt.Printf("%+v\n", resp)

	// SUBMIT WEBBHOOK SUBSSCRIPTION
	// resp, _ := c.PostWebhookSubscription(&helix.WebhookSubscriptionPayload{
	// 	Callback:     "https://webhooks.chatstatz.com/twitch/webhooks",
	// 	LeaseSeconds: 864000, // 10 days
	// 	Mode:         "subscribe",
	// 	Topic:        "https://api.twitch.tv/helix/streams?user_id=104137656",
	// 	Secret:       "poop",
	// })
	// fmt.Printf("%+v\n", resp)

	// SUBMIT WEBBHOOK UNSUBSCRIBE
	// resp, _ := c.PostWebhookSubscription(&helix.WebhookSubscriptionPayload{
	// 	Callback:     "https://webhooks.chatstatz.com/twitch/webhooks",
	// 	LeaseSeconds: 864000, // 10 days
	// 	Mode:         "unsubscribe",
	// 	Topic:        "https://api.twitch.tv/helix/streams?user_id=104137656",
	// 	Secret:       "poop",
	// })
	// fmt.Printf("%+v\n", resp)

	// GET WEBHOOK SUBSCRIPTIONS
	// resp, err := c.GetWebhookSubscriptions(&helix.WebhookSubscriptionsParams{
	// 	First: 10,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v\n", resp)
}
