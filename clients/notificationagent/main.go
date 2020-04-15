package notificationagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cyverse-de/requests/clients/util"
	"github.com/pkg/errors"
)

// Client describes a single instance of this client library.
type Client struct {
	baseURL string
}

// NewClient creates a new instance of this client library.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// NotificationRequest represents a request for a notification.
type NotificationRequest struct {
	Type          string      `json:"type"`
	User          string      `json:"user"`
	Subject       string      `json:"subject"`
	Message       string      `json:"message"`
	Email         bool        `json:"email"`
	EmailTemplate string      `json:"email_template"`
	Payload       interface{} `json:"payload"`
}

// buildURL builds the URL to use for the given path components.
func (c *Client) buildURL(pathComponents ...string) (string, error) {
	fullURL, err := util.BuildURL(c.baseURL, pathComponents)
	if err != nil {
		return "", err
	}
	return fullURL.String(), nil
}

// SendNotification sends a notification request to the notificaiton agent.
func (c *Client) SendNotification(requestBody *NotificationRequest) error {
	errorMessage := "unable to send notification"
	var err error

	// Build the request URL.
	requestURL, err := c.buildURL("notification")
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}

	// Serialize the request body.
	body, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}

	// Submit the request.
	resp, err := http.Post(requestURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}
	defer resp.Body.Close()

	// Check the HTTP Status code.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, errorMessage)
		}
		return fmt.Errorf("%s: %s", errorMessage, respBody)
	}

	return nil
}
