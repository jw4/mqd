package mqd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"

	mqd_smtp "github.com/johnweldon/mqd/smtp"
)

type SmtpAuthType string

const (
	LoginAuth SmtpAuthType = "LOGIN"
	PlainAuth SmtpAuthType = "PLAIN"
)

type Settings struct {
	Connections map[string]ConnectionDetails `json:"connections"`
	MailQueue   string                       `json:"mailqueue"`
	BadMail     string                       `json:"badmail"`
}

func NewSettings(mailqueue, badmail string) Settings {
	return Settings{Connections: map[string]ConnectionDetails{}, MailQueue: mailqueue, BadMail: badmail}
}

func ReadSettings(path string) (Settings, error) {
	s := Settings{}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return s, err
	}

	err = unmarshalSettings(raw, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func WriteSettings(path string, s Settings) error {
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	_, err = fi.Write(bytes)
	if err != nil {
		return err
	}

	return nil
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

func unmarshalSettings(data []byte, s *Settings) error { return json.Unmarshal(data, s) }
