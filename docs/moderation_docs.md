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

## Ban User

This is an example of how to ban or timeout a user.

To use this function you need a user access token with the `moderator:manage:banned_users` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

// Ban user permanently
resp, err := client.BanUser(&helix.BanUserParams{
    BroadcasterID: "54946241",
    ModeratorId:   "14532827",
    Body: helix.BanUserRequestBody{
        UserId: "23981723",
        Reason: "no reason",
    },
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)

// Put user in a timeout
resp, err = client.BanUser(&helix.BanUserParams{
    BroadcasterID: "54946241",
    ModeratorId:   "14532827",
    Body: helix.BanUserRequestBody{
        UserId:   "23981723",
        Duration: 3600,
        Reason:   "no reason",
    },
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Unban User

This is an example of how to unban a user.

To use this function you need a user access token with the `moderator:manage:banned_users` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.UnbanUser(&helix.UnbanUserParams{
    BroadcasterID: "54946241",
    ModeratorID:   "14532827",
    UserID:        "23981723",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Blocked Terms

This is an example of how to get the list of non-private, blocked words or phrases.

To use this function you need a user access token with the `moderator:read:blocked_terms` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetBlockedTerms(&helix.BlockedTermsParams{
    BroadcasterID: "54946241",
    ModeratorID:   "14532827",
    First:         100,                                 // optional
    After:         "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6I", // optional
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Add Blocked Term

This is an example of how to add a word or phrase to the list of blocked terms.

To use this function you need a user access token with the `moderator:manage:blocked_terms` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.AddBlockedTerm(&helix.AddBlockedTermParams{
    BroadcasterID: "54946241",
    ModeratorID:   "14532827",
    Text:          "crac*",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Remove Blocked Term

This is an example of how to remove a word or phrase to the list of blocked terms.

To use this function you need a user access token with the `moderator:manage:blocked_terms` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.RemoveBlockedTerm(&helix.RemoveBlockedTermParams{
    BroadcasterID: "54946241",
    ModeratorID:   "14532827",
    ID:            "c9fc79b8-0f63-4ef7-9d38-efd811e74ac2",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Delete Specific Chat Message

This is an example of how to delete a specific chat message. 

To use this function you need a user access token with the `moderator:manage:chat_messages` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.DeleteChatMessage(&helix.DeleteChatMessageParams{
    BroadcasterID: "54946241",
    ModeratorID:   "145328278",
    MessageID:     "885196de-cb67-427a-baa8-82f9b0fcd05f",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Delete All Chat Messages

This is an example of how to delete all chat message.

To use this function you need a user access token with the `moderator:manage:chat_messages` scope.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.DeleteAllChatMessages(&helix.DeleteAllChatMessagesParams{
    BroadcasterID: "54946241",
    ModeratorID:   "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Manage Held AutoMod Messages

This is an example of how to manage held automod message in a channel.

To use this function you need a user access token with the `moderator:manage:automod` scope.

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

## Get Moderators

To use this function you need a user access token with the `moderation:read` scope.
`BroadcasterID` is required and need to be the same as the user id of the user access token.

This is an example of how to get moderators of a channel.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetModerators(&helix.GetModeratorsParams{
    BroadcasterID: "145328278",
    First: 10
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Moderated Channels

To use this function you need a user access token with the `user:read:moderated_channels` scope.
`UserID` is required and must match the user ID of the user access token.

This is an example of how to get channels the user has moderator privileges in.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetModeratedChannels(&helix.GetModeratedChannelsParams{
    UserID: "154315414",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Warn Chat User

Requires a user access token that includes the moderator:manage:warnings scope. Query parameter moderator_id must match the user_id in the user access token.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.SendModeratorWarnMessage(
    &SendModeratorWarnChatMessageParams{
        BroadcasterID: "1234",
        ModeratorID:   "5678",
        UserID:        "9876",
        Reason:        "Test warning message",
    },
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```