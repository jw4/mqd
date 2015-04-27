package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"

	config "github.com/johnweldon/mqd/config"
)

var logger = log.New(os.Stderr, "mqd.mailer: ", log.Lshortfile)

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
		logger.Printf("ERROR: parsing email: %q", err)
		return false
	}
	sender := findSender(eml)
	recipients := findRecipients(eml)
	if err := m.Send(sender, recipients, message); err != nil {
		logger.Printf("ERROR: sending from %q to %v: %q", sender, recipients, err)
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
	senders := getMatchingSlices(eml.Header, "X-Sender", "From")
	if len(senders) < 1 {
		return ""
	}
	return senders[0]
}

func findRecipients(eml *mail.Message) []string {
	recipients := []string{}
	if eml == nil {
		return recipients
	}
	return getMatchingSlices(eml.Header, "To", "Cc", "Bcc")
}

func getMatchingSlices(in map[string][]string, keys ...string) []string {
	slice := []string{}
	for _, key := range keys {
		if items, ok := in[key]; ok {
			slice = append(slice, items...)
		}
	}
	return slice
}
