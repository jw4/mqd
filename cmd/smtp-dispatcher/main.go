package main

import (
	"fmt"
	"os"
	"time"

	"github.com/johnweldon/mqd/dispatcher"
	"github.com/johnweldon/mqd/mailer"
)

const (
	// TODO(jw4) fix
	mailqueuefolder = "C:\\build\\temp\\mailqueue"
	settingsfile    = ".smtp-dispatcher.settings"
)

func main() {
	q := dispatcher.NewPickupFolderQueue(mailqueuefolder)
	m := mailer.NewMailer()
	err := m.LoadSettings(mailer.ReadSettingsFile(settingsfile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: couldn't read settings %q\n", err)
		os.Exit(-1)
	}

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
