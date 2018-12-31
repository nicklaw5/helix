# Clips Documentation

## Get Clips

### Example 1 - Individual Clips

This is an example of how to get a multiples clip via their IDs.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetClips(&helix.ClipsParams{
    IDs: []string{"EncouragingPluckySlothSSSsss", "PatientBlindingChamoisSmoocherZ"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Example 2 - Broadcaster Clips

This is an example of how to get multiple clips from a single broadcaster.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetClips(&helix.ClipsParams{
    BroadcasterID: "26490481", // summit1g
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Example 3 - Game Clips

This is an example of how to get multiple clips from a single game.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetClips(&helix.ClipsParams{
    GameID: "490377", // Sea of Thieves
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
