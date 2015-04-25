package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"

	js "github.com/johnweldon/mqd/smtp"
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
		return js.LoginAuth(d.Username, d.Password), nil
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

func (m *smtpMailer) LoadSettings(data []byte) error {
	s := NewSettings()
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	for key, val := range s.Connections {
		m.settings.Connections[key] = val
	}
	return nil
}

func (m *smtpMailer) Send(sender string, message []byte) error {
	if connection, ok := m.settings.Connections[sender]; ok {
		auth, err := connection.Auth()
		if err != nil {
			return err
		}
		return smtp.SendMail(connection.Server, auth, connection.Sender, findRecipients(&message), message)
	}
	return fmt.Errorf("no connection settings found for %q", sender)
}

func findRecipients(msg *[]byte) []string {
	//logger.Printf("findRecipients(msg: %q)", string(*msg))
	recipients := []string{}
	if msg == nil {
		return recipients
	}

	buf := bytes.NewBuffer(*msg)
	eml, err := mail.ReadMessage(buf)
	if err != nil {
		return recipients
	}

	return getPossibleSlices(eml.Header, "To", "Cc", "Bcc")
}

func getPossibleSlices(in map[string][]string, keys ...string) []string {
	slice := []string{}
	for _, key := range keys {
		if items, ok := in[key]; ok {
			slice = append(slice, items...)
		}
	}
	return slice
}
