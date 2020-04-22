package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

var appVersion = `0.0.0`

var usage = `
Usage:
  gctunnel auth
  gctunnel list-mails <INBOX>
  gctunnel get-mail <MSGID>
  gctunnel register-event --subject <SUBJECT> --body <BODY> --time <TIME>

Options:
  -h, --help      Show help.
  -V, --version   Show version`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["auth"] == true {
		fmt.Println("this is auth")
		client := auth.client()
	} else if args["list-mails"] == true {
		fmt.Println("this is list-mails")
	} else if args["get-mail"] == true {
		fmt.Println("this is get-mail")
	} else if args["register-event"] == true {
		fmt.Println("this is register-event")
	}
}
