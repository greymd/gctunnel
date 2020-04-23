package auth

import (
	"fmt"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

var userID = "me"

// GetProfile ... Print profile on stdout
func GetProfile(client *http.Client) (*gmail.Profile, error) {
	srv, err := gmail.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Gmail service: %v", err)
	}
	res, err := srv.Users.GetProfile(userID).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrive user profile: %v", err)
	}
	return res, nil
}
