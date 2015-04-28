// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package dispatcher_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/johnweldon/mqd/dispatcher"
)

func TestProcess(t *testing.T) {
	tc := testContext{}
	defer initializeAndTearDown(t, &tc)()

	tc.addFile(t, "From: foo@bar.com\r\nTo: baz@bar.com\r\nSubject: Hello\r\n\r\nMessage body here\r\n")

	if tc.mailqueue == "" {
		t.Fatal("temp mailqueue folder not created")
	}

	q := dispatcher.NewPickupFolderQueue(tc.mailqueue, tc.badmail)
	err := q.Process(testCallback(t))
	if err != nil {
		t.Fatalf("Process call failed: %q", err)
	}
}

func testCallback(t *testing.T) dispatcher.MailQueueCallbackFn {
	return func(data []byte) bool {
		t.Logf("got data: %q", string(data))
		return true
	}
}

func initializeAndTearDown(t *testing.T, tc *testContext) func() {
	q, err := ioutil.TempDir("", "dispatcher_test_mailqueue")
	if err != nil {
		t.Fatalf("error creating temp folder: %q", err)
		return func() {}
	}
	tc.mailqueue = q

	b, err := ioutil.TempDir("", "dispatcher_test_badmail")
	if err != nil {
		t.Fatalf("error creating temp folder: %q", err)
		return func() {}
	}
	tc.badmail = b

	return func() {
		if err := os.RemoveAll(q); err != nil {
			t.Errorf("problem removing %q: %q", q, err)
		}
		if err := os.RemoveAll(b); err != nil {
			t.Errorf("problem removing %q: %q", b, err)
		}
	}
}

type testContext struct {
	mailqueue string
	badmail   string
}

func (tc *testContext) addFile(t *testing.T, contents string) {
	f, err := ioutil.TempFile(tc.mailqueue, "test_file")
	if err != nil {
		t.Fatalf("problem creating temp file: %q", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(contents))
	if err != nil {
		t.Fatalf("problem writing temp file: %q", err)
	}
}
