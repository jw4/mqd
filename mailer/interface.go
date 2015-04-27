package mailer

type Mailer interface {
	LoadSettings(*Settings) error
	Send(sender string, recipients []string, message []byte) error
	ConvertAndSend(email []byte) bool
}
