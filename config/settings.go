// Copyright 2015-2017 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	mqd_smtp "github.com/jw4/mqd/smtp"
)

// SMTPAuthType names smtp authentication methods
type SMTPAuthType string

const (
	LoginAuth SMTPAuthType = "LOGIN"
	PlainAuth SMTPAuthType = "PLAIN"
)

// Settings holds the configuration parameters used by the mail queue
// dispatcher.
type Settings struct {
	C         map[string]ConnectionDetails `json:"connections"`
	MailQueue string                       `json:"mailqueue"`
	BadMail   string                       `json:"badmail"`
	SentMail  string                       `json:"sentmail"`
	Interval  int                          `json:"interval"`
}

// NewSettings generates a new Settings configuration, initializing it
// with the supplied mailqueue and badmail folders, and an empty map
// of connection details.
func NewSettings(mailqueue string, badmail string) *Settings {
	return &Settings{C: map[string]ConnectionDetails{}, MailQueue: mailqueue, BadMail: badmail, Interval: 30}
}

// String fulfills the fmt.Stringer interface
func (s *Settings) String() string {
	return fmt.Sprintf(
		"mailqueue: %s, badmail: %s, interval: %d seconds, details: %s",
		s.MailQueue, s.BadMail, s.Interval, s.C)
}

// ConnectionForSender uses the supplied sender and tries to find
// ConnectionDetails that match the email address.
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

// ReadSettingsFrom uses an io.Reader to read in JSON that is parsed
// into a Settings configuration object.
func ReadSettingsFrom(r io.Reader) (*Settings, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return UnmarshalSettings(raw)
}

// UnmarshalSettings converts raw bytes that are expected to be in JSON
// format into the Settings configuration object.
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

// WriteSettingsTo marshals the Settings config struct and sends it to
// the io.Writer, returning any errors encountered along the way.
func WriteSettingsTo(w io.Writer, s *Settings) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

// WriteSettings is a helper function to open a file and send it to the
// WriteSettingsTo function.
func WriteSettings(path string, s *Settings) error {
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	return WriteSettingsTo(fi, s)
}

// ConnectionDetails describe the metadata needed to complete an SMTP
// connection for a given sender.
type ConnectionDetails struct {
	// Sender should be the plain, lowercase email address of the
	// account that the email will originate from.
	Sender string `json:"sender,omitempty"`
	// Server is the string representing the smtp host and port joined
	// by a colon.  e.g: smtp.example.com:25
	Server string `json:"server,omitempty"`
	// Host is just the host portion.  e.g: smtp.example.com
	Host string `json:"host,omitempty"`
	// AuthType shows which authentication mechanism should be used
	// when connecting to this Server.
	AuthType SMTPAuthType `json:"authtype"`
	// Username of the Sender account.
	Username string `json:"username,omitempty"`
	// Password of the Sender account.
	Password string `json:"password,omitempty"`
}

// Auth returns an implementation of the smtp.Auth interface that can
// be used to perform the smtp authentication for this ConnectionDetails
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

// String implements the fmt.Stringer interface
func (d *ConnectionDetails) String() string {
	return fmt.Sprintf(
		"sender: %s, authtype: %s, server: %s, host: %s, username: %s, password: ******",
		d.Sender, d.AuthType, d.Server, d.Host, d.Username)
}

func unmarshalSettings(data []byte, s *Settings) error {
	return json.Unmarshal(data, s)
}
