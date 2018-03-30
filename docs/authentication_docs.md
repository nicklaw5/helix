# Authentication Documentation

## Get Authorization URL

The below `GetAuthorizationURL` method returns a URL, based on your Client ID, redirect URI, and scopes, that can be used to authenticate users with their Twitch accounts. After the user authorizes your application, they will be redirected to the provided redirect URI with an authorization code that can be used to generate access tokens for API consumption. See the Twitch [authentication docs](https://dev.twitch.tv/docs/authentication) for more information.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:    "your-client-id",
    RedirectURI: "https://example.com/auth/callback",
    Scopes:      []string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
})
if err != nil {
    // handle error
}

url := client.GetAuthorizationURL("your-state", true)

fmt.Printf("%s\n", url)
```

## Get App Access Token

Here's an example of how to create an app access token:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})
if err != nil {
    // handle error
}

resp, err := client.GetAppAccessToken()
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get User Access Token

After obtaining an authentication code, you can submit a request for a user access token which can then be used to submit API requests on behalf of a user. Here's an example of how to create a user access token:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURI:  "https://example.com/auth/callback",
})
if err != nil {
    // handle error
}

code := "your-authentication-code"

resp, err := client.GetUserAccessToken(code)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Refresh User Access Token

You can refresh a user access token in the following manner:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})
if err != nil {
    // handle error
}

refreshToken := "your-refresh-token"

resp, err := client.RefreshUserAccessToken(refreshToken)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Revoke User Access Token

You can revoke a user access token in the following manner:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:     "your-client-id",
})
if err != nil {
    // handle error
}

accessToken := "your-access-token-to-revoke"

resp, err := client.RevokeUserAccessToken(accessToken)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
