package messages

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

type messageObject struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Snippet string `json:"snippet"`
	Body    string `json:"body"`
}

var userID = "me"

func createMessage(msg *gmail.Message, beLong bool) messageObject {
	subject := ""
	to := ""
	from := ""
	date := ""
	body := ""
	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "Date":
			date = h.Value
			break
		case "To":
			to = h.Value
			break
		case "From":
			from = h.Value
			break
		case "Subject":
			subject = h.Value
			break
		}
	}
	if beLong {
		bodyEncoded := msg.Payload.Body.Data
		if bodyEncoded == "" {
			bodyEncoded = extractTextFromPart(msg.Payload.Parts)
		}
		bodyBinary, err := base64.URLEncoding.DecodeString(bodyEncoded)
		if err != nil {
			log.Fatalf("Unable to decode message body %v: %v", msg.Id, err)
		} else {
			body = string(bodyBinary)
		}
	}
	return messageObject{
		ID:      msg.Id,
		Date:    date,
		To:      to,
		From:    from,
		Subject: subject,
		Snippet: msg.Snippet,
		Body:    body,
	}
}

// Get body as text as much as possible.
func extractTextFromPart(parts []*gmail.MessagePart) string {
	for _, p := range parts {
		if p.MimeType == "text/html" {
			return p.Body.Data
		}
		if p.MimeType == "text/plain" {
			return p.Body.Data
		}
		if p.MimeType == "multipart/alternative" {
			return extractTextFromPart(p.Parts)
		}
	}
	return ""
}

func printMsg(msg messageObject) {
	j, _ := json.Marshal(msg)
	fmt.Println(string(j))
}

// Get ... Get single message
func Get(client *http.Client, id string) {
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	msg, err := srv.Users.Messages.Get(userID, id).Format("full").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve message %v: %v", id, err)
	}
	printMsg(createMessage(msg, true))
}

// List ... Get multiple messages
func List(client *http.Client, query string) {
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	pageToken := ""
	for {
		req := srv.Users.Messages.List(userID).Q(query)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}

		for _, m := range r.Messages {
			msg, err := srv.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}

			printMsg(createMessage(msg, false))
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
}
