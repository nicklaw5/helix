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
