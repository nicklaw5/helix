# Drops Documentation

## Get Drops

### Example 1 - Get drops entitlements for a game

This is an example of how to get entitlements for all users of a game.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetDropsEntitlements(&helix.GetDropEntitlementsParams{
	GameID: "your-game-id",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Example 2 - Get drops entitlements for a user

This is an example of how to get entitlements for a specific user.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetDropsEntitlements(&helix.GetDropEntitlementsParams{
	UserID: "your-game-id",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Example 3 - Get drops entitlements next page of results

This is an example of how to get multiple pages of entitlement results.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

// First page 
resp, err := client.GetDropsEntitlements(&helix.GetDropEntitlementsParams{
	GameID: "your-game-id",
    After: "",
    First: 50,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)

// Next page
resp, err = client.GetDropsEntitlements(&helix.GetDropEntitlementsParams{
	GameID: "your-game-id",
    After: resp.Data.Pagination.Cursor,
    First: 50,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

### Example 4 - Get drops entitlements filtered by fulfillment status

This is an example of how to get entitlements for a specific fulfillment status.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetDropsEntitlements(&helix.GetDropEntitlementsParams{
	UserID: "your-game-id",
    FulfillmentStatus: "CLAIMED"
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Update Drops Entitlements

### Example 1 - Update drops entitlemetns

This is an example of how to update drops entitlements status.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.UpdateDropsEntitlements(&helix.UpdateDropsEntitlementsParams{
	EntitlementIDs: ["entitlement-id-1", "entitlement-id-2"],
	FulfillmentStatus: "FULFILLED",
})
if err != nil {
    // handle error
}

// Check which of the requested ids were updated
for _, es := range resp.Data.EntitlementSets {
    if es.Status == "SUCCESS" {
        fmt.Printf("%+v\n", es.IDs)
    }
}
```
