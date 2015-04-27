package main

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"

	config "github.com/johnweldon/mqd/config"
	"github.com/johnweldon/mqd/dispatcher"
	"github.com/johnweldon/mqd/mailer"
)

const (
	settingsfile = ".smtp-dispatcher.settings"
)

func init() {}

func main() {
	flag.Parse()
	settings, err := config.ReadSettings(settingsfile)
	if err != nil {
		glog.Fatalf("couldn't read settings %q\n", err)
		os.Exit(-1)
	}
	q := dispatcher.NewPickupFolderQueue(settings.MailQueue, settings.BadMail)
	m := mailer.NewMailer(settings)

	loop(settings.Interval, q, m)
}

func loop(interval int, q dispatcher.MailQueueDispatcher, m mailer.Mailer) {
	q.Process(m.ConvertAndSend)
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			q.Process(m.ConvertAndSend)
		}
	}
}
