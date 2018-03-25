# Streams Documentation

## Get Streams

This is an example of how to get streams. Here we are requesting the first two streams from the English language.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetStreams(&helix.StreamsParams{
    First: 2,
    Language: "en",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Streams Metadata

This is an example of how to get streams metadata. Here we are requesting the first two streams from Hearthstone.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, _ := client.GetStreamsMetadata(&helix.StreamsMetadataParams{
    First:   2,
    GameIDs: []string{"138585"}, // Hearthstone
})
fmt.Printf("%+v\n", resp)
```
