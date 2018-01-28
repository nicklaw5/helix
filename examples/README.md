# Examples

Here you'll find a number of examples and use cases that can assist you for getting started with [helix](https://github.com/nicklaw5/helix).

## Clips

Get a single clip.

```go
twitch, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

clip, err := twitch.GetClip("clip-id")
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", clip)
```

## Users

Get users *without* an authorization access token:

```go
twitch, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

users, err := twitch.GetUsers(&helix.UsersRequest{
    IDs:    []string{"1", "2"},
    Logins: []string{"summit1g", "lirik"},
})

if err != nil {
    // handle error
}

fmt.Printf("%+v\n", users)
```

Get users *with* an authorization access token:

```go
twitch, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

twitch.SetAccessToken("your-access-token")

users, err := twitch.GetUsers(&helix.UsersRequest{
    Logins: []string{"summit1g"},
})

if err != nil {
    // handle error
}

fmt.Printf("%+v\n", users)
```

## Other

### User-Agent Header

It's entirely possible that you may want to set or change the *User-Agent* header value that is sent with each request. You can do so with the `SetUserAgent()` method before sending a request. For example:

```go
twitch, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

twitch.SetUserAgent("my-user-agent-value")

users, err := twitch.GetUsers(&helix.UsersRequest{
    Logins: []string{"summit1g"},
})

if err != nil {
    // handle error
}

fmt.Printf("%+v\n", users)
```
