// Copyright 2015-2016 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/golang/glog"

	config "github.com/jw4/mqd/config"
)

type senderFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

type smtpMailer struct {
	sendFn   senderFunc
	settings *config.Settings
}

// NewMailer returns a Mailer implementation using config.Settings
// to transmit emails.
func NewMailer(s *config.Settings) Mailer {
	return &smtpMailer{settings: s, sendFn: smtp.SendMail}
}

// LoadSettings updates the Mailer configuration given the supplied
// config.Settings.
func (m *smtpMailer) LoadSettings(s *config.Settings) error {
	m.settings = s
	return nil
}

// ConvertAndSend takes in a raw []byte message, parses it to discover
// the sender (preferring From: then X-Sender:) and the recipients
// and then transmits the email through the configured smtp settings
// for the sender.
func (m *smtpMailer) ConvertAndSend(message []byte) bool {
	eml, err := parseEmail(message)
	if err != nil {
		glog.Errorf("parsing email: %q", err)
		return false
	}
	sender := findSender(eml)
	recipients := findRecipients(eml)
	if err := m.send(sender, recipients, message); err != nil {
		glog.Errorf("sending from %q to %v: %q", sender, recipients, err)
		return false
	}
	return true
}

// SendMail fulfills the EmailSender interface.  It wraps an internal
// function that can be swapped out to use different sending techniques.
func (m *smtpMailer) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return m.sendFn(addr, a, from, to, msg)
}

func (m *smtpMailer) send(sender string, recipients []string, message []byte) error {
	connection, err := m.settings.ConnectionForSender(sender)
	if err != nil {
		return err
	}

	auth, err := connection.Auth()
	if err != nil {
		return err
	}

	return m.SendMail(connection.Server, auth, connection.Sender, recipients, message)
}

func parseEmail(msg []byte) (*mail.Message, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	return mail.ReadMessage(bytes.NewBuffer(msg))
}

func findSender(eml *mail.Message) string {
	if eml == nil {
		return ""
	}
	// prefer X-Sender
	senders := getEmails(eml.Header, "X-Sender")
	if len(senders) > 0 {
		for _, sender := range senders {
			addr, err := mail.ParseAddress(sender)
			if err == nil {
				return addr.Address
			}
		}
	}
	// fallback to From
	senders = getEmails(eml.Header, "From")
	if len(senders) > 0 {
		for _, sender := range senders {
			addr, err := mail.ParseAddress(sender)
			if err == nil {
				return addr.Address
			}
		}
	}
	return ""
}

func findRecipients(eml *mail.Message) []string {
	recipients := []string{}
	if eml == nil {
		return recipients
	}
	return getEmails(eml.Header, "To", "Cc", "Bcc", "X-Receiver")
}

func getEmails(in map[string][]string, keys ...string) []string {
	itemMap := map[string]interface{}{}
	for _, key := range keys {
		if items, ok := in[key]; ok {
			for _, item := range items {
				list, err := mail.ParseAddressList(item)
				if err != nil {
					glog.Warningf("failed to email list %q: %v", item, err)
				} else {
					for _, addr := range list {
						itemMap[addr.Address] = nil
					}
				}
			}
		}
	}
	slice := []string{}
	for k := range itemMap {
		slice = append(slice, k)
	}
	return slice
}
