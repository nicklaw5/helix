# Logging Documentation

## Adding a custom logger

You can add a custom logger to the client that will be used to print outgoing requests and incoming responses. Debug
mode must be set to `true` to enable logging.

This is an example of how to enable logging and log to a file:

```go
// Create a logger, which outputs to file.
f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
if err != nil {
    log.Fatalf("error opening file: %v", err)
}
defer f.Close()

fileLogger := &log.Logger{}
fileLogger.SetOutput(f)

client, err := helix.NewClient(&helix.Options{
    logger: fileLogger,	
    debug: true,	
})
if err != nil {
    // handle error
}
```
