# Docs

## Endpoints

Follow the links below to their respective documentation:

- [Clips](clips_docs.md)
- [Games](games_docs.md)
- [Users](users_docs.md)
- [Videos](videos_docs.md)

## Getting Started

It's recommended that you use a dependency management tool such as [Dep](https://github.com/golang/dep). If you are using Dep you can import helix by running:

```bash
dep ensure -add github.com/nicklaw5/helix
```

Or you can simply import using the Go toolchain:

```bash
go get -u github.com/nicklaw5/helix
```

## Creating A New API Client

The only requirement for creating a new API client is your Twitch Client-ID. See the [Twitch authentication docs](https://dev.twitch.tv/docs/authentication) on how to obtain a Client-ID. Once you have a Client-ID, to create a new client simply pass the Client-ID through as the first argument of the `NewClient` function, like so:

```go
client := helix.NewClient("your-client-id", nil)
```

If you'd like to pass in your own `http.Client`, you can do so by passing it through as an option when creating a new client, like so:

```go
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
    },
    Timeout: 10 * time.Second,
}

client := helix.NewClient("your-client-id", &helix.Options{
    HTTPClient: httpClient,
})
```

If no custom `http.Client` is provided, `http.DefaultClient` is used by default.

## Responses

It is common for a Twitch API request to simply fail sometimes. Occasionally a request gets hung up and eventually fails with a 500 internal server error. It's also possible that an invalid request was sent and Twitch responded with an error. To assist in circumstances such as these, the HTTP status code is returned with each API request, along with any error that may been encountered. For example, notice below that the `UsersResponse` struct, which is returned with the `GetUsers()` method, includes fields from the `ResponseCommon` struct.

```go
type UsersResponse struct {
    ResponseCommon
    Data ManyUsers
}

type ManyUsers struct {
    Users []User `json:"data"`
}

type ResponseCommon struct {
    StatusCode   int
    Error        string `json:"error"`
    ErrorStatus  int    `json:"status"`
    ErrorMessage string `json:"message"`
    RateLimit    RateLimit
}

type RateLimit struct {
    Limit     int
    Remaining int
    Reset     int64
}
```

Also note from above that the `ResponseCommon` struct includes the rate limit header results returned with each request. See below on how you may want to limit your requests based on the rate limit results returned with each request.

## Request Rate Limiting

Twitch enforces strict request rate limits for their API. See [their documentation](https://dev.twitch.tv/docs/api#rate-limits) for the specific rate limit values. At the time of writing this, requests are limited to 30 queries per minute (if a Bearer token is not provided) or 120 queries per minute (if a Bearer token is provided).

This package allows users to provide a rate limit callback that will be executed just before a request is sent. That way you can provide some sort of functionality for limiting the requests sent and prevent spamming Twitch with requests.

The below snippet provides an example of how you might structure your rate limit callback to approach limiting requests. In this example, once we've reached our rate limit, we'll simply wait for the limit to pass before sending the next request.

```go
func rateLimitCallback(lastResponse *helix.Response) error {
    if lastResponse.RateLimit.Remaining > 0 {
        return nil
    }

    currentTime := time.Now().Unix()

    if currentTime < lastResponse.RateLimit.Reset {
        timeDiff := time.Duration(lastResponse.RateLimit.Reset - currentTime)
        if timeDiff > 0 {
            fmt.Printf("Waiting on rate limit to pass before sending next request (%d seconds)\n", timeDiff)
            time.Sleep(timeDiff * time.Second)
        }
    }

    return nil
}

client := helix.NewClient("your-client-id", &helix.Options{
    RateLimitFunc: rateLimitCallback,
})
```

If a `RateLimitFunc` is provided, the client will re-attempt to send a failed request if said request received a 429 (Too Many Requests) response. Before retrying the request, the `RateLimitFunc` will be applied. This functionality is enabled by default but can be disabled by setting the `RetryRateLimitedRequests` option to false.

## Access Token Header

Some endpoints require that you have a valid access token in order to fulfill the request.

In order to set the access token for a request, use the `SetAccessToken` method. For example:

```go
client := helix.NewClient("your-client-id", nil)
client.SetAccessToken("your-access-token")

// send API request...
```

Note that any subsequent API requests will utilize this same access token. So it is necessary to unset the access token when you are finished with it. To do so, simply pass an empty string to the `SetAccessToken` method.

## User-Agent Header

It's entirely possible that you may want to set or change the *User-Agent* header value that is sent with each request. You can do so by it through as option when creating a new client, like so:

with the `SetUserAgent()` method before sending a request. For example:

```go
client := helix.NewClient("your-client-id", &helix.Options{
    UserAgent: "your-user-agent-value",
})

// send API request...
```

Alternatively, you can set by calling the `SetUserAgent()` method before sending a request. For example:

```go
client := helix.NewClient("your-client-id", nil)
client.SetUserAgent("your-user-agent-value")

// send API request...
```
