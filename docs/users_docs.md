# Users Documentation

## Get Users

This is an example of how to get users. Note that you don't need to provide both a list of ids and logins, one or the other will suffice.

```go
client, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

resp, err := client.GetUsers(&helix.UsersParams{
    IDs:    []string{"26301881", "18074328"},
    Logins: []string{"summit1g", "lirik"},
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
