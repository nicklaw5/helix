# Games Documentation

## Get Games

This is an example of how to get games.

```go
client, err := helix.NewClient("your-client-id", nil)
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
