package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"

	mqd_smtp "github.com/johnweldon/mqd/smtp"
)

var logger = log.New(os.Stderr, "mqd.mailer: ", log.Lshortfile)

type SmtpAuthType string

const (
	LoginAuth SmtpAuthType = "LOGIN"
	PlainAuth SmtpAuthType = "PLAIN"
)

type Settings struct {
	Connections map[string]ConnectionDetails `json:"connections"`
}

func NewSettings() Settings {
	return Settings{Connections: map[string]ConnectionDetails{}}
}

type ConnectionDetails struct {
	Sender   string       `json:"sender,omitempty"`
	Server   string       `json:"server,omitempty"`
	Host     string       `json:"host,omitempty"`
	AuthType SmtpAuthType `json:"authtype"`
	Username string       `json:"username,omitempty"`
	Password string       `json:"password,omitempty"`
}

func (d *ConnectionDetails) Auth() (smtp.Auth, error) {
	if d == nil {
		return nil, fmt.Errorf("invalid connection info; nil")
	}

	switch d.AuthType {
	case LoginAuth:
		return mqd_smtp.LoginAuth(d.Username, d.Password), nil
	case PlainAuth:
		return smtp.PlainAuth("", d.Username, d.Password, d.Host), nil
	}

	return nil, fmt.Errorf("unknown auth type %q", string(d.AuthType))
}

type smtpMailer struct {
	settings Settings
}

func NewMailer() Mailer {
	return &smtpMailer{settings: NewSettings()}
}

func (m *smtpMailer) LoadSettings(s *Settings) error {
	if s == nil {
		return fmt.Errorf("passed in settings pointer is nil")
	}
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

func ReadSettingsFile(path string) *Settings {
	s := NewSettings()
	if raw, err := ioutil.ReadFile(path); err == nil {
		if err := unmarshalSettings(raw, &s); err != nil {
			logger.Printf("ERROR: unmarshalSettings(raw: %q) : %q", string(raw), err)
		}
	}
	return &s
}

func WriteSettingsFile(path string, s *Settings) error {
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	bytes, err := json.MarshalIndent(*s, "", "  ")
	if err != nil {
		return err
	}

	_, err = fi.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalSettings(data []byte, s *Settings) error { return json.Unmarshal(data, s) }

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
