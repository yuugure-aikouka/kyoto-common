package mailer

import (
	"bytes"
	"log"
	"text/template"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	fromEmail string
	client    *mail.Client
	template  *template.Template
}

type SMTPConfig struct {
	Host     string
	Username string
	Password string
}

func NewMailer(fromEmail string, smtpConfig SMTPConfig) *Mailer {
	client, err := mail.NewClient(
		smtpConfig.Host,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpConfig.Username),
		mail.WithPassword(smtpConfig.Password),
	)
	if err != nil {
		log.Fatalf("Failed to create mail client: %s\n", err)
	}

	tmpl := template.Must(template.ParseGlob("mailer/templates/*.html"))

	return &Mailer{
		fromEmail: fromEmail,
		client:    client,
		template:  tmpl,
	}
}

func (m *Mailer) sendEmail(to, subject, templateName string, data any) error {
	message := mail.NewMsg()
	if err := message.From(m.fromEmail); err != nil {
		return err
	}

	if err := message.To(to); err != nil {
		return err
	}

	content, err := m.createMailContent(templateName, data)
	if err != nil {
		return err
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, content)

	if err := m.client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) createMailContent(templateName string, data any) (string, error) {
	var body bytes.Buffer
	if err := m.template.ExecuteTemplate(&body, templateName+".html", data); err != nil {
		return "", err
	}

	return body.String(), nil
}
