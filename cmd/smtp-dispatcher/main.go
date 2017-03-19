// Copyright 2015-2017 John Weldon. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// +build !windows

package main

import (
	"os"

	"github.com/golang/glog"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	if err := dispatcherApp().Run(os.Args); err != nil {
		glog.Fatalf("error running dispatcher app: %v", err)
		os.Exit(1)
	}
}

var (
	appVersion         = "0.1.1"
	dispatcherCommands = []*cli.Command{
		generateCmd(),
	}
)

func dispatcherApp() *cli.App {
	app := &cli.App{
		Name:     "MailQueueDispatcher",
		Commands: dispatcherCommands,
		Version:  appVersion,
	}
	return app
}
