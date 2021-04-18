# Channels Documentation

## Search Channels

This is an example of how to search channels. Here we are requesting the first two streams from the English language. SearchChannels returns live as well as offline channels.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.SearchChannels(&helix.SearchChannelsParams{
    First: 2,
    Language: []string{"en"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Channel Information

This is an example of how to get channel informations.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetChannelInformation(&helix.GetChannelInformationParams{
    BroadcasterID: "123456",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Modify Channel Information

This is an example of how to modify channel informations.
The `Delay` param is a Twitch Partner feature.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetChannelInformation(&helix.GetChannelInformationParams{
    BroadcasterID       : "123456",
    GameID              : "456789",
    BroadcasterLanguage : "en",
    Title               : "Your stream title",
    Delay               : 0,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
