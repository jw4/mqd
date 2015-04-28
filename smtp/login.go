// Copyright 2015 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package smtp

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

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username: []byte(username), password: []byte(password)}
}

func (l *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	glog.Infof("Start(server: %+v)", *server)
	return "LOGIN", []byte{}, nil
}

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
