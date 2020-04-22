package messages

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

// Get ... Get single message
func Get(client *http.Client, id string) {
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	user := "me"
	msg, err := srv.Users.Messages.Get(user, id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve message %v: %v", id, err)
	}
	date := ""
	for _, h := range msg.Payload.Headers {
		if h.Name == "Date" {
			date = h.Value
			break
		}
	}
	fmt.Printf("ID: %v, Date: %v, Snippet: %q\n", msg.Id, date, msg.Snippet)
}

// List ... Get multiple messages
func List(client *http.Client, query string) {
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	user := "me"
	pageToken := ""
	for {
		// https://support.google.com/mail/answer/7190
		// i.e in:inbox
		req := srv.Users.Messages.List(user).Q(query)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}

		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {
			msg, err := srv.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}
			date := ""
			for _, h := range msg.Payload.Headers {
				if h.Name == "Date" {
					date = h.Value
					break
				}
			}
			fmt.Printf("ID: %v, Date: %v, Snippet: %q\n", msg.Id, date, msg.Snippet)
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
}
