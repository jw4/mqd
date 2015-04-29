// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mqd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	mqd_smtp "gopkg.in/mail-queue-dispatcher/dispatcher.v0/smtp"
)

type SmtpAuthType string

const (
	LoginAuth SmtpAuthType = "LOGIN"
	PlainAuth SmtpAuthType = "PLAIN"
)

type Settings struct {
	C         map[string]ConnectionDetails `json:"connections"`
	MailQueue string                       `json:"mailqueue"`
	BadMail   string                       `json:"badmail"`
	Interval  int                          `json:"interval"`
}

func NewSettings(mailqueue, badmail string) *Settings {
	return &Settings{C: map[string]ConnectionDetails{}, MailQueue: mailqueue, BadMail: badmail, Interval: 30}
}

func (s *Settings) String() string {
	return fmt.Sprintf(
		"mailqueue: %s, badmail: %s, interval: %d seconds, details: %s",
		s.MailQueue, s.BadMail, s.Interval, s.C)
}

func (s *Settings) ConnectionForSender(sender string) (ConnectionDetails, error) {
	if details, ok := s.C[sender]; ok {
		return details, nil
	}
	addr, err := mail.ParseAddress(sender)
	if err != nil {
		return ConnectionDetails{}, err
	}
	if details, ok := s.C[addr.Address]; ok {
		return details, nil
	}
	if details, ok := s.C[strings.ToLower(addr.Address)]; ok {
		return details, nil
	}
	return ConnectionDetails{}, fmt.Errorf("connection details not found for %q", sender)
}

func ReadSettingsFrom(r io.Reader) (*Settings, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return UnmarshalSettings(raw)
}

func UnmarshalSettings(raw []byte) (*Settings, error) {
	s := &Settings{}
	err := unmarshalSettings(raw, s)
	if err != nil {
		return s, err
	}

	if s.Interval < 5 || s.Interval > 3600 {
		s.Interval = 30
	}
	return s, nil
}

func WriteSettings(path string, s *Settings) error {
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

func (d *ConnectionDetails) String() string {
	return fmt.Sprintf(
		"sender: %s, authtype: %s, server: %s, host: %s, username: %s, password: ******",
		d.Sender, d.AuthType, d.Server, d.Host, d.Username)
}

func unmarshalSettings(data []byte, s *Settings) error { return json.Unmarshal(data, s) }
