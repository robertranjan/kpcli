package cmd

import (
	"fmt"
	// "log"
	"os"

	gomail "gopkg.in/gomail.v2"
)

// TODO: below configs should moved out to a config file
var (
	// Configuration
	from         = "yourEmail@gmail.com"
	password     = os.Getenv("keepass_gmail_app_password")
	to           = []string{"yourEmail@gmail.com"}
	smtpHost     = "smtp.gmail.com"
	smtpPort     = 587
	subject      = "here are the KDBX changes since last backup!"
	emailContent = "will be generated during execution"
)

func (d *Diff) Notify(contentFile string) {

	emailContentByte, err := os.ReadFile(contentFile)
	if err != nil {
		log.Fatalf("failed to open file: %v, err: %v", contentFile, err)
	}
	emailContent = "<pre>" + string(emailContentByte) + "</pre>"

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to[0])
	msg.SetHeader("Subject", subject)
	// msg.SetBody("text/html", "<b>This is the body of the mail</b>")
	msg.SetBody("text/html", emailContent)
	msg.Attach(contentFile)

	n := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		fmt.Printf("failed to notify: %v", err)
	}
}
