# Clips Documentation

## Get Clip

This is an example of how to get a single clip.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetClips(&helix.ClipsParams{
    IDs: []string{"EncouragingPluckySlothSSSsss"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
