package email

type Mailer interface {
	SendSimpleMessage(string)
}
