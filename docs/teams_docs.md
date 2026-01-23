# Teams Documentation

## Get Teams

Gets information about the specified Twitch team. You can specify the team by its name or ID, but not both.

### Get Team by Name

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetTeams(&helix.GetTeamsParams{
    Name: "weightedblanket",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Get Team by ID

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetTeams(&helix.GetTeamsParams{
    ID: "1234567",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
