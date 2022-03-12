package email

import "time"

type EmailLog struct {
	To         string
	Content    string
	Provider   string
	Error      string
	StatusCode int
	ReceivedAt time.Time
	CreatedAt  time.Time
}
