package sesops

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggarcia209/go-aws/go-ses/goses"
	"github.com/ggarcia209/portfolio/service/util/awsops"
)

const (
	// EnvarServiceEmailAddress contains the env variable name for the service email.
	// This automated no-reply email programitcally sends customer-facing and internal email notifications.
	EnvarServiceEmailAddress = "SES_SERVICE_EMAIL_ADDRESS"
	// EnvarAdminEmailAddress contains the env variable name for the admin email.
	EnvarAdminEmailAddress = "SES_ADMIN_EMAIL_ADDRESS"
)

// ServiceEmailAddress returns the env variable value for the service email.
func ServiceEmailAddress() string { return os.Getenv(EnvarServiceEmailAddress) }

// AdminEmailAddress returns the env variable value for the admin email.
func AdminEmailAddress() string { return os.Getenv(EnvarAdminEmailAddress) }

// InitSesh encapsulates the gosns.InitSesh() method and returns the SNS service
// as an interface{} type.
func InitSesh() interface{} {
	svc := goses.InitSesh()
	return svc
}

// InitSesh encapsulates the goses.InitSesh() method and returns the SES service
// as an interface{} type.
func InitClient(session awsops.AWSSession) interface{} {
	svc := goses.NewSESClient(session.Session)
	return svc
}

// SendShippingNotification sends an order notification email intended for the business admin and/or fulfillment team.
// 'from' specifies the 'from' address (ex: orders@store.com), 'notifyEmail' specifies the 'to' address (ex: fulfillment@store.com).
func SendContactRequest(svc interface{}, src, target, from, subject, msg string) error {
	subjectLine := "gcodes.io | New Contact Request"
	body := fmt.Sprintf("From: %s\nSubject: %s\nMessage: %s\n", from, subject, msg)
	err := sendEmailWithRetry(svc, src, target, subjectLine, body)
	if err != nil {
		log.Printf("SendContactRequest failed: %v", err)
		return err
	}

	return nil
}

// TO DO: add error handling for bounced messages
func sendEmailWithRetry(svc interface{}, from, to, subject, text string) error {
	// poll for messages with exponential backoff for errors & empty responses
	retries := 0
	maxRetries := 4
	backoff := 1000.0
	for {
		// receive messages from queue
		err := goses.SendPlainTextEmail(svc, []string{to}, []string{}, []string{}, from, subject, text)
		if err != nil {
			// retry with backoff if error
			if retries > maxRetries {
				log.Printf("sendEmailWithRetry failed: %v -- max retries exceeded", err)
				return err
			}
			log.Printf("sendEmailWithRetry failed: %v -- retrying...", err)
			time.Sleep(time.Duration(backoff) * time.Millisecond)
			backoff = backoff * 2
			retries++
			continue
		}

		return nil
	}
}
