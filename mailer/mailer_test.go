// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer_test

import (
	"testing"

	"github.com/johnweldon/mqd/mailer"
)

func TestFindSender(t *testing.T) {
	msg, err := mailer.ParseEmail([]byte("X-Sender: \"asdf\" <asdf@qwer.ty>\r\nFrom: \"qwer\" <asdf@qwer.ty>\r\n\r\nasdf\r\n"))
	if err != nil {
		t.Error(err)
	}

	sender := mailer.FindSender(msg)
	if sender != "asdf@qwer.ty" {
		t.Errorf("sender incorrect, got %q", sender)
	}
}

func TestFindRecipients(t *testing.T) {
	tests := []struct {
		message    []byte
		recipients []string
	}{{
		message:    []byte("To: asdf@wert.yo, qwer@asdf.gh\r\nCc: qwer@qwer.yt\r\nBcc: zxcv@zxcv.as\r\n\r\nasdf\r\n"),
		recipients: []string{"asdf@wert.yo", "qwer@asdf.gh", "qwer@qwer.yt", "zxcv@zxcv.as"},
	}}

	for ix, test := range tests {
		t.Logf("Test %d", ix)
		msg, err := mailer.ParseEmail(test.message)
		if err != nil {
			t.Error(err)
		}
		recipients := mailer.FindRecipients(msg)
		if len(recipients) != len(test.recipients) {
			t.Errorf("mismatch, expected %d recipients, got %d", len(test.recipients), len(recipients))
		}
	}
}
