// Copyright 2015-2016 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	"net/smtp"
	"testing"

	config "github.com/jw4/mqd/config"
)

func TestFindSender(t *testing.T) {
	tests := []struct {
		message []byte
		sender  string
	}{{
		message: []byte("X-Sender: \"asdf\" <asdf@qwer.ty>\r\nFrom: \"qwer\" <asdf@qwer.ty>\r\n\r\nasdf\r\n"),
		sender:  "asdf@qwer.ty",
	}, {
		message: []byte("X-Sender: asdf@qwer.ty\r\nFrom: qwer@qwer.ty\r\n\r\nasdf\r\n"),
		sender:  "asdf@qwer.ty",
	}, {
		message: []byte("X-Sender: qwer@qwer.ty\r\nFrom: asdf@qwer.ty\r\n\r\nasdf\r\n"),
		sender:  "qwer@qwer.ty",
	}}

	for ix, test := range tests {
		t.Logf("Test %d", ix)
		msg, err := parseEmail(test.message)
		if err != nil {
			t.Error(err)
		}

		sender := findSender(msg)
		if sender != test.sender {
			t.Errorf("sender incorrect, expected %q got %q", test.sender, sender)
		}
	}
}

func TestFindRecipients(t *testing.T) {
	tests := []struct {
		message    []byte
		recipients []string
	}{{
		message:    []byte("To: asdf@wert.yo, qwer@asdf.gh\r\nCc: qwer@qwer.yt\r\nBcc: zxcv@zxcv.as\r\n\r\nasdf\r\n"),
		recipients: []string{"asdf@wert.yo", "qwer@asdf.gh", "qwer@qwer.yt", "zxcv@zxcv.as"},
	}, {
		message:    []byte("To: asdf@wert.yo\r\nCc: yummy@gummy.br\r\nBcc: zxcv@zxcv.as\r\nX-Receiver: \"Yummy Gummy\" <yummy@gummy.br>\r\n\r\nasdf\r\n"),
		recipients: []string{"asdf@wert.yo", "zxcv@zxcv.as", "yummy@gummy.br"},
	}}

	for ix, test := range tests {
		t.Logf("Test %d", ix)
		msg, err := parseEmail(test.message)
		if err != nil {
			t.Error(err)
		}
		recipients := findRecipients(msg)
		if len(recipients) != len(test.recipients) {
			t.Errorf("mismatch, expected %d recipients, got %d", len(test.recipients), len(recipients))
		}
	}
}

func TestConvertAndSend(t *testing.T) {
	tests := []struct {
		message    []byte
		shouldpass bool
	}{{
		message:    []byte("To: asdf@qwer.ty\r\nFrom: qwer@asdf.gh\r\nSubject: hello\r\n\r\nqwer\r\n"),
		shouldpass: true,
	}, {
		message:    []byte("To: asdf@qwer.ty\r\nFrom: \"BAZ\" <baz@foo.com>\r\nSubject: hello\r\n\r\nqwer\r\n"),
		shouldpass: true,
	}}
	m := testMailer(t)

	for ix, test := range tests {
		t.Logf("Test %d", ix)
		pass := m.ConvertAndSend(test.message)
		if pass != test.shouldpass {
			t.Errorf("ConvertAndSend returned %t, expected %t", pass, test.shouldpass)
		}
	}
}

var (
	testConfig = []byte(`{
    "interval": 45,
    "mailqueue": "mailqueue",
    "badmail": "badmail",
    "connections": {
        "qwer@asdf.gh": {
        "sender": "foo@bar.com",
        "server": "localhost:587",
        "host": "localhost",
        "authtype": "PLAIN",
        "username": "asdf",
        "password": "qwer"
        },
        "\"BAZ\" <baz@foo.com>": {
        "sender": "baz@foo.com",
        "server": "smtp.foo.com:587",
        "host": "smtp.foo.com",
        "authtype": "LOGIN",
        "username": "baz@foo.com",
        "password": "pazzwerd"
        },
        "baz@foo.com": {
        "sender": "baz@foo.com",
        "server": "smtp.foo.com:587",
        "host": "smtp.foo.com",
        "authtype": "LOGIN",
        "username": "baz@foo.com",
        "password": "pazzwerd"
        }}}`)
)

func testMailer(t *testing.T) Mailer {
	cfg, err := config.UnmarshalSettings(testConfig)
	if err != nil {
		t.Errorf("unmarshal fail: %v", err)
	}

	m := NewMailer(cfg)

	fn := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		t.Logf("Sending.. addr: %q, auth: %s, from: %q, to: %v, len: %d", addr, a, from, to, len(msg))
		return nil
	}
	dummySender(m, fn)
	return m
}

// dummySender stubs out the actual smtp.Send function to facilitate
// testing.  It returns a func() that can be used in a defer statement
// to restore the original Send function.
func dummySender(m Mailer, sf senderFunc) func() {
	sm := m.(*smtpMailer)
	orig := sm.sendFn
	sm.sendFn = senderFunc(sf)
	return func() { sm.sendFn = orig }
}
