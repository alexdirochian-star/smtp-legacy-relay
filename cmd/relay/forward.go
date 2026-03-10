package main

import (
	"net/smtp"
	"os"
	"strings"
)

const smtpServer = "smtp.gmail.com:587"

var (
	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASS")
)

func forwardMessage(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	msg := string(data)

	rcpt := extractRecipient(msg)

	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPass,
		"smtp.gmail.com",
	)

	err = smtp.SendMail(
		smtpServer,
		auth,
		smtpUser,
		[]string{rcpt},
		data,
	)

	return err
}

func extractRecipient(msg string) string {

	for _, line := range strings.Split(msg, "\n") {

		if strings.HasPrefix(line, "X-Relay-RcptTo:") {

			addr := strings.TrimSpace(
				strings.TrimPrefix(line, "X-Relay-RcptTo:"),
			)

			addr = strings.Trim(addr, "<>")

			return addr
		}
	}

	return smtpUser
}