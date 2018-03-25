# Authentication Documentation

## Get Authorization URL

The below `GetAuthorizationURL` method returns a URL, based on your Client ID, redirect URI, and scopes, that can be used to authenticate users with their Twitch accounts. After the user authorizes your application, they will be redirected to the provided redirect URI with an authorization code that can be used to generate access tokens for API consumption. See the Twitch [authentication docs](https://dev.twitch.tv/docs/authentication) for more information.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    RedirectURI: "https://example.com/auth/callback",
    Scopes:      []string{"analytics:read:games", "bits:read", "clips:edit", "user:edit", "user:read:email"},
})
if err != nil {
    // handle error
}

url := client.GetAuthorizationURL("your-state", true)

fmt.Printf("%s\n", url)
```
