// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/golang/glog"

	config "github.com/johnweldon/mqd/config"
)

type smtpMailer struct {
	settings config.Settings
}

func NewMailer(s config.Settings) Mailer {
	return &smtpMailer{settings: s}
}

func (m *smtpMailer) LoadSettings(s config.Settings) error {
	for key, val := range s.Connections {
		m.settings.Connections[key] = val
	}
	return nil
}

func (m *smtpMailer) Send(sender string, recipients []string, message []byte) error {
	if connection, ok := m.settings.Connections[sender]; ok {
		auth, err := connection.Auth()
		if err != nil {
			return err
		}
		return smtp.SendMail(connection.Server, auth, connection.Sender, recipients, message)
	}
	return fmt.Errorf("no connection settings found for %q", sender)
}

func (m *smtpMailer) ConvertAndSend(message []byte) bool {
	eml, err := parseEmail(message)
	if err != nil {
		glog.Errorf("parsing email: %q", err)
		return false
	}
	sender := findSender(eml)
	recipients := findRecipients(eml)
	if err := m.Send(sender, recipients, message); err != nil {
		glog.Errorf("sending from %q to %v: %q", sender, recipients, err)
		return false
	}
	return true
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
	senders := getEmails(eml.Header, "X-Sender", "From")
	if len(senders) < 1 {
		return ""
	}
	for _, sender := range senders {
		addr, err := mail.ParseAddress(sender)
		if err == nil {
			return addr.Address
		}
	}
	return ""
}

func findRecipients(eml *mail.Message) []string {
	recipients := []string{}
	if eml == nil {
		return recipients
	}
	return getEmails(eml.Header, "To", "Cc", "Bcc")
}

func getEmails(in map[string][]string, keys ...string) []string {
	slice := []string{}
	for _, key := range keys {
		if items, ok := in[key]; ok {
			for _, item := range items {
				list, err := mail.ParseAddressList(item)
				if err != nil {
					glog.Warningf("failed to email list %q: %v", item, err)
				} else {
					for _, addr := range list {
						slice = append(slice, addr.Address)
					}
				}
			}
		}
	}
	return slice
}
