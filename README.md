# Email drip campaign with Temporal

Demo app that shows how to create a subscribe/unsubscribe feature that uses a 
long-running Temporal Workflow to send weekly email messages.

## Running

Start a Temporal server. 

Then run the following command to start the web server:

```
go run server/main.go
```

In another window, run the following command to start a Worker process:

```
go run server/main.go
```

Visit `http://localhost:4000` and add your email to subscribe.

Visit `http://localhost:4000/unsubscribe` and add your email to unsubscribe.


## How it works.

The repository contains the following files:

```
├── README.md         <- This file.
├── go.mod
├── go.sum
├── mails             <- Emails for the campaign
│   ├── 1.md
│   ├── 2.md
│   ├── 3.md
│   ├── goodbye.md    <- The unsubscribe message
│   └── welcome.md    <- A welcome message
├── server
│   └── main.go       <- The web server
├── worker
│   └── main.go       <- The Temporal worker
└── subscribe.go      <- The application code including Workflows and Activities.
```


The `server/main.go` file contains a basic Go web server that exposes three endpoints:
* `/`: Displays the signup form
* `/subscribe`: Processes the signup request and executes the Temporal Workflow to handle subscriptions.
* `/unsubscribe`: Supports `GET` and `POST`:
  * `GET /unsubscribe`: Displays the unsubscribe form
  * `POST /unsubscribe`: processes the unsubscribe request for the form by sending a cancellation request to the Workflow.

The `worker/main.go` file contains the Temporal Worker that executes Workflows and Activities.

The `subscribe.go` file contains the Temporal Workflow and Activity for subscriptions.
* `UserSubscriptionWorkflow` accepts a `Subscription` which is made up of the email address and the email campaign to send. The campaign is a struct that specifies the email messages that make up the campaign, including the welcome message, the unsubscribe message, and the messages that make up the campaign's contents.
* `SendContentEmail` is an Activity that accepts an `EmailInfo` struct which contains the message to send and email address.
* The `UserSubscriptionWorkflow` sends the welcome message and the first message of the campaign, and then sleeps for the duration specified by the `duration` value in the `UserSubscriptionWorkflow` function.


