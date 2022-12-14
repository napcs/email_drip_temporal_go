package emaildrips

import (
	"bufio"
	"context"
	"errors"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// EmailInfo is the data that the SendContentEmail uses to send the message.
type EmailInfo struct {
	EmailAddress string
	Mail         string
}

// Campaign is the info about the email campaign.
type Campaign struct {
	Name             string
	WelcomeEmail     string
	UnsubscribeEmail string
	Mails            []string
}

// Subscription is the user email and the campaign they'll receive.
type Subscription struct {
	EmailAddress string
	Campaign     Campaign
}

// UserSubscriptionWorkflow handles subscribing users. Accepts a Subsription.
func UserSubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + subscription.EmailAddress)

	// How frequently to send the messages
	duration := time.Minute
	// duration := (24 * 7) * time.Hour

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

		// Cancellation received, which will trigger an unsubscribe email.

		newCtx, _ := workflow.NewDisconnectedContext(ctx)

		data := EmailInfo{
			EmailAddress: subscription.EmailAddress,
			Mail:         subscription.Campaign.UnsubscribeEmail,
		}

		logger.Info("Sending unsubscribe email to %s", subscription.EmailAddress)
		err := workflow.ExecuteActivity(newCtx, SendContentEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()

	logger.Info("Sending welcome email to %s", subscription.EmailAddress)

	data := EmailInfo{
		EmailAddress: subscription.EmailAddress,
		Mail:         subscription.Campaign.WelcomeEmail,
	}

	err := workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	}

	for _, mail := range subscription.Campaign.Mails {

		data := EmailInfo{
			EmailAddress: subscription.EmailAddress,
			Mail:         mail,
		}

		err = workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email %s", "Error", mail, err)
		}

		logger.Info("sent content email %s to %s", mail, subscription.EmailAddress)

		workflow.Sleep(ctx, duration)
	}

	return nil
}

// SendContentEmail is the activity that sends the email to the customer.
func SendContentEmail(ctx context.Context, emailInfo EmailInfo) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending email %s to %s", emailInfo.Mail, emailInfo.EmailAddress)

	// call mailer api here.
	message, err := getEmailFromFile(emailInfo.Mail)

	if err != nil {
		return sendMail(message, emailInfo.EmailAddress)
	}

	logger.Error("Failed getting email", err)
	return errors.New("unable to locate message to send")

}

// getEmailFromFile gets the email from the specified text file.
func getEmailFromFile(filename string) (string, error) {

	file, err := os.Open(filename)

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	return scanner.Text(), scanner.Err()
}

// mocked mail.
func sendMail(message string, email string) error {
	return nil
}
