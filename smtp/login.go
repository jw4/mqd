package smtp

import (
	"fmt"
	"log"
	ns "net/smtp"
	"os"
	"strings"
)

var logger = log.New(os.Stderr, "dispatcher.smtp: ", log.Lshortfile)

type loginAuth struct {
	username []byte
	password []byte
}

func LoginAuth(username, password string) ns.Auth {
	return &loginAuth{username: []byte(username), password: []byte(password)}
}

func (l *loginAuth) Start(server *ns.ServerInfo) (string, []byte, error) {
	//logger.Printf("Start(server: %+v)", *server)
	return "LOGIN", []byte{}, nil
}

func (l *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	//logger.Printf("Next(fromServer: %q, more: %t)", string(fromServer), more)
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
