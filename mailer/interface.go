package mailer

type Mailer interface {
	LoadSettings(json []byte) error
	Send(sender string, message []byte) error
}
