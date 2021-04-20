# Users Documentation

## Get Users

This is an example of how to get users. Note that you don't need to provide both a list of ids and logins, one or the other will suffice.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
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

## Update User

This is an example of how to update a users description:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.UpdateUser(&UpdateUserParams{
  Description: "New description",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Users Follows

This is an example of how to get users follows.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetUsersFollows(&helix.UsersFollowsParams{
    FromID:  "23161357",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get User Block List

This is an example of how to get users blocked

```
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetUsersBlocked(&helix.UsersBlockedParams{
    BroadcasterID: "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Block User

This is an example of how to block user

```
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.BlockUser(&helix.BlockUserParams{
    TargetUserID:  "677636701",
    SourceContext: "chat",
    Reason:        "spam",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Unblock User

This is an example of how to unblock user

```
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.UnblockUser(&helix.UnblockUserParams{
    TargetUserID: "677636701",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
