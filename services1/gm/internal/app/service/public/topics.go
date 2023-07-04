package public

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	TopicNone  Topic = ""
	GMTopicBan Topic = "ban"
)
