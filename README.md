# Mail Queue Dispatcher

[![GoDoc](https://godoc.org/github.com/jw4/mqd?status.svg)](https://godoc.org/github.com/jw4/mqd)

The Mail Queue Dispatcher is a simple Windows service that watches a
mailqueue folder, looks up sender information, and transmits an email
for each message it finds in the folder.

If there is a problem, the message will be moved to the configured
badmail folder.


## Usage

To get the program, just run `go get github.com/jw4/mqd`

To install: `go install github.com/jw4/mqd/cmd/smtp-dispatcher`
This will install the binary in your GOPATH, but to install the windows
service, I recommend copying the executable into its own folder and then
creating or generating the `.smtp-dispatcher.settings` file, and
modifying it to match your settings first. Then to install the service
run `./smtp-dispatcher.exe install`, and `./smtp-dispatcher.exe start`
to start monitoring the mailqueue folder and sending emails.


## Building

To generate the windows binary with the icon and resource info you can
use `go generate github.com/jw4/mqd/cmd/smtp-dispatcher` after
installing the fine tool by Joseph Spurrier:

  `go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo`

This should generate a .syso file which `go build` will use to
incorporate the resource info into the binary when you finish up with:

  `go build github.com/jw4/mqd/cmd/smtp-dispatcher`



![gopher mascot](img/smtp-dispatcher-gopher.png)

The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/) 
The design is licensed under the Creative Commons 3.0 Attributions license. 
Read this article for more details: http://blog.golang.org/gopher
