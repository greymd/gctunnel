package events

import (
	"fmt"
	"net/http"

	"google.golang.org/api/calendar/v3"
)

// CreateEvent ... Create event
func CreateEvent(client *http.Client, calendarID string, summary string, description string, startDateTime string, endDateTime string, timezone string, colorID string) error {
	srv, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Calendar client: %v", err)
	}

	event := &calendar.Event{
		Summary:     summary,
		Location:    "",
		Description: description,
		ColorId:     colorID,
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
		return fmt.Errorf("Unable to create event: %v", err)
	}
	fmt.Printf("Event created: %s\n", fixedEvent.HtmlLink)
	return nil
}

// ListCalendars ... list calenders
func ListCalendars(client *http.Client) error {
	svc, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("Unable to create Calendar service: %v", err)
	}
	listRes, err := svc.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve list of calendars: %v", err)
	}

	for _, v := range listRes.Items {
		cal, _ := svc.Calendars.Get(v.Id).Do()
		j, _ := cal.MarshalJSON()
		fmt.Println(string(j))
	}
	return nil
}

// ListEvents ... List events
func ListEvents(client *http.Client, calendarID string, since string, end string) error {
	svc, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("Unable to create Calendar service: %v", err)
	}
	pageToken := ""
	for {
		req := svc.Events.List(calendarID).Fields("items(id,colorId,start,end,summary)", "nextPageToken").TimeMin(since).TimeMax(end).SingleEvents(true)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		res, err := req.Do()
		if err != nil {
			return fmt.Errorf("Unable to retrieve calendar events list: %v", err)
		}
		for _, v := range res.Items {
			j, _ := v.MarshalJSON()
			fmt.Println(string(j))
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return nil
}

// GetEvent ... Get events
func GetEvent(client *http.Client, calendarID string, eventID string) error {
	svc, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("Unable to create Calendar service: %v", err)
	}
	res, err := svc.Events.Get(calendarID, eventID).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve event: %v", err)
	}
	j, _ := res.MarshalJSON()
	fmt.Println(string(j))
	return nil
}
