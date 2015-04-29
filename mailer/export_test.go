// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

var (
	ParseEmail     = parseEmail
	FindSender     = findSender
	FindRecipients = findRecipients
)

type SenderFunc senderFunc

func DummySender(m Mailer, sf SenderFunc) func() {
	sm := m.(*smtpMailer)
	orig := sm.sendFn
	sm.sendFn = senderFunc(sf)
	return func() { sm.sendFn = orig }
}
