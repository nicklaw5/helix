# Categories Documentation

## Search Categories

This is an example of how to search categories. Here we are requesting the first two categories that match the query string `pokemon`. 

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}
resp, err := client.SearchCategories(&helix.SearchCategoriesParams{
    First: 2,
    Query: "pokemon",
})
if err != nil {
    // handle error
}
fmt.Printf("%+v\n", resp)
```