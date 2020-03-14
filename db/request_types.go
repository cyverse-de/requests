package db

import (
	"database/sql"
	"fmt"

	"github.com/cyverse-de/requests/model"
)

func requestTypesFromRows(rows *sql.Rows) ([]*model.RequestType, error) {
	requestTypes := make([]*model.RequestType, 0)

	// Build the list of request types.
	for rows.Next() {
		var rt model.RequestType
		err := rows.Scan(&rt.ID, &rt.Name)
		if err != nil {
			return nil, err
		}
		requestTypes = append(requestTypes, &rt)
	}

	return requestTypes, nil
}

// ListRequestTypes returns a listing of request types from the database sorted by name.
func ListRequestTypes(tx *sql.Tx) ([]*model.RequestType, error) {
	query := "SELECT id, name FROM request_types ORDER BY name"

	// Query the database.
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the list of request types.
	return requestTypesFromRows(rows)
}

// GetRequestType returns the request type with the given name if it exists.
func GetRequestType(tx *sql.Tx, name string) (*model.RequestType, error) {
	query := "SELECT id, name FROM request_types WHERE name = $1"

	// Query the database.
	rows, err := tx.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the request type if one was found.
	requestTypes, err := requestTypesFromRows(rows)
	if err != nil {
		return nil, err
	}
	if len(requestTypes) > 0 {
		return requestTypes[0], nil
	}
	return nil, nil
}

// AddRequestType adds a request type with the given name.
func AddRequestType(tx *sql.Tx, name string) (*model.RequestType, error) {
	query := `INSERT INTO request_types (name)
			  VALUES ($1)
			  RETURNING id, name`

	// Insert the new request type.
	rows, err := tx.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the request type information.
	requestTypes, err := requestTypesFromRows(rows)
	if err != nil {
		return nil, err
	}

	// We should have a result.
	if len(requestTypes) == 0 {
		return nil, fmt.Errorf("unable to retrieve request type information after registration")
	}
	return requestTypes[0], nil
}
