# helix

A Twitch Helix API client written in Go (Golang).

[![Tests and Coverage](https://github.com/nicklaw5/helix/workflows/Tests%20and%20Coverage/badge.svg)](https://github.com/nicklaw5/helix/actions?query=workflow%3A%22Tests+and+Coverage%22)
[![Coverage Status](https://coveralls.io/repos/github/nicklaw5/helix/badge.svg)](https://coveralls.io/github/nicklaw5/helix)

## Package Status

This project is a work in progress. Twitch has not finished all available endpoints/features for the Helix
API, but as these get released they are likely to be implemented in this package.

## Documentation & Examples

All documentation and usage examples for this package can be found in the [docs directory](docs). If you are
looking for the Twitch API docs, see the [Twitch Developer website](https://dev.twitch.tv/docs/api).

## Supported Endpoints & Features

**Authentication:**

- [x] Generate Authorization URL ("code" or "token" authorization)
- [x] Get App Access Tokens (OAuth Client Credentials Flow)
- [x] Get User Access Tokens (OAuth Authorization Code Flow)
- [x] Refresh User Access Tokens
- [x] Revoke User Access Tokens
- [x] Validate Access Token

**API Endpoint:**

- [x] Start Commercial
- [x] Get Extension Analytics
- [x] Get Game Analytics
- [x] Get Bits Leaderboard
- [x] Get Cheermotes
- [ ] Get Extension Transactions
- [x] Get Channel Information
- [x] Modify Channel Information
- [x] Get Channel Editors
- [ ] Create Custom Rewards
- [ ] Delete Custom Reward
- [ ] Get Custom Reward
- [ ] Get Custom Reward Redemption
- [ ] Update Custom Reward
- [ ] Update Redemption Status
- [x] Create Clip
- [x] Get Clip
- [x] Get Clips
- [x] Create Entitlement Grants Upload URL
- [x] Get Code Status
- [x] Get Drops Entitlements
- [x] Redeem Code
- [x] Create / Remove / List EventSub Subscriptions
- [x] Get Top Games
- [x] Get Games
- [ ] Get Hype Train Events
- [ ] Check AutoMod Status
- [ ] Get Banned Events
- [x] Get Banned Users
- [ ] Get Moderators
- [ ] Get Moderator Events
- [ ] Search Categories
- [ ] Search Channels
- [ ] Get Stream Key
- [x] Get Streams
- [x] Get Followed Streams
- [x] Create Stream Marker
- [x] Get Stream Markers
- [x] Get Broadcaster Subscriptions
- [x] Check User Subscription
- [ ] Get All Stream Tags
- [ ] Get Stream Tags
- [ ] Replace Stream Tags
- [ ] Get Channel Teams
- [ ] Get Teams
- [x] Get Users
- [x] Update User
- [x] Get Users Follows
- [ ] Create User Follows
- [ ] Delete User Follows
- [x] Get User Block List
- [x] Block User
- [x] UnBlock User
- [ ] Get User Extensions
- [ ] Get User Active Extensions
- [ ] Update User Extensions
- [x] Get Videos
- [x] Delete Videos
- [x] Get Webhook Subscriptions
- [ ] Create Extension Secret
- [ ] Get Extension Secret
- [ ] Revoke Extension Secrets
- [ ] Get Live Channels with Extension Activated
- [ ] Set Extension Required Configuration
- [ ] Set Extension Configuration Segment
- [ ] Get Extension Channel Configuration
- [ ] Get Extension Configuration Segment
- [ ] Send Extension PubSub Message
- [ ] Send Extension Chat Message

## Quick Usage Example

This is a quick example of how to get users. Note that you don't need to provide both a list of ids and logins,
one or the other will suffice.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetUsers(&helix.UsersParams{
    IDs:    []string{"26301881", "18074328"},
    Logins: []string{"summit1g", "lirik"},
})
if err != nil {
    // handle error
}

fmt.Printf("Status code: %d\n", resp.StatusCode)
fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())

for _, user := range resp.Data.Users {
    fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
}
```

Output:

```txt
Status code: 200
Rate limit: 30
Rate limit remaining: 29
Rate limit reset: 1517695315

ID: 26301881 Name: sodapoppin
ID: 18074328 Name: destiny
ID: 26490481 Name: summit1g
ID: 23161357 Name: lirik
```

## Contributions

PRs are very much welcome. Where possible, please write tests for any code that is introduced by your PRs.

## License

This package is distributed under the terms of the [MIT](License) License.
