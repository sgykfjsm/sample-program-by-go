package main

import (
	"fmt"
	"log"

	"github.com/go-gomail/gomail"
)

type Email struct {
	To      []string
	Cc      []string
	Subject string
	Body    string
}

type EmailConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Sender   string
}

type EmailSender struct {
	conf EmailConfig
}

func (e *EmailSender) Send(email Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.conf.Sender)
	m.SetHeader("To", email.To...)
	if len(email.Cc) == 1 {
		for _, c := range email.Cc {
			m.SetAddressHeader("Cc", c, c)
		}
	} else if len(email.Cc) > 1 {
		m.SetHeader("Cc", email.Cc...)
	}
	m.SetHeader("Subject", email.Subject)
	m.SetBody("plain/text", email.Body)
	// m.Attach("/tmp/attach.txt")

	if e.conf.Username != "" && e.conf.Password != "" {
		d := gomail.NewDialer(e.conf.Hostname, e.conf.Port, e.conf.Username, e.conf.Password)
		return d.DialAndSend(m)
	}

	d := gomail.Dialer{
		Host: e.conf.Hostname,
		Port: e.conf.Port,
	}
	return d.DialAndSend(m)
}

func NewEmailSender(conf EmailConfig) EmailSender {
	return EmailSender{conf}
}

func FatalIfNotNil(message string, e error) {
	if e != nil {
		log.Fatalf(message, e.Error())
	}
}

func main() {
	conf := EmailConfig{
		Hostname: "127.0.0.1",
		Port:     1025, // See https://github.com/mailhog/MailHog#getting-started
		Username: "test-user",
		Password: "password",
		Sender:   "test.sender@example.com",
	}
	e := NewEmailSender(conf)

	email := Email{
		To:      []string{"me@example.com", "you@example.com"},
		Cc:      []string{"he@example.com", "she@example.com"},
		Subject: "this is test",
		Body:    "Hello World!",
	}
	err := e.Send(email)
	FatalIfNotNil("Failed to send the mail: %s", err)
	fmt.Printf("Sent mail to %q succesfully\n", email.To)
	fmt.Println("Go http://127.0.0.1:8025 to confirm")
}
