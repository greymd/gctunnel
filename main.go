package main

import (
	"os"
	"time"

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
  gctunnel msgs [<QUERY>]
  gctunnel msg <MSG_ID>
  gctunnel cals
  gctunnel evts <CAL_ID> [--since=<datetime>]
  gctunnel new-evt <CAL_ID> --summary=<text> --start=<datetime> --end=<datetime> [--color=<COL_ID>] [--description=<text>] [--timezone=<tz>]

Options:
  -h, --help      Show help.
  -V, --version   Show version

Commands:
  auth        credentials.json file is required. Visit Ref[1] and download under [OAuth 2.0 Client ID] and rename the file
  msgs        Search messages
  msg         Show body of single message
  cals        List Calendars
  evts        List events, starting within last month by default, on specified calendar
  new-evt     Create new event on the specified calendar

Arguments:
  <QUERY>        Query for searching Gmail messages. See Ref[2]. [default: in:inbox]
  <MSG_ID>       Identifier of the message in Gmail.
  <CAL_ID>       Identifier of the calendar.
  <text>         Free format text.
  <datetime>     A combined date-time value formatted according to RFC3339, e.g "2020-04-23T00:00:00Z"
  <tz>           Time zone name with IANA Time Zone Database name, e.g "Asia/Tokyo". [default: UTC]

Ref:
  [1] https://console.developers.google.com/apis/credentials
  [2] https://support.google.com/mail/answer/7190
`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["auth"].(bool) {
		auth.RefreshToken()
		auth.GetClient()
		os.Exit(0)
	}
	// fmt.Println(args)
	// os.Exit(0)
	client := auth.GetClient()
	if args["msgs"].(bool) {
		query := args["<QUERY>"]
		if query == nil {
			query = "in:inbox"
		}
		messages.List(client, query.(string))
	} else if args["msg"].(bool) {
		messages.Get(client, args["<MSG_ID>"].(string))
	} else if args["cals"].(bool) {
		events.ListCalendars(client)
	} else if args["evts"].(bool) {
		since := args["--since"]
		if since == nil {
			since = getLastMonth()
		}
		events.ListEvents(client, args["<CAL_ID>"].(string), since.(string))
	} else if args["new-evt"].(bool) {
		description := args["--description"]
		if description == nil {
			description = ""
		}
		timezone := args["--timezone"]
		if timezone == nil {
			timezone = "UTC"
		}
		colorID := args["--color"]
		if colorID == nil {
			colorID = ""
		}
		events.CreateEvent(
			client,
			args["<CAL_ID>"].(string),
			args["--summary"].(string),
			description.(string),
			args["--start"].(string),
			args["--end"].(string),
			timezone.(string),
			colorID.(string))
	}
	os.Exit(0)
}

func getLastMonth() string {
	now := time.Now()
	then := now.AddDate(0, -1, 0)
	return then.Format(time.RFC3339)
}
