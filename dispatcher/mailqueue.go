package dispatcher

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var logger = log.New(os.Stderr, "mqd.dispatcher: ", log.Lshortfile)

type folderQueue struct {
	mailqueue string
	badmail   string
}

func NewPickupFolderQueue(mailqueue, badmail string) MailQueueDispatcher {
	return &folderQueue{mailqueue: mailqueue, badmail: badmail}
}

func (q *folderQueue) Process(callbackFn MailQueueCallbackFn) error {
	return filepath.Walk(q.mailqueue, q.processItem(callbackFn))
}

func (q *folderQueue) processItem(fn MailQueueCallbackFn) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err2 error) error {
		logger.Printf("INFO: processing %q", path)
		if info.IsDir() {
			if path == q.mailqueue {
				return nil
			}
			return filepath.SkipDir
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			q.markBad(path, info)
			return nil
		}

		if fn(raw) {
			q.markComplete(path)
		} else {
			q.markBad(path, info)
		}
		return nil
	}
}

func (q *folderQueue) markBad(path string, info os.FileInfo) {
	logger.Printf("ERROR: processing %q was unsuccessful", path)
	target := filepath.Join(q.badmail, info.Name())
	err := os.Rename(path, target)
	if err != nil {
		logger.Printf("ERROR: moving %q to %q: %q", path, target, err)
	}
}

func (q *folderQueue) markComplete(path string) {
	if err := os.Remove(path); err != nil {
		logger.Printf("ERROR: removing file %q: %q", path, err)
	}
}
