// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	"net/smtp"

	config "github.com/johnweldon/mqd/config"
)

type SenderFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

type EmailSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type Mailer interface {
	EmailSender
	LoadSettings(config.Settings) error
	Send(sender string, recipients []string, message []byte) error
	ConvertAndSend(email []byte) bool
}
