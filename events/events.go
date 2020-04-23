package events

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/calendar/v3"
)

// CreateEvent ... Create event
func CreateEvent(client *http.Client, calendarID string, summary string, description string, startDateTime string, endDateTime string, timezone string) {
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	event := &calendar.Event{
		Summary:     summary,
		Location:    "",
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startDateTime,
			TimeZone: timezone,
		},
		End: &calendar.EventDateTime{
			DateTime: endDateTime,
			TimeZone: timezone,
		},
	}

	fixedEvent, err := srv.Events.Insert(calendarID, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", fixedEvent.HtmlLink)
}

// ListCalendars ... list calenders
func ListCalendars(client *http.Client) {
	svc, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}
	listRes, err := svc.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}

	for _, v := range listRes.Items {
		cal, _ := svc.Calendars.Get(v.Id).Do()
		fmt.Printf("Calendar ID: %v, %v\n", v.Id, cal.Summary)
	}
}

// ListEvents ... List events
func ListEvents(client *http.Client, calendarID string, since string) {
	svc, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}
	pageToken := ""
	for {
		req := svc.Events.List(calendarID).Fields("items(start,end,summary)", "nextPageToken").TimeMin(since)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		res, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar events list: %v", err)
		}
		for _, v := range res.Items {
			fmt.Printf("Calendar ID %q start: %v:, end: %v, summary:%q\n", calendarID, v.Start, v.End, v.Summary)
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
}
