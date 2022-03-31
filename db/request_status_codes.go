package db

import (
	"context"
	"database/sql"

	"github.com/cyverse-de/requests/model"
)

// requestStatusCodesFromRows converts SQL rows to an array of request status codes. The column order is always
// expected to be ID, name, display_name, and email_template.
func requestStatusCodesFromRows(rows *sql.Rows) ([]*model.RequestStatusCode, error) {
	requestStatusCodes := make([]*model.RequestStatusCode, 0)

	// Build the list of request status codes.
	for rows.Next() {
		var rsc model.RequestStatusCode
		err := rows.Scan(&rsc.ID, &rsc.Name, &rsc.DisplayName, &rsc.EmailTemplate)
		if err != nil {
			return nil, err
		}
		requestStatusCodes = append(requestStatusCodes, &rsc)
	}

	return requestStatusCodes, nil
}

// ListRequestStatusCodes lists all of the currently available request status codes.
func ListRequestStatusCodes(ctx context.Context, tx *sql.Tx) ([]*model.RequestStatusCode, error) {
	query := "SELECT id, name, display_name, email_template FROM request_status_codes"

	// Query the database.
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the list of rows.
	return requestStatusCodesFromRows(rows)
}

// GetRequestStatusCode looks the status code with the given name.
func GetRequestStatusCode(ctx context.Context, tx *sql.Tx, status string) (*model.RequestStatusCode, error) {
	query := "SELECT id, name, display_name, email_template FROM request_status_codes WHERE name = $1"

	// Query the database.
	rows, err := tx.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the request status code information.
	requestStatusCodes, err := requestStatusCodesFromRows(rows)
	if err != nil {
		return nil, err
	}

	// Return the matching status code if there is one.
	if len(requestStatusCodes) == 0 {
		return nil, nil
	}
	return requestStatusCodes[0], nil
}
