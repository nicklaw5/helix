# Chat Documentation

## Get Channel Chat Badges

This is an example of how to get channel chat badges

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetChannelChatBadges(&helix.GetChatBadgeParams{
    BroadcasterID: "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Global Chat Badges

This is an example of how to get global chat badges

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetGlobalChatBadges()
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Channel Emotes

This is an example of how to get channel emotes

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetChannelEmotes(&helix.GetChannelEmotesParams{
    BroadcasterID: "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Global Emotes

This is an example of how to get global emotes

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetGlobalEmotes()
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Emote Sets

This is an example of how to get a set of emotes

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetEmoteSets(&helix.GetEmoteSetsParams{
    EmoteSetIDs: []string{"300678379"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Chat Settings

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    AppAccessToken: "your-app-user-token",
    // Optionally, a moderator's user token with the `moderator:read:chat_settings` scope can be specified to read some more settings
})
if err != nil {
    // handle error
}

resp, err := client.GetChatSettings(&helix.GetChatSettingsParams{
    BroadcasterID: "22484632",
    // ModeratorID should be specified matching the UserAccessToken if you want the extended information
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get User Chat Color
Gets the color used for the user’s name in chat.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetUserChatColor(&GetUserChatColorParams{
    UserID: "22484632",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Update User Chat Color
Updates the color used for the user’s name in chat.
```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.UpdateUserChatColor(&UpdateUserChatColorParams{
    UserID: "22484632",
    Color: "blue"
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```