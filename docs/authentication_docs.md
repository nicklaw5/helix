# Authentication Documentation

## Get Authorization URL

The below `GetAuthorizationURL` method returns a URL, based on your Client ID, and redirect URI,
that can be used to authenticate users with their Twitch accounts. After the user authorizes your
application, they will be redirected to the provided redirect URI with an authorization code that
can be used to generate access tokens for API consumption. See the Twitch
[authentication docs](https://dev.twitch.tv/docs/authentication) for more information.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:    "your-client-id",
    RedirectURI: "https://example.com/auth/callback",
})
if err != nil {
    // handle error
}

url := client.GetAuthorizationURL(&helix.AuthorizationURLParams{
    ResponseType: "code", // or "token"
    Scopes:       []string{"user:read:email"},
    State:        "some-state",
    ForceVerify:  false,
})

fmt.Printf("%s\n", url)
```

## Get User Access Token

After obtaining an authentication code, you can submit a request for a user access token which can
then be used to submit API requests on behalf of a user. Here's an example of how to create a user
access token:

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

resp, err := client.RequestUserAccessToken(code)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)

// Set the access token on the client
client.SetUserAccessToken(resp.Data.AccessToken)
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
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

userAccessToken := client.GetUserAccessToken()

resp, err := client.RevokeUserAccessToken(userAccessToken)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Validate User Access Token

You can validate an access token and get token details in the following manner:

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

userAccessToken := client.GetUserAccessToken()

isValid, resp, err := client.ValidateToken(userAccessToken)
if err != nil {
    // handle error
}

if isValid {
    fmt.Println("%s access token is valid!", userAccessToken)
}

fmt.Println("%+v", resp)
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

resp, err := client.RequestAppAccessToken([]string{"user:read:email"})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)

// Set the access token on the client
client.SetAppAccessToken(resp.Data.AccessToken)
```
