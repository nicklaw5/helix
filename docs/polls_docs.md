# Polls Documentation

## Get Polls

This is an example of how to get polls.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.GetPolls(&helix.PollsParams{
    BroadcasterID: "145328278",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Create Poll

This is an example of how to create a poll.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.CreatePoll(&helix.CreatePollParams{
    BroadcasterID: "145328278",
    Title: "Test",
    Choices: []helix.PollChoiceParam{
        helix.PollChoiceParam{ Title: "choice 1" },
        helix.PollChoiceParam{ Title: "choice 2" },
    },
    Duration: 30,
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## End Poll

This is an example of how to end a poll.

```go
client, err := helix.NewClient(&helix.Options{
    ClientID: "your-client-id",
})
if err != nil {
    // handle error
}

resp, err := client.EndPoll(&helix.EndPollParams{
    BroadcasterID: "145328278",
    ID: "25b14b42-d4d8-4756-86ce-842bf76f82a0",
    Status: "TERMINATED",
})
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
