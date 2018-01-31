# helix

A Twitch Helix API client written in Go. If you are looking for a client for Twitch's Kraken API, see [kraken](https://github.com/nicklaw5/kraken).

## Package Status

This project is a work in progress. Below is a list of currently supported endpoints. Happy for others to contribute.

## Supported Endpoints

- [x] GET /clips
- [ ] POST /clips
- [ ] POST /entitlements/upload
- [ ] GET /games
- [ ] GET /games/top
- [ ] GET /streams
- [ ] GET /streams/metadata
- [x] GET /users
- [ ] GET /users/follows
- [ ] PUT /users
- [ ] GET /videos

## Getting Started

It's recommended that you use a dependency management tool such as [Dep](https://github.com/golang/dep). If you are using Dep you can import helix by running:

```bash
dep ensure -add github.com/nicklaw5/helix
```

Or you can simply import using the Go toolchain:

```bash
go get -u github.com/nicklaw5/helix
```

## Usage

This is a quick example of how to get users. Note that you don't need to provide both a list of ids and logins, one or the other will suffice.

```go
twitch, err := helix.NewClient("your-client-id", nil)
if err != nil {
    // handle error
}

users := twitch.GetUsers(&helix.UsersRequest{
    IDs: []string{"1", "2"},
    Logins: []string{"summit1g", "lirik"},
})

fmt.Printf("%+v\n", users)
```

## Responses & Rate Limits

It is common for a Twitch API request to simply fail sometimes. Occasionally a request gets hung up and eventually fails with a 500 internal server error. It's also possible that an invalid request was sent and Twitch responded with an error. To assist in circumstances such as these, the HTTP status code is returned with each API request, along with any error that may been encountered. For example, notice below that the `UsersResponse` struct, which is returned with the `GetUsers()` method, includes fields from the `ResponseCommon` struct.

```go
type UsersResponse struct {
    ResponseCommon
    Data []User `json:"data"`
}

type ResponseCommon struct {
    Error              string `json:"error"`
    ErrorStatus        int    `json:"status"`
    ErrorMessage       string `json:"message"`
    RatelimitLimit     int
    RatelimitRemaining int
    RatelimitReset     int
    StatusCode         int
}
```

Also note from above that the `ResponseCommon` struct includes the rate limit header results returned with each request. This package makes no attempt to manage the sending of request based on these rate limit values. That is something your application will need to concur on it's own.

## Examples

See the [examples](examples) page for other use cases.

## License

This package is distributed under the terms of the [MIT](License) License.
