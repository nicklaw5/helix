# Clips Documentation

## Get Clip

This is an example of how to get a single clip.

```go
client := helix.NewClient("your-client-id", nil)

resp, err := client.GetClips(&helix.ClipsParams{
    IDs: []string{"EncouragingPluckySlothSSSsss"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
