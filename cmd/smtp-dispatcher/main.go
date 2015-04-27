package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/johnweldon/mqd/config"
	"github.com/johnweldon/mqd/dispatcher"
	"github.com/johnweldon/mqd/mailer"
)

const (
	settingsfile = ".smtp-dispatcher.settings"
)

func main() {
	settings, err := config.ReadSettings(settingsfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: couldn't read settings %q\n", err)
		os.Exit(-1)
	}
	q := dispatcher.NewPickupFolderQueue(settings.MailQueue, settings.BadMail)
	m := mailer.NewMailer(settings)

	loop(q, m)
}

func loop(q dispatcher.MailQueueDispatcher, m mailer.Mailer) {
	for {
		select {
		case <-time.After(15 * time.Second):
			q.Process(m.ConvertAndSend)
		}
	}
}
