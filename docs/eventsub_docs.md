# EventSub Documentation

EventSub implements the Twitch EventSub notifications. It should cover everything found at https://dev.twitch.tv/docs/eventsub

In contrary to webhooks, eventsub subscriptions do not expire based on time. They will expire when the token used is revoked or if the notification response failure rate is to high.

## Get EventSub Subscriptions

This is an example of how to get eventsub subscriptions.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    AppAccessToken: "your-app-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
    Status: helix.EventSubStatusEnabled, // This is optional.
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Create EventSub Subscription

To create a subscription call CreateEventSubSubscription with a pointer to a subscription. As of writing, Version should always be "1" except for the Channel moderator add / remove events which are still in beta and therefore you need to use Version "beta".
Within the Transport the only supported Method currently is "webhook". Callback needs to be a https link on port 443. With the secret you can verify if notifications came from twitch. See (#verify-eventSub-notification)

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    AppAccessToken: "your-app-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.CreateEventSubSubscription(&helix.EventSubSubscription{
    Type: helix.EventSubTypeChannelFollow,
    Version: "1",
    Condition: helix.EventSubCondition{
        BroadcasterUserID: "1337",
    },
    Transport: helix.EventSubTransport{
        Method: "webhook",
        Callback: "https://example.com/follow",
        Secret: "s3cre7w0rd",
    },
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Delete EventSub Subscription

To delete a subscription you need to call RemoveEventSubSubscription with the subscription id as parameter.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
    AppAccessToken: "your-app-access-token",
})
if err != nil {
    // handle error
}

resp, err := client.RemoveEventSubSubscription("subscription-id")
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Example for handling a notification

```go

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func eventsubFollow(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Println(err)
        return
    }
    defer r.Body.Close()
    // verify that the notification came from twitch using the secret.
    if !helix.VerifyEventSubNotification("s3cre7w0rd", r.Header, string(body)) {
        log.Println("no valid signature on subscription")
        return
    } else {
        log.Println("verified signature for subscription")
    }
    var vals eventSubNotification
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
    if err != nil {
        log.Println(err)
        return
    }
    // if there's a challenge in the request, respond with only the challenge to verify your eventsub.
    if vals.Challenge != "" {
        w.Write([]byte(vals.Challenge))
        return
    }
    var followEvent helix.EventSubChannelFollowEvent
    err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&followEvent)

    log.Printf("got follow webhook: %s follows %s\n", followEvent.UserName, followEvent.BroadcasterUserName)
    w.WriteHeader(200)
    w.Write([]byte("ok"))
}
```
