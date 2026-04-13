# Upgrading Guide

This document describes breaking changes and the steps required to upgrade between major versions of this package.

---

## Upgrading from v2 to v3

### Import path

Update your import path from `github.com/nicklaw5/helix/v2` to `github.com/nicklaw5/helix/v3`:

```go
// Before
import "github.com/nicklaw5/helix/v2"

// After
import "github.com/nicklaw5/helix/v3"
```

---

### All endpoint methods now return a non-nil error for Twitch API errors

**Before (v2):**

Every endpoint method always returned `nil` as the Go `error` value, even when the Twitch API
responded with an HTTP 4xx or 5xx error. Callers had to inspect `resp.ErrorStatus` manually:

```go
resp, err := client.GetUsers(&helix.UsersParams{Logins: []string{"twitchdev"}})
if err != nil {
    // Only transport / network errors reached here.
    log.Fatal(err)
}
if resp.ErrorStatus != 0 {
    // Twitch returned an API-level error — had to be checked explicitly.
    log.Fatalf("twitch error %d: %s", resp.ErrorStatus, resp.ErrorMessage)
}
```

**After (v3):**

When the Twitch API responds with a non-zero `status` field in the response body (i.e. any API-level
error), the endpoint method now returns a **non-nil** Go `error` that wraps the exported sentinel
`helix.ErrTwitchResponseError`. The `ResponseCommon` fields (`ErrorStatus`, `Error`,
`ErrorMessage`) are still populated and accessible through the response struct.

```go
resp, err := client.GetUsers(&helix.UsersParams{Logins: []string{"twitchdev"}})
if err != nil {
    if errors.Is(err, helix.ErrTwitchResponseError) {
        // Twitch returned an API-level error.
        // resp.ErrorStatus / resp.Error / resp.ErrorMessage are still accessible.
        log.Fatalf("twitch error %d: %s", resp.ErrorStatus, resp.ErrorMessage)
    }
    // Some other transport / network error.
    log.Fatal(err)
}
```

**What you need to change:**

1. Any code that checked `resp.ErrorStatus != 0` to detect Twitch errors should now check `err != nil` (or `errors.Is(err, helix.ErrTwitchResponseError)`) instead.
2. Any code that assumed `err == nil` after a successful HTTP round-trip must now handle the possibility of a Twitch API-level error being surfaced as a Go error.
3. If you only care whether the call succeeded at the Twitch level and still want to read the response body, check `errors.Is(err, helix.ErrTwitchResponseError)` and then access `resp.ErrorStatus` / `resp.ErrorMessage` as before.

**Sentinel error:**

```go
// ErrTwitchResponseError is returned (wrapped) when the Twitch API responds
// with a non-zero status in the response body.
var ErrTwitchResponseError = errors.New("twitch returned an error")
```

The wrapped error message includes the HTTP status code, error string, and message from Twitch, for example:

```
twitch returned an error: Bad Request (status 400, message "Missing required parameter \"broadcaster_id\"")
```
