package api

import (
	"net/http"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/labstack/echo/v4"
)

// GetRequestStatusCodesHandler handles GET requests to the /request-status-codes endpoint.
func (a *API) GetRequestStatusCodesHandler(c echo.Context) error {
	ctx := c.Request().Context()

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			c.Logger().Errorf("unable to roll back the transaction: %s", err)
		}
	}()

	// Obtain the list of request status codes.
	requestStatusCodes, err := db.ListRequestStatusCodes(ctx, tx)
	if err != nil {
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Return the response.
	return c.JSON(http.StatusOK, model.RequestStatusCodeListing{
		RequestStatusCodes: requestStatusCodes,
	})
}
