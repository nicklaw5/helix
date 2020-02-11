# User Extensions Documentation

## Get User Extensions

This is an example of how to get user extensions

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetUserExtensions()
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get User Active Extensions

This is an example of how to get active user extensions

Using UserAccessToken:
```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    UserAccessToken: "user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetUserExtensions(nil)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

Using user_id query parameter:
```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetUserExtensions(&helix.UserActiveExtensionsParams{
    UserID: "user-id"
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
