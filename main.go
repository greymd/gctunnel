package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/greymd/gctunnel/auth"
	"github.com/greymd/gctunnel/events"
	"github.com/greymd/gctunnel/messages"
)

var appVersion = `0.0.0`

var usage = `
Usage:
  gctunnel [Options]
  gctunnel auth
  gctunnel msgs <QUERY>
  gctunnel msg <MSGID>
  gctunnel cals
  gctunnel events --calendar-id <CALENDER_ID>
  gctunnel create-event --calendar-id <CALENDER_ID> --summary <SUMMARY> --start <START_DATE_TIME> --end <END_DATE_TIME> [--description <DESCRIPTION>] [--timezone <TZ>]

Options:
  -h, --help      Show help.
  -V, --version   Show version

Default values:
  DESCRIPTION     Empty
  TZ              UTC
`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["auth"].(bool) {
		auth.GetClient()
		os.Exit(0)
	}
	fmt.Println(args)
	client := auth.GetClient()
	if args["msgs"].(bool) {
		messages.List(client, args["<QUERY>"].(string))
	} else if args["msg"].(bool) {
		messages.Get(client, args["<MSGID>"].(string))
	} else if args["cals"] == true {
		events.ListCalendars(client)
	} else if args["events"].(bool) {
		events.ListEvents(client, args["<CALENDER_ID>"].(string))
	} else if args["create-event"].(bool) {
		fmt.Println(args)
		description := args["<DESCRIPTION>"]
		if description == nil {
			description = ""
		}
		timezone := args["<TZ>"]
		if timezone == nil {
			timezone = "UTC"
		}
		events.CreateEvent(client, args["<CALENDER_ID>"].(string), args["<SUMMARY>"].(string), description.(string), args["<START_DATE_TIME>"].(string), args["<END_DATE_TIME>"].(string), timezone.(string))
	}
	os.Exit(0)
}
