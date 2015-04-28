# Mail Queue Dispatcher

The Mail Queue Dispatcher is a simple Windows service that watches a 
mailqueue folder, looks up sender information, and transmits an email
for each message it finds in the folder.

If there is a problem, the message will be moved to the configured 
badmail folder.


## Usage

To get the program, just run `go get github.com/johnweldon/mqd`

To install: `go install github.com/johnweldon/mqd/cmd/smtp-dispatcher`
This will install the binary in your GOPATH, but to install the windows
service, I recommend copying the executable into its own folder and then
creating or generating the `.smtp-dispatcher.settings` file, and
modifying it to match your settings first. Then to install the service
run `./smtp-dispatcher.exe install`, and `./smtp-dispatcher.exe start`
to start monitoring the mailqueue folder and sending emails.


![gopher mascot](img/smtp-dispatcher-gopher.png)
