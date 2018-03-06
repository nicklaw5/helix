# Streams Documentation

## Get Streams

This is an example of how to get streams. Here we are requesting the first two streams from the English language.

```go
client := helix.NewClient("your-client-id", nil)

resp, err := client.GetStreams(&helix.StreamsParams{
    First: 2,
    Language: "en",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
