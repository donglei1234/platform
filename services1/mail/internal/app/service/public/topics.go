package public

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	MailTopic        Topic = "mail"
	MailPrivateTopic Topic = "mail.private"
	MailPublicTopic  Topic = "mail.public"
)
