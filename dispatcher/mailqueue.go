// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package dispatcher

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

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
		glog.Infof("processing %q", path)
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
	glog.Errorf("processing %q was unsuccessful", path)
	target := filepath.Join(q.badmail, info.Name())
	err := os.Rename(path, target)
	if err != nil {
		glog.Errorf("moving %q to %q: %q", path, target, err)
	}
}

func (q *folderQueue) markComplete(path string) {
	if err := os.Remove(path); err != nil {
		glog.Errorf("removing file %q: %q", path, err)
	}
}
