// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package dispatcher

type MailQueueCallbackFn func([]byte) bool

type MailQueueDispatcher interface {
	Process(MailQueueCallbackFn) error
}
