package main

import gomail "gopkg.in/gomail.v2"

type EmailRequest struct {
	ReplyTo  string
	To       string
	Subject  string
	HTMLBody string
	Images   []string
}

func (config EmailerConfig) Email(request EmailRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.SendingAddress)
	m.SetHeader("Reply-To", request.ReplyTo)
	m.SetHeader("To", request.To)
	m.SetHeader("Subject", request.Subject)
	m.SetBody("text/html", request.HTMLBody)

	for _, Image := range request.Images {
		m.Attach(Image)
	}

	d := gomail.NewDialer(config.URL, 587, config.SendingAddress, config.SendingPassword)

	return d.DialAndSend(m)
}
