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

## Update Custom Rewards

This is an example of how to update a custom reward.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.UpdateCustomReward(&helix.UpdateChannelCustomRewardsParams{
    ID            : "6741db51-bc4e-4f0e-b96b-d79eafe227f3",
    BroadcasterID : "145328278",
    Title         : "game analysis 1v1",
    Cost          : 50000,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Delete Custom Rewards

This is an example of how to delete a custom rewards.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.DeleteCustomRewards(&helix.DeleteCustomRewardsParams{
    BroadcasterID : "145328278",
    ID            : "84da6b13-efe1-4a82-91d0-25260aeb6a9b",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Custom Rewards

This is an example of how to get a custom rewards.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetCustomRewards(&helix.GetCustomRewardsParams{
    BroadcasterID : "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Custom Reward Redemption Status

This is an example of how to get the status of a custom reward redemption.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetCustomRewardsRedemptions(&helix.GetCustomRewardsRedemptionsParams{
    BroadcasterID : "145328278",
    RewardID      : "48d3a271-0571-49e0-91db-d9919274bbd3",
    Status        : "UNFULFILLED",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Update Custom Reward Redemption Status

This is an example of how to update the status of a custom reward redemption.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.UpdateChannelCustomRewardsRedemptionStatus(&helix.UpdateChannelCustomRewardsRedemptionStatusParams{
    ID            : "17fa2df1-ad76-4804-bfa5-a40ef63efe63",
    BroadcasterID : "274637212",
    RewardID      : "92af127c-7326-4483-a52b-b0da0be61c01",
    Status        : "FULFILLED", // The only acceptable values are "CANCELED" and "FULFILLED"
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
