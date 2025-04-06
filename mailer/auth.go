package mailer

const (
	SEND_MAGIC_LINK_TEMPLATE = "send_magic_link"
)

func (m *GoMailer) SendMagicLinkMail(to, magicLink string) error {
	return m.sendEmail(
		to,
		"Your magic link 🚀",
		SEND_MAGIC_LINK_TEMPLATE,
		map[string]any{"Link": magicLink},
	)
}
