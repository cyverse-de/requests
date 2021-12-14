package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/cyverse-de/requests/query"
	"github.com/labstack/echo/v4"
)

// validateRequestTypeName returns an error if a request type name is invalid.
func validateRequestTypeName(name string) error {
	re := regexp.MustCompile("^[\\w-]+$")
	if !re.MatchString(name) {
		return fmt.Errorf("request type names may only contain alphanumerics, underscores, and hyphens")
	}
	return nil
}

// GetRequestTypesHandler handles GET requests to the /request-types endpoint.
func (a *API) GetRequestTypesHandler(ctx echo.Context) error {

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Obtain the list of request types.
	requestTypes, err := db.ListRequestTypes(tx)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, model.RequestTypeListing{
		RequestTypes: requestTypes,
	})
}

// RegisterRequestTypeHandler handles POST requests to the /request-types/{name} endpoint.
func (a *API) RegisterRequestTypeHandler(ctx echo.Context) error {
	name := ctx.Param("name")

	// Validate the request type name.
	err := validateRequestTypeName(name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get the maximum number of requests per user for this request type and validate it.
	maximumRequestsPerUser, err := query.ValidateOptionalIntQueryParam(ctx, "maximum-requests-per-user")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}
	if maximumRequestsPerUser != nil && *maximumRequestsPerUser <= 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "maximum-requests-per-user must be a positive integer if specified",
		})
	}

	// Get the maximum number of concurrent requestws per user for this request type and validate it.
	maximumConcurrentRequestsPerUser, err := query.ValidateOptionalIntQueryParam(
		ctx, "maximum-concurrent-requests-per-user",
	)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}
	if maximumConcurrentRequestsPerUser != nil && *maximumConcurrentRequestsPerUser <= 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "maximum-concurrent-requests-per-user must be a positive integer if specified",
		})
	}

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// If a request type with the same name already exists, return it.
	requestType, err := db.GetRequestType(tx, name)
	if err != nil {
		return err
	}
	if requestType != nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, requestType)
	}

	// The request type doesn't exist yet. Add it.
	requestType, err = db.AddRequestType(tx, name, maximumRequestsPerUser, maximumConcurrentRequestsPerUser)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, requestType)
}

// UpdateRequestTypesHandler handles PATCH requests to the /request-types/{name} endpoint.
func (a *API) UpdateRequestTypesHandler(ctx echo.Context) error {
	name := ctx.Param("name")

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the request type.
	requestType, err := db.GetRequestType(tx, name)
	if err != nil {
		return err
	}
	if requestType == nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("request type %s not found", name),
		})
	}

	// Get the maximum number of requests per user for this request type and validate it.
	maximumRequestsPerUser, err := query.ValidateOptionalIntQueryParam(ctx, "maximum-requests-per-user")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}
	if maximumRequestsPerUser != nil && *maximumRequestsPerUser <= 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "maximum-requests-per-user must be a positive integer if specified",
		})
	}

	// Get the maximum number of concurrent requests per user for this request type and validate it.
	maximumConcurrentRequestsPerUser, err := query.ValidateOptionalIntQueryParam(
		ctx, "maximum-concurrent-requests-per-user",
	)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}
	if maximumConcurrentRequestsPerUser != nil && *maximumConcurrentRequestsPerUser <= 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "maximum-concurrent-requests-per-user must be a positive integer if specified",
		})
	}

	// Simply return the existing object if there were no updates.
	if maximumRequestsPerUser == nil && maximumConcurrentRequestsPerUser == nil {
		return ctx.JSON(http.StatusOK, requestType)
	}

	// Update the request type.
	requestType, err = db.UpdateRequestType(tx, name, maximumRequestsPerUser, maximumConcurrentRequestsPerUser)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, requestType)
}

// GetRequestTypeHandler handles GET requests to the /request-types/{name} endpoint.
func (a *API) GetRequestTypeHandler(ctx echo.Context) error {
	name := ctx.Param("name")

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the request type.
	requestType, err := db.GetRequestType(tx, name)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	if requestType == nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("request type %s not found", name),
		})
	}
	return ctx.JSON(http.StatusOK, requestType)
}
