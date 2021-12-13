package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/cyverse-de/requests/clients/iplantemail"
	"github.com/cyverse-de/requests/clients/iplantgroups"
	"github.com/cyverse-de/requests/clients/notificationagent"
	"github.com/labstack/echo/v4"
)

// API defines REST API of the requests service.
type API struct {
	Echo                    *echo.Echo
	Title                   string
	Version                 string
	DB                      *sql.DB
	UserDomain              string
	AdminEmail              string
	IPlantEmailClient       *iplantemail.Client
	IPlantGroupsClient      *iplantgroups.Client
	NotificationAgentClient *notificationagent.Client
}

// RootResponse describes the response of the root endpoint.
type RootResponse struct {

	// The name of the service
	Service string `json:"service"`

	// The service title
	Title string `json:"title"`

	// The service version
	Version string `json:"Version"`
}

// ErrorResponse describes an error response for any endpoint. This type implements the error interface so that it can
// be returned as an error from existing functions.
type ErrorResponse struct {
	Message   string                  `json:"message"`
	ErrorCode string                  `json:"error_code,omitempty"`
	Details   *map[string]interface{} `json:"details,omitempty"`
}

// ErrorBytes returns a byte-array representation of an ErrorResponse.
func (e ErrorResponse) ErrorBytes() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return make([]byte, 0)
	}
	return bytes
}

// Error returns a string representation of an ErrorResponse.
func (e ErrorResponse) Error() string {
	return string(e.ErrorBytes())
}

// RootHandler handles GET requests to the / endpoint.
func (a API) RootHandler(ctx echo.Context) error {
	resp := RootResponse{
		Service: "requests",
		Title:   a.Title,
		Version: a.Version,
	}
	return ctx.JSON(http.StatusOK, resp)
}
