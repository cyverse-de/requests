package model

import "time"

// RequestSubmission represents a new request being submitted to this service.
type RequestSubmission struct {
	// The request type
	RequestType string `json:"request_type" validate:"required"`

	// Arbitrary JSON object describing the request details
	Details any `json:"details" validate:"required"`
}

// RequestSummary represents a brief overview of an administrative request.
type RequestSummary struct {
	// The request ID
	ID string `json:"id"`

	// The username of the requesting user
	RequestingUser string `json:"requesting_user"`

	// The request type
	RequestType string `json:"request_type"`

	// The time the request was submitted
	CreatedDate time.Time `json:"created_date"`

	// The request status
	Status string `json:"status"`

	// The time the request was most recently updated
	UpdatedDate time.Time `json:"updated_date"`

	// Arbitrary JSON object describing the request details
	Details any `json:"details"`
}

// RequestListing represents a list of requests.
type RequestListing struct {
	// The list of requests
	Requests []*RequestSummary `json:"requests"`
}

// RequestUpdate represents a request status update.
type RequestUpdate struct {
	// The request update ID
	ID string `json:"id"`

	// The request status code
	StatusCode string `json:"status"`

	// The username of the updating user
	UpdatingUser string `json:"updating_user"`

	// The timestamp corresponding to when the request was updated
	CreatedDate time.Time `json:"created_date"`

	// The message that was entered when the update was created
	Message string `json:"message"`
}

// RequestDetails represents the details of a request.
type RequestDetails struct {
	// The request ID
	ID string `json:"id"`

	// The username of the requesting user
	RequestingUser string `json:"requesting_user"`

	// The request type
	RequestType string `json:"request_type"`

	// Arbitrary JSON object describing the request details
	Details any `json:"details"`

	// The status updates for this request
	Updates []*RequestUpdate `json:"updates"`
}

// RequestUpdateSubmission represents information that should be submitted with a request update.
type RequestUpdateSubmission struct {
	// The request status code
	StatusCode string `json:"status" validate:"required"`

	// The message to associate with the request update
	Message string `json:"message"`
}
