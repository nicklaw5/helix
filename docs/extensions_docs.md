# Extensions Documentation

## Extension Helix Requests

### Generate PUBSUB JWT Permissions
> relevant PUBSUB permission must be passed to the 'ExtensionCreateClaims()' func, in order to correctly publish a pubsub message of a particular type

Broadcast pubsub type
```go
client.FormBroadcastSendPubSubPermissions()
```

Global pubsub type
```go
perms := client.FormGlobalSendPubSubPermissions()
```

Whisper User type
```go
client.FormWhisperSendPubSubPermissions(userId)
```

### JWT ROLES
> Note:- Currently only the 'external' role is supported by helix endpoints


### EBS JWT
this is used to set the correct header for any Extension helix requests

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
    ExtensionOpts: helix.ExtensionOptions{
        OwnerUserID: os.Getenv(""),
        Secret: os.Getenv(""),
        ConfigurationVersion:  os.Getenv(""),
        Version: os.Getenv(""),
    },
})


// see docs below to see what permissions and roles you can pass 
claims, err := client.ExtensionCreateClaims(broadcasterID, ExternalRole, client.FormBroadcastSendPubSubPermissions(), 0)
if err != nil {
    // handle err
}

jwt,err  := client.ExtensionJWTSign(claims)
if err != nil {
    // handle err
}

// set this before doing extension endpoint requests
client.SetExtensionSignedJWTToken(jwt)
```
## Get Extension Configuration Segments

```go

client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
    ExtensionOpts: helix.ExtensionOptions{
        OwnerUserID: os.Getenv("EXT_OWNER_ID"),
        Secret: os.Getenv("EXT_SECRET"),
        ConfigurationVersion:  os.Getenv("EXT_CFG_VERSION"),
        Version: os.Getenv("EXT_VERSION"),
    },
})
if err != nil {
    // handle error
}

claims, err := client.ExtensionCreateClaims(broadcasterID, ExternalRole, FormBroadcastSendPubSubPermissions(), 0)
if err != nil {
    // handle error
}

// set the JWT token to be used as in the Auth bearer header
jwt := client.ExtensionJWTSign(claims)
client.SetExtensionSignedJWTToken(jwt)

params := helix.ExtensionGetConfigurationParams{
    ExtensionID: "some-extension-id",                              // Required
    Segments:          []helix.ExtensionSegmentType{helix.GlobalSegment}, // Optional
}
resp, err := client.GetExtensionConfigurationSegment
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Extension Transactions

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:        "your-client-id",
    UserAccessToken: "your-user-access-token",
})
if err != nil {
    // handle error
}

params := helix.ExtensionTransactionsParams{
    ExtensionID: "some-extension-id",                              // Required
    ID:          []string{"74c52265-e214-48a6-91b9-23b6014e8041"}, // Optional
    First:       1,                                                // Optional
}

resp, err := client.GetExtensionTransactions(&params)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
