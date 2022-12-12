package main

import (
	"context"
	"emaildrips"
	"fmt"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var workflowClient client.Client

func main() {
	var err error

	workflowClient, err = client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	w := worker.New(workflowClient, "email_drips", worker.Options{})
	w.RegisterWorkflow(emaildrips.UserSubscriptionWorkflow)
	w.RegisterActivity(emaildrips.WelcomeEmail)
	w.RegisterActivity(emaildrips.TipsEmail)
	w.RegisterActivity(emaildrips.UnsubscribeEmail)

	go func() {
		err = w.Start()
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	}()

	fmt.Println("Starting dummy server...")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	_ = http.ListenAndServe(":4000", nil)

}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Sign up</h1>")
	_, _ = fmt.Fprint(w, "<form method='post' action='subscribe'><input required name='email' type='email'><input type='submit' value='Subscribe'>")
}

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

	workflowOptions := client.StartWorkflowOptions{
		ID:        "email_drip_" + email,
		TaskQueue: "email_drips",
	}

	_, err = workflowClient.ExecuteWorkflow(context.Background(), workflowOptions, emaildrips.UserSubscriptionWorkflow, email)

	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Couldn't sign up</h1>")
		log.Print(err)
	} else {
		_, _ = fmt.Fprint(w, "<h1>Signed up!</h1>")
	}

}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		// http.ServeFile(w, r, "form.html")
		_, _ = fmt.Fprint(w, "<h1>Unsubscribe</h1>")
		_, _ = fmt.Fprint(w, "<form method='post' action='unsubscribe'><input required name='email' type='email'><input type='submit' value='Unsubscribe'>")
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

		workflowID := "subscription_" + email

		err = workflowClient.CancelWorkflow(context.Background(), workflowID, "")

		if err != nil {
			_, _ = fmt.Fprint(w, "<h1>Couldn't unsubscribe you</h1>")
			log.Fatalln("Unable to cancel Workflow Execution", err)
		} else {
			_, _ = fmt.Fprint(w, "<h1>Unsubscribed you from our emails. Sorry to see you go.</h1>")
			log.Println("Workflow Execution cancelled", "WorkflowID", workflowID)
		}
	}
}
