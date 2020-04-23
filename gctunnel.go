package main

import (
	"fmt"
	"log"
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
  gctunnel authorize
  gctunnel profile
  gctunnel messages [<QUERY>]
  gctunnel message <MESSAGE_ID>
  gctunnel calendars
  gctunnel events [<CALENDAR_ID>] [--since=<datetime>]
  gctunnel create-event [<CALENDAR_ID>] --summary=<text> --start=<datetime> --end=<datetime> [--color=<COLOR_ID>] [--description=<text>] [--timezone=<tz>]

Options:
  -h, --help      Show help.
  -V, --version   Show version

Commands:
  authorize     credentials.json file is required. Visit Ref[1] and download under "OAuth 2.0 Client ID" and rename the file
  profile       Show authorized account information
  messages      List messages
  message       Show body of single message
  calendars     List Calendars
  events        List events, starting within last month by default, on specified calendar
  create-event  Create new event on the specified calendar

Arguments:
  <QUERY>        Query for searching Gmail messages. See Ref[2]. [default: in:inbox]
  <MESSAGE_ID>   Identifier of the message in Gmail. Value of "id" in results of 'messages'.
  <CALENDAR_ID>  Identifier of the calendar. Value of "id" in results of 'calenders'. [default: (authorized GMail address)]
  <COLOR_ID>     Specify color of the event. Value of "colorID" in results of "events".
  <text>         Free format text.
  <datetime>     A combined date-time value formatted according to RFC3339, e.g "2020-04-23T00:00:00Z"
  <tz>           Time zone name with IANA Time Zone Database name, e.g "Asia/Tokyo". [default: UTC]

Ref:
  [1] https://console.developers.google.com/apis/credentials
  [2] https://support.google.com/mail/answer/7190
`

func main() {
	args, _ := docopt.ParseArgs(usage, nil, appVersion)
	if args["authorize"].(bool) {
		auth.RefreshToken()
		auth.GetClient()
		os.Exit(0)
	}

	client := auth.GetClient()

	if args["profile"].(bool) {
		profile, err := auth.GetProfile(client)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		j, _ := profile.MarshalJSON()
		fmt.Println(string(j))

	} else if args["messages"].(bool) {
		query := args["<QUERY>"]
		if query == nil {
			query = "in:inbox"
		}
		err := messages.List(client, query.(string))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else if args["message"].(bool) {
		err := messages.Get(client, args["<MESSAGE_ID>"].(string))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else if args["calendars"].(bool) {
		err := events.ListCalendars(client)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else if args["events"].(bool) {
		calendarID := args["<CALENDAR_ID>"]
		if calendarID == nil {
			profile, err := auth.GetProfile(client)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			calendarID = profile.EmailAddress
		}
		since := args["--since"]
		if since == nil {
			since = getLastMonth()
		}
		err := events.ListEvents(client, calendarID.(string), since.(string))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else if args["create-event"].(bool) {
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
		calendarID := args["<CALENDAR_ID>"]
		if calendarID == nil {
			profile, err := auth.GetProfile(client)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			calendarID = profile.EmailAddress
		}
		err := events.CreateEvent(
			client,
			calendarID.(string),
			args["--summary"].(string),
			description.(string),
			args["--start"].(string),
			args["--end"].(string),
			timezone.(string),
			colorID.(string))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func getLastMonth() string {
	now := time.Now()
	then := now.AddDate(0, -1, 0)
	return then.Format(time.RFC3339)
}
