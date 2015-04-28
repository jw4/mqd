package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/winsvc/svc"

	"github.com/golang/glog"

	config "github.com/johnweldon/mqd/config"
)

const (
	svcName        = "MailQueueDispatcher"
	svcFriendly    = "Mail Queue Dispatcher"
	svcDescription = "Mail Queue Dispatcher watches a mailqueue folder and uses configured smtp connections " +
		"to relay the mail on.  If there are problems with delivery, or if the sender does not have a configured " +
		"smtp connection, the mail file will be moved to a badmail folder for manual review"
)

var (
	settingsfile     = ".smtp-dispatcher.settings"
	generateSettings = false
)

func init() {
	flag.BoolVar(&generateSettings, "generate", generateSettings, "generate settings file")
	flag.BoolVar(&generateSettings, "g", generateSettings, "generate settings file (short version)")
}

func main() {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		glog.Fatalf("couldn't find program folder %q\n", err)
		os.Exit(1)
	}
	settingsfile = filepath.Join(pwd, settingsfile)
	flag.Set("log_dir", pwd)
	flag.Parse()

	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		glog.Fatalf("couldn't tell if we're in an interactive session: %v", err)
		os.Exit(2)
	}

	if !interactive {
		runService(svcName, false)
		return
	}

	if generateSettings {
		generate()
		return
	}

	if len(flag.Args()) < 1 {
		usage("no command specified")
	}

	cmd := strings.ToLower(flag.Arg(0))
	switch cmd {
	case "debug":
		runService(svcName, true)
		return
	case "install":
		err = installService(svcName, svcFriendly, svcDescription)
	case "uninstall", "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		glog.Fatalf("failed to %s %s: %v", cmd, svcName, err)
		os.Exit(3)
	}
	return
}

func usage(message string) {
	fmt.Fprintf(os.Stderr, "\n%s\n\n"+
		"usage: %s <command>\n"+
		"    where <command> is one of\n"+
		"    install, remove, debug, start, stop, pause, or continue.\n\n",
		message, os.Args[0])
	os.Exit(8)
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
		os.Exit(9)
	}
	os.Exit(0)
}
