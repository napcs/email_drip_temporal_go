package emaildrips

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

/* UserSubscriptionWorkflow handles subscribing users */
func UserSubscriptionWorkflow(ctx workflow.Context, email string) (result string, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + email)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// Handle any cleanup, including cancellations..
	defer func() {

		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}

		// Cancellation received, which will be an unsubscribe request.

		newCtx, _ := workflow.NewDisconnectedContext(ctx)

		logger.Info("Sending unsubscribe email to "+email, err)
		err := workflow.ExecuteActivity(newCtx, UnsubscribeEmail, email)
		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()

	logger.Info("Sending tips email to " + email)
	err = workflow.ExecuteActivity(ctx, WelcomeEmail, email).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	}

	workflow.Sleep(ctx, (24*7)*time.Hour)

	logger.Info("Sending tips email to " + email)
	err = workflow.ExecuteActivity(ctx, TipsEmail, email).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send tips email", "Error", err)
	}

	workflow.Sleep(ctx, (24*7)*time.Hour)

	logger.Info("Sending sales email to " + email)
	err = workflow.ExecuteActivity(ctx, SalesEmail, email).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send sales email", "Error", err)
	}

	return "", nil
}

/* WelcomeEmail is the activity that sends the welcome email to the customer. */
func WelcomeEmail(ctx context.Context, email string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending welcome email to " + email)

	// call mailer api here.

	return nil

}

/* TipsEmail is the activity that sends the email with tips and tricks to the customer. */
func TipsEmail(ctx context.Context, email string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending tips email to " + email)

	// call mailer api here.

	return nil

}

/*  SalesEmail is the activity that sends the email that asks to talk to the customer. */
func SalesEmail(ctx context.Context, email string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending tips email to " + email)

	// call mailer api here.

	return nil

}

/* UnsubscribeEmail is sent when someone unsubscribes */
func UnsubscribeEmail(ctx context.Context, email string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending tips email to " + email)

	// call mailer api here.

	return nil

}
