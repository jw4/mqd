// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package mailer

import (
	config "github.com/johnweldon/mqd/config"
)

type Mailer interface {
	LoadSettings(config.Settings) error
	Send(sender string, recipients []string, message []byte) error
	ConvertAndSend(email []byte) bool
}
