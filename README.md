# GMail / Google Calendar CLI client for event tunneling 

## Installation

```
$ go get -u github.com/greymd/gctunnel
```

## Usage

```
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
  authorize     credential.json file is required. Visit Ref[1] and download under "OAuth 2.0 Client ID" and rename the file
  profile       Show authorized account information
  messages      List messages
  message       Show body of single message
  calendars     List Calendars
  events        List events, starting within last month by default, on specified calendar
  create-event  Create new event on the specified calendar

Arguments:
  <QUERY>        Query for searching GMail messages. See Ref[2]. [default: in:inbox]
  <MESSAGE_ID>   Identifier of the message in GMail. Value of "id" in results of 'messages'.
  <CALENDAR_ID>  Identifier of the calendar. Value of "id" in results of 'calenders'. [default: (authorized GMail address)]
  <COLOR_ID>     Specify color of the event. Value of "colorID" in results of "events".
  <text>         Free format text.
  <datetime>     A combined date-time value formatted according to RFC3339, e.g "2020-04-23T00:00:00Z"
  <tz>           Time zone name with IANA Time Zone Database name, e.g "Asia/Tokyo". [default: UTC]

Ref:
  [1] https://console.developers.google.com/apis/credentials
  [2] https://support.google.com/mail/answer/7190
```

## Getting Started

(1) Go to https://console.developers.google.com/ and create new project

(2) Create new OAuth 2.0 Client ID and download JSON file including client ID / secret. Rename it to `credential.json`

(3) Run `gctunnel authorize` and follows the introduction
```
$ gctunnel authorize
Go to the following link in your browser then type the authorization code:
https://accounts.google.com/o/oauth2/auth?access_type=...
!! Paste your credentials provided by the browser
```
=> `token.json` file will be created on the current directory.

(4) Run commands

```
$ gctunnel profile
=> Your email address will be shown.

$ gctunnel messages
=> Your emails in "inbox" will be shown.
```

Please make sure `credential.json` and `token.json` are need to be located on the current directory.

## Purpose of this tool
Convert date/time in GMail message into Google Calender's event!

```bash
$ gctunnel messages 'from:noreply@example.com subject:Appointment newer_than:1d' \
  | head -n 1 | jq .id \
  | xargs -I@ go run ./gctunnel.go message @ | jq -r .body
Thank you for your appointoment.
You appointment time is [2020/01/01 17:30].
```
=> Check GMail message includes date or time.

If this script is running once a day, you won't miss your important appointoments!
```bash
#!/bin/bash
eventDate=$(gctunnel messages 'from:noreply@example.com subject:Appointment newer_than:1d' \
  | head -n 1 | jq .id \
  | xargs -I@ go run ./gctunnel.go message @ | jq -r .body \
  | grep -o '\[[^]]*\]' | tr -d '[]' | date -f- --rfc-3339=seconds)

gtunnel create-event --summary="Appointment" --start="$eventDate" --end="$(date -d "$eventDate 30 minutes" --rfc-3339=seconds)"
```
=> New event titled "Appointment" will be created on your Google Calendar.


## Why I made it ?
I know there are several CLI tools providing GMail or Google Calendar functionalities.

i.e:
* [insanum/gcalcli](https://github.com/insanum/gcalcli)
* [ThomasHabets/cmdg](https://github.com/ThomasHabets/cmdg)
* [yoshinari-nomura/glima](https://github.com/yoshinari-nomura/glima)

Technically, the purpose of this tool can be done by combining above tools.

However, combining multiple tools may cause many troubles in the future.
For example, we have to authorize multiple applications and requires multiple Client ID / Secret.

And also, I want to avoid combining multiple tools developed by individuals or multiple language like Ruby, Python, Node.js like above.
Because even one of them loses compatibility or stops maintenance, entire script will stop to work.

Unfortunately, as far as I searched, there is no CLI tool providing entire functionalities of GMail/Google Calendar both.
