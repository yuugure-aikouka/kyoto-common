package helper

import (
	"github.com/yuugure-aikouka/kyoto-common/mailer"
)

type SentMail struct {
	To   string
	Link string
}

type MockMailer struct {
	SentEmails []SentMail
}

func NewMockMailer() mailer.Mailer {
	return &MockMailer{}
}

func (m *MockMailer) SendMagicLinkMail(to, magicLink string) error {
	m.SentEmails = append(m.SentEmails, SentMail{To: to, Link: magicLink})
	return nil
}

func (m *MockMailer) Clear() {
	m.SentEmails = m.SentEmails[:0]
}
