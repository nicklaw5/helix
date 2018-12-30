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

## Create Clip

This is an example of how to create a clip:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-acceess-token",
})
if err != nil {
    // handle error
}

resp, err := client.CreateClip(&helix.CreateClipParams{
    BroadcasterID: "26490481", // summit1g
    HasDelay: true, // optional, defaults to false
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
