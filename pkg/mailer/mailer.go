package mailer

import (
	gomail "gopkg.in/mail.v2"
)

type Mailer struct {
	email    string
	password string
	host     string
	port     int
}

func NewMailer(email, password, host string, port int) *Mailer {
	return &Mailer{
		email:    email,
		password: password,
		host:     host,
		port:     port,
	}
}

func (m *Mailer) Send(recipient, subject, body string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.email)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	d := gomail.NewDialer(m.host, m.port, m.email, m.password)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
