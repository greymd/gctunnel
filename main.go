package main

import (
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
  gctunnel msg <MSG_ID>
  gctunnel cals
  gctunnel evts <CAL_ID> [--since=<datetime>]
  gctunnel new-evt <CAL_ID> --summary=<text> --start=<datetime> --end=<datetime> [--description=<text>] [--timezone=<tz>]

Options:
  -h, --help      Show help.
  -V, --version   Show version

Commands:
  auth        Authentication Google API. Client ID/Secret are required. Visit Ref[1].
  msgs        Search messages.
  msg         Show body of single message.
  cals        List calendars.
  evts        List events, starting within last a week by default, on specified calendar.
  new-evt     Create new event on the specified calendar.

Arguments:
  <QUERY>        Query for searching Gmail messages e.g "in:inbox". See Ref[2].
  <MSG_ID>       Identifier of the message in Gmail. See Ref[3].
  <CAL_ID>       Identifier of the calendar. See Ref[4].
  <text>         Free format text.
  <datetime>     A combined date-time value formatted according to RFC3339, e.g "2020-04-23T00:00:00Z"
  <tz>           Time zone name with IANA Time Zone Database name, e.g "Asia/Tokyo". [default: UTC]

Ref:
  [1] https://console.developers.google.com
  [2] https://support.google.com/mail/answer/7190
  [3] https://developers.google.com/gmail/api/v1/reference/users/messages
  [4] https://developers.google.com/calendar/v3/reference/calendars
`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["auth"].(bool) {
		// Delete token.json
		auth.GetClient()
		os.Exit(0)
	}
	// fmt.Println(args)
	// os.Exit(0)
	client := auth.GetClient()
	if args["msgs"].(bool) {
		messages.List(client, args["<QUERY>"].(string))
	} else if args["msg"].(bool) {
		messages.Get(client, args["<MSG_ID>"].(string))
	} else if args["cals"].(bool) {
		events.ListCalendars(client)
	} else if args["evts"].(bool) {
		events.ListEvents(client, args["<CAL_ID>"].(string))
	} else if args["new-evt"].(bool) {
		description := args["--description"]
		if description == nil {
			description = ""
		}
		timezone := args["--timezone"]
		if timezone == nil {
			timezone = "UTC"
		}
		os.Exit(0)
		events.CreateEvent(client, args["<CAL_ID>"].(string), args["--summary"].(string), description.(string), args["--start"].(string), args["--end"].(string), timezone.(string))
	}
	os.Exit(0)
}
