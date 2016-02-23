// Copyright 2015-2016 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	"net/smtp"

	config "github.com/jw4/mqd/config"
)

// EmailSender is an ad-hoc interface to describe the SendMail function
// that is in the net/smtp package.
type EmailSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// Mailer describes an object that is able to send emails via the
// EmailSender interface, and that can load settings and convert
// raw email messages into sent mail.
type Mailer interface {
	EmailSender
	LoadSettings(*config.Settings) error
	ConvertAndSend(email []byte) bool
}
