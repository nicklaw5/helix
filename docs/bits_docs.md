# Bits Documentation

## Get Bits Leaderboard

This is an example of how to get the last 20 top bits contributers over the past week.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetBitsLeaderboard(&helix.BitsLeaderboardParams{
    Count:  20,
    Period: "week",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
