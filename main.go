package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/greymd/gctunnel/auth"
	"github.com/greymd/gctunnel/messages"
)

var appVersion = `0.0.0`

var usage = `
Usage:
  gctunnel auth
  gctunnel list-msgs <QUERY>
  gctunnel get-msg <MSGID>
  gctunnel register-event --subject <SUBJECT> --body <BODY> --time <TIME>

Options:
  -h, --help      Show help.
  -V, --version   Show version`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["auth"].(bool) {
		auth.GetClient()
		os.Exit(0)
	}
	client := auth.GetClient()
	if args["list-msgs"].(bool) {
		messages.List(client, args["<QUERY>"].(string))
	} else if args["get-msg"] == true {
		messages.Get(client, args["<MSGID>"].(string))
	} else if args["register-event"] == true {
		fmt.Println("this is register-event")
	}
	os.Exit(0)
}
