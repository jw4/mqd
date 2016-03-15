// Copyright 2015-2016 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

/*
Package dispatcher provides the mechanism for watching and processing
emails that appear in a designated pickup folder, and removes them
after successfully sending them, or else moves the failed messages to
the configured badmail folder.
*/
package dispatcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/glog"
)

type folderQueue struct {
	mailqueue string
	badmail   string
	sentmail  string
}

// NewPickupFolderQueue returns an implementation of MailQueueDispatcher
// that watches a mailqueue folder and consumes messages that are left
// there and if it fails to consume the message, moves the message to a
// badmail folder. If sentmail is a valid folder, successful emails will
// be moved there after sending.
func NewPickupFolderQueue(mailqueue string, badmail string, sentmail string) MailQueueDispatcher {
	return &folderQueue{mailqueue: mailqueue, badmail: badmail, sentmail: sentmail}
}

// Process implements the MailQueueDispatcher interface, and walks the
// mailqueue folder and sends the found messages to the callbackFn.
func (q *folderQueue) Process(callbackFn MailQueueCallbackFn) error {
	return filepath.Walk(q.mailqueue, q.processItem(callbackFn))
}

func (q *folderQueue) processItem(fn MailQueueCallbackFn) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err2 error) error {
		glog.V(2).Infof("processing %q", path)
		if info == nil {
			glog.Warningf("path not found for %q", path)
			return nil
		}
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
			q.markComplete(path, info)
		} else {
			q.markBad(path, info)
		}
		return nil
	}
}

func (q *folderQueue) markBad(path string, info os.FileInfo) {
	glog.Errorf("processing %q was unsuccessful", path)
	target := filepath.Join(q.badmail, info.Name())
	if err := os.Rename(path, target); err != nil {
		glog.Errorf("moving %q to %q: %q", path, target, err)
	}
}

func (q *folderQueue) markComplete(path string, info os.FileInfo) {
	if sm, err := os.Stat(q.sentmail); err == nil && sm.IsDir() {
		target := filepath.Join(q.sentmail, time.Now().Format("2006-01-02_150405-")+info.Name()+".sent")
		if err = os.Rename(path, target); err != nil {
			glog.Errorf("moving %q to %q: %q", path, target, err)
		}
	} else if err := os.Remove(path); err != nil {
		glog.Errorf("removing file %q: %q", path, err)
	}
}
