package dispatcher_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/johnweldon/mqd/dispatcher"
)

type testContext struct {
	folder string
}

func (tc *testContext) addFile(t *testing.T, contents string) {
	f, err := ioutil.TempFile(tc.folder, "test_file")
	if err != nil {
		t.Fatalf("problem creating temp file: %q", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(contents))
	if err != nil {
		t.Fatalf("problem writing temp file: %q", err)
	}
}

func TestProcess(t *testing.T) {
	tc := testContext{}
	defer initializeAndTearDown(t, &tc)()

	tc.addFile(t, "From: foo@bar.com\r\nTo: baz@bar.com\r\nSubject: Hello\r\n\r\nMessage body here\r\n")

	if tc.folder == "" {
		t.Fatal("temp mailqueue folder not created")
	}

	q := dispatcher.NewPickupFolderQueue(tc.folder)
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
	folder, err := ioutil.TempDir("", "dispatcher_test_queue")
	if err != nil {
		t.Fatalf("error creating tempdir %q", err)
		return func() {}
	}
	t.Logf("new folder: %q", folder)
	tc.folder = folder
	return func() {
		if err := os.RemoveAll(folder); err != nil {
			t.Errorf("problem removing %q: %q", folder, err)
		}
	}
}
