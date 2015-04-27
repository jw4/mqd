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
	return filepath.Walk(q.mailqueue, processItem(q.mailqueue, callbackFn))
}

func processItem(root string, fn MailQueueCallbackFn) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err2 error) error {
		logger.Printf("INFO: processing %q", path)
		if info.IsDir() {
			if path == root {
				return nil
			}
			return filepath.SkipDir
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			// TODO(jw4) move to badmail queue, maybe after a few retries
			logger.Printf("ERROR: reading file %q: %q", path, err)
			return nil
		}

		if fn(raw) {
			defer func() {
				if err = os.Remove(path); err != nil {
					// TODO(jw4)  possibly mark for deletion somehow
					logger.Printf("ERROR: removing file %q: %q", path, err)
				}
			}()
		} else {
			// TODO(jw4) move to badmail queue, or possibly retry a few times
			logger.Printf("ERROR: processing %q was unsuccessful", path)
		}
		return nil
	}
}
