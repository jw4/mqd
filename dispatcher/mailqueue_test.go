package dispatcher_test

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
}

func initializeAndTearDown() func() {
	folder, err := ioutil.TempDir("", "dispatcher_test_queue")
	if err != nil {
		return func() {}
	}
	return func() { os.RemoveAll(folder) }
}
