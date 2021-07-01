# Extensions Documentation

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
