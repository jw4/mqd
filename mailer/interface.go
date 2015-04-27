package mailer

import (
	config "github.com/johnweldon/mqd/config"
)

type Mailer interface {
	LoadSettings(config.Settings) error
	Send(sender string, recipients []string, message []byte) error
	ConvertAndSend(email []byte) bool
}
