// Copyright 2015-2017 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package dispatcher // import "jw4.us/mqd/dispatcher"

// MailQueueCallbackFn describes the callback mechanism the dispatcher
// uses to transmit raw bytes representing an email to the mailer to
// actually send.
type MailQueueCallbackFn func([]byte) bool

// MailQueueDispatcher describes the interface a dispatcher must
// fulfill in order to use the MailQueueCallbackFn
type MailQueueDispatcher interface {
	Process(MailQueueCallbackFn) error
}
