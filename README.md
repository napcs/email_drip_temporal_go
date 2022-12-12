# Email drip campaign with Temporal

Demo app that shows how to create a subscribe/unsubscribe feature that uses a 
long-running Temporal Workflow to send weekly email messages.

## Running

Start a Temporal server.

Then run the following command to start the server:

```
go run server/main.go
```

Visit `http://localhost:3000` and add your email to subscribe.

Visit `http://localhost:3000/unsubscribe` and add your email to unsubscribe.



