# Raid Documentation

## Start Raid

This is an example of how to start a raid.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.StartRaid(&helix.StartRaidParams{
    FromBroadcasterID: "22484632",
  	ToBroadcasterID: "71092938" 
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Cancel Raid

This is an example of how to cancel a raid.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := resp, err := client.CancelRaid(&helix.StartRaidParams{
    BroadcasterID: "22484632", 
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```