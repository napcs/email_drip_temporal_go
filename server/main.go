package main

import (
	"context"
	"emaildrips"
	"fmt"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

// global client.
var temporalClient client.Client
var taskQueueName string

func main() {
	port := "4000"
	taskQueueName = "email_drips"

	var err error

	temporalClient, err = client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting the web server on port %s\n", port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	_ = http.ListenAndServe(":"+port, nil)

}

// Index page shows the subscription form.
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Sign up</h1>")
	_, _ = fmt.Fprint(w, "<form method='post' action='subscribe'><input required name='email' type='email'><input type='submit' value='Subscribe'>")
}

// Handle subscriptions from the form.
func subscribeHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
		return
	}

	email := r.PostForm.Get("email")

	if email == "" {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
		return
	}

	// use the email as the id in the workflow. This may leak PII.
	workflowOptions := client.StartWorkflowOptions{
		ID:        "email_drip_" + email,
		TaskQueue: taskQueueName,
	}

	// Define the subscription
	subscription := emaildrips.Subscription{
		EmailAddress: email,
		Campaign: emaildrips.Campaign{
			Name:             "Temporal Tips and Tricks",
			WelcomeEmail:     "../mails/welcome.md",
			UnsubscribeEmail: "../mails/goodbye.md",
			Mails:            []string{"../mails/1.md", "../mails/2.md", "../mails/3.md"},
		},
	}

	// execute the Temporal Workflow to start the subscription.
	_, err = temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, emaildrips.UserSubscriptionWorkflow, subscription)

	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Couldn't sign up</h1>")
		log.Print(err)
	} else {
		_, _ = fmt.Fprint(w, "<h1>Signed up!</h1>")
	}

}

// Handle unsubscribe requests.
func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		// http.ServeFile(w, r, "form.html")
		_, _ = fmt.Fprint(w, "<h1>Unsubscribe</h1><form method='post' action='unsubscribe'><input required name='email' type='email'><input type='submit' value='Unsubscribe'>")

	case "POST":

		err := r.ParseForm()

		if err != nil {
			// in case of any error
			_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
			return
		}

		email := r.PostForm.Get("email")

		if email == "" {
			// in case of any error
			_, _ = fmt.Fprint(w, "<h1>Email is blank</h1>")
			return
		}

		workflowID := "email_drip_" + email

		err = temporalClient.CancelWorkflow(context.Background(), workflowID, "")

		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Couldn't unsubscribe you</h1>")
			log.Fatalln("Unable to cancel Workflow Execution", err)
		} else {
			_, _ = fmt.Fprint(w, "<h1>Unsubscribed you from our emails. Sorry to see you go.</h1>")
			log.Println("Workflow Execution cancelled", "WorkflowID", workflowID)
		}
	}
}
