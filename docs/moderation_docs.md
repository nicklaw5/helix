# Moderation Documentation

## Get Banned Users

To use this function you need a user access token with the `moderation:read` scope.
`BroadcasterID` is required and need to be the same as the user id of the user access token.

This is an example of how to get banned users in a channel.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetBannedUsers(&helix.BannedUsersParams{
    BroadcasterID: "54946241",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Manage Held AutoMod Messages

This is an example of how to manage held automod message in a channel.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.ModerateHeldMessage(&helix.HeldMessageModerationParams{
    UserID : "145328278",
    MsgID  : "19fe2618-df5f-45d3-a210-aeda6f6c6d9e",
    Action : "ALLOW",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

