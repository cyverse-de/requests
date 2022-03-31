package iplantgroups

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cyverse-de/requests/clients/util"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

// Subject represents a single subject returned by iplant-groups.
type Subject struct {
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	LastName    *string `json:"last_name"`
	Description *string `json:"description"`
	ID          *string `json:"id"`
	Institution *string `json:"institution"`
	FirstName   *string `json:"first_name"`
	SourceID    *string `json:"source_id"`
}

// Client describes a single instance of this client library.
type Client struct {
	baseURL       string
	deGrouperUser string
}

// NewClient creates a new instance of this client library.
func NewClient(baseURL, deGrouperUser string) *Client {
	return &Client{
		baseURL:       baseURL,
		deGrouperUser: deGrouperUser,
	}
}

// buildURL builds the URL to use for the given path components.
func (c *Client) buildURL(pathComponents ...string) (string, error) {
	var err error

	// Build the URL with the full path.
	fullURL, err := util.BuildURL(c.baseURL, pathComponents)
	if err != nil {
		return "", err
	}

	// Add the user query argument.
	query := fullURL.Query()
	query.Set("user", c.deGrouperUser)
	fullURL.RawQuery = query.Encode()

	// Return the updated URL.
	return fullURL.String(), nil
}

// GetUserInfo looks up information for a single user by calling iplant-groups.
func (c *Client) GetUserInfo(ctx context.Context, username string) (*Subject, error) {
	errorMessage := fmt.Sprintf("unable to look up userinformation for %s", username)
	var err error

	// Build the request URL.
	requestURL, err := c.buildURL("subjects", username)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}

	// Submit the request.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}
	defer resp.Body.Close()

	// Check the HTTP status code.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, errorMessage)
		}
		return nil, fmt.Errorf("%s: %s", errorMessage, respBody)
	}

	// Parse the response body.
	var subject Subject
	err = json.NewDecoder(resp.Body).Decode(&subject)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}

	return &subject, nil
}
