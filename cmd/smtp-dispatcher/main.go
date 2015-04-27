package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/golang/glog"

	config "github.com/johnweldon/mqd/config"
	"github.com/johnweldon/mqd/dispatcher"
	"github.com/johnweldon/mqd/mailer"
)

var (
	settingsfile     = ".smtp-dispatcher.settings"
	generateSettings = false
	settings         config.Settings
)

func init() {
	flag.StringVar(&settingsfile, "settingsfile", settingsfile, "json encoded settings file")
	flag.StringVar(&settingsfile, "sf", settingsfile, "json encoded settings file (short version)")

	flag.BoolVar(&generateSettings, "generate", generateSettings, "generate settings file")
	flag.BoolVar(&generateSettings, "g", generateSettings, "generate settings file (short version)")

	flag.Set("log_dir", ".")
}

func main() {
	flag.Parse()

	if generateSettings {
		generate()
		return
	}

	s, err := config.ReadSettings(settingsfile)
	if err != nil {
		glog.Fatalf("couldn't read settings %q\n", err)
		os.Exit(2)
	}

	settings = s
	runLoop()
	os.Exit(0)
}

func runLoop() {
	q := dispatcher.NewPickupFolderQueue(settings.MailQueue, settings.BadMail)
	m := mailer.NewMailer(settings)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	var t tomb.Tomb
	t.Go(func() error {
		q.Process(m.ConvertAndSend)
		for {
			select {
			case s := <-c:
				glog.Infof("got signal %v", s)
				glog.Flush()
				return nil
			case <-t.Dying():
				return nil
			case <-time.After(time.Duration(settings.Interval) * time.Second):
				q.Process(m.ConvertAndSend)
				glog.Flush()
			}
		}
	})
	t.Wait()
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
