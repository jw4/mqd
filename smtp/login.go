// Copyright 2015-2017 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package smtp // import "jw4.us/mqd/smtp"

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/golang/glog"
)

type loginAuth struct {
	username []byte
	password []byte
}

// LoginAuth returns an smtp.Auth implementation for the LOGIN smtp
// mechanism.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username: []byte(username), password: []byte(password)}
}

// Start fulfills the smtp.Auth interface.  It returns information to
// identify it's capabilities.
func (l *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	glog.Infof("Start(server: %+v)", *server)
	return "LOGIN", []byte{}, nil
}

// Next fulfills the smtp.Auth interface.  It responds to inputs from
// the remote SMTP server until the authentication is complete or fails.
func (l *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	glog.Infof("Next(fromServer: %q, more: %t)", string(fromServer), more)
	response := strings.ToLower(string(fromServer[:9]))
	switch response {
	case "username:":
		return l.username, nil
	case "password:":
		return l.password, nil
	case "2.7.0 aut":
		return nil, nil
	}
	return nil, fmt.Errorf("unknown prompt: %q", fromServer)
}
