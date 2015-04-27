package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"

	config "github.com/johnweldon/mqd/config"
	"github.com/johnweldon/mqd/dispatcher"
	"github.com/johnweldon/mqd/mailer"
)

var (
	settingsfile = ".smtp-dispatcher.settings"
)

func init() {
	flag.StringVar(&settingsfile, "settingsfile", settingsfile, "json encoded settings file")
	flag.StringVar(&settingsfile, "sf", settingsfile, "json encoded settings file (short version)")
}

func printHelp() {
	fmt.Fprintf(os.Stdout, "Usage: %s [ run | generate ]\n\n run: run the mailqueue dispatcher\n generate: generate settings file\n\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		printHelp()
	}

	switch strings.ToLower(flag.Arg(0)) {
	case "run":
		runLoop()
	case "generate":
		generate()
	default:
		printHelp()
	}
}

func runLoop() {
	settings, err := config.ReadSettings(settingsfile)
	if err != nil {
		glog.Fatalf("couldn't read settings %q\n", err)
		os.Exit(2)
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

func generate() {
	settings := config.NewSettings("mailqueue_folder_path", "badmail_folder_path")
	settings.Connections["foo@bar.com"] = config.ConnectionDetails{
		AuthType: config.PlainAuth,
		Sender:   "foo@bar.com",
		Server:   "smtp.example.com:587",
		Host:     "smtp.example.com",
		Username: "username",
		Password: "P455w0rd",
	}

	if err := config.WriteSettings(settingsfile, settings); err != nil {
		glog.Fatalf("couldn't write settings %q\n", err)
		os.Exit(3)
	}
	os.Exit(0)
}
