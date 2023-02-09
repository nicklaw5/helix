# Vips Documentation

## Get VIPs

Gets a list of the broadcaster’s VIPs.

To use this function you need a user access token with the `channel:read:vips` scope.
`BroadcasterID` is required and need to be the same as the user id of the user access token.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetChannelVips(&helix.GetChannelVipsParams{
    BroadcasterID: "54946241",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Add Channel VIP

Adds the specified user as a VIP in the broadcaster’s channel.

To use this function you need a user access token with the `channel:manage:vips` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

//Add Vip
resp, err := client.AddChannelVip(&helix.AddChannelVipParams{
    UserID:        "23981723",
    BroadcasterID: "54946241",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Remove Channel VIP

Removes the specified user as a VIP in the broadcaster’s channel

To use this function you need a user access token with the `channel:manage:vips` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.RemoveChannelVip(&helix.RemoveChannelVipParams{
    UserID:        "23981723",
    BroadcasterID: "54946241",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

