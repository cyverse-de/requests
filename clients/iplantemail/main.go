package iplantemail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

// EmailRequestBody represents a request body sent to iplant-email.
type EmailRequestBody struct {
	To       string      `json:"to"`
	Template string      `json:"template"`
	Subject  string      `json:"subject"`
	Values   interface{} `json:"values"`
}

// Client describes a single instance of this client library.
type Client struct {
	baseURL string
}

// NewClient creates a new instance of this client library.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// SendEmail sends an arbitrary email.
func (c *Client) sendEmail(ctx context.Context, requestBody *EmailRequestBody) error {
	errorMessage := "unable to send email"
	var err error

	// Serialize the request body.
	body, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}

	req.Header.Set("content-type", "application/json")

	// Submit the request.
	resp, err := httpClient.Do(req)
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

// SendRequestSubmittedEmail sends an email corresponding to a request.
func (c *Client) SendRequestSubmittedEmail(ctx context.Context, emailAddress, templateName string, requestDetails interface{}) error {
	return c.sendEmail(ctx, &EmailRequestBody{
		To:       emailAddress,
		Template: templateName,
		Subject:  "New Administrative Request",
		Values:   requestDetails,
	})
}

// SendRequestUpdatedEmail sends an email corresponding to a request status update.
func (c *Client) SendRequestUpdatedEmail(ctx context.Context, emailAddress, templateName string, requestDetails interface{}) error {
	return c.sendEmail(ctx, &EmailRequestBody{
		To:       emailAddress,
		Template: templateName,
		Subject:  "Administrative Request Updated",
		Values:   requestDetails,
	})
}
