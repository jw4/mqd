package main

import (
	"errors"
	"os"
	"path"

	"github.com/golang/glog"
	"gopkg.in/urfave/cli.v2"

	config "github.com/jw4/mqd/config"
)

var (
	settingsfile     = ".smtp-dispatcher.settings"
	generateSettings = false
)

func generate() {
	generateSettingsFile(`c:\`)
}

func generateSettingsFile(basePath string) {
	settings := &config.Settings{
		Interval:  47,
		MailQueue: path.Join(basePath, "mailqueue"),
		SentMail:  path.Join(basePath, "sentmail"),
		BadMail:   path.Join(basePath, "badmail"),
		C: map[string]config.ConnectionDetails{
			"foo@bar.com": {
				Sender:   "foo@bar.com",
				Server:   "localhost:587",
				Host:     "localhost",
				AuthType: config.PlainAuth,
				Username: "asdf",
				Password: "qwer",
			},
			"baz@foo.com": {
				Sender:   "baz@foo.com",
				Server:   "smtp.foo.com:587",
				Host:     "smtp.foo.com",
				AuthType: config.LoginAuth,
				Username: "baz@foo.com",
				Password: "pazzwerd",
			},
		},
	}

	if err := config.WriteSettings(settingsfile, settings); err != nil {
		glog.Fatalf("couldn't write settings %q\n", err)
		os.Exit(10)
	}
	os.Exit(0)
}

func generateFn(ctx *cli.Context) error {
	return errors.New("generateFn not implemented")
}

func generateCmd() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Usage:   "create a properly formatted settings file",
		Aliases: []string{"g"},
		Action:  generateFn}
}
