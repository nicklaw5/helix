# Clips Documentation

## Get Clip

This is an example of how to get a single clip.

```go
client, err := helix.NewClient("your-client-id", nil)
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
