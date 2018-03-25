# Games Documentation

## Get Games

This is an example of how to get games.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetGames(&helix.GamesParams{
    Names: []string{"Sea of Thieves", "Fortnite"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Top Games

This is an example of how to get top games.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetTopGames(&helix.TopGamesParams{
    First: 20,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
