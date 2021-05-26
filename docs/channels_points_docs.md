# Channels Points Documentation

## Create Custom Rewards

This is an example of how to create a custom reward.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.CreateCustomReward(&helix.ChannelCustomRewardsParams{
    BroadcasterID : "145328278",
    Title         : "game analysis 1v1",
    Cost          : 50000,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
