package main

import (
	"emaildrips"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	taskQueueName := "email_drips"

	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	// Create worker and register Workflow and Activity
	w := worker.New(temporalClient, taskQueueName, worker.Options{})
	w.RegisterWorkflow(emaildrips.UserSubscriptionWorkflow)
	w.RegisterActivity(emaildrips.SendContentEmail)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
