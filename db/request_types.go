package db

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/cyverse-de/requests/model"
)

// requestTypesFromRows converts SQL rows to an array of request types. The column order is always expected to be ID,
// and name.
func requestTypesFromRows(rows *sql.Rows) ([]*model.RequestType, error) {
	requestTypes := make([]*model.RequestType, 0)

	// Build the list of request types.
	for rows.Next() {
		var rt model.RequestType
		err := rows.Scan(&rt.ID, &rt.Name, &rt.MaximumRequestsPerUser, &rt.MaximumConcurrentRequestsPerUser)
		if err != nil {
			return nil, err
		}
		requestTypes = append(requestTypes, &rt)
	}

	return requestTypes, nil
}

// ListRequestTypes returns a listing of request types from the database sorted by name.
func ListRequestTypes(ctx context.Context, tx *sql.Tx) ([]*model.RequestType, error) {
	query := `SELECT id, name, maximum_requests_per_user, maximum_concurrent_requests_per_user
	          FROM request_types
			  ORDER BY name`

	// Query the database.
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Return the list of request types.
	return requestTypesFromRows(rows)
}

// GetRequestType returns the request type with the given name if it exists.
func GetRequestType(ctx context.Context, tx *sql.Tx, name string) (*model.RequestType, error) {
	query := `SELECT id, name, maximum_requests_per_user, maximum_concurrent_requests_per_user
	          FROM request_types
			  WHERE name = $1`

	// Query the database.
	rows, err := tx.QueryContext(ctx, query, name)
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
func AddRequestType(ctx context.Context, tx *sql.Tx, name string, maximumRequestsPerUser, maximumConcurrentRequestsPerUser *int32) (
	*model.RequestType, error,
) {
	query := `INSERT INTO request_types (name, maximum_requests_per_user, maximum_concurrent_requests_per_user)
			  VALUES ($1, $2, $3)
			  RETURNING id, name, maximum_requests_per_user, maximum_concurrent_requests_per_user`

	// Insert the new request type.
	rows, err := tx.QueryContext(ctx, query, name, maximumRequestsPerUser, maximumConcurrentRequestsPerUser)
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

// UpdateRequestType updates the request type with the given name.
func UpdateRequestType(ctx context.Context, tx *sql.Tx, name string, maximumRequestsPerUser, maximumConcurrentRequestsPerUser *int32) (
	*model.RequestType, error,
) {

	// Build the query.
	builder := psql.Update("request_types").
		Where(sq.Eq{"name": name}).
		Suffix("RETURNING id, name, maximum_requests_per_user, maximum_concurrent_requests_per_user")
	if maximumRequestsPerUser != nil && *maximumRequestsPerUser >= 0 {
		builder = builder.Set("maximum_requests_per_user", *maximumRequestsPerUser)
	}
	if maximumConcurrentRequestsPerUser != nil && *maximumConcurrentRequestsPerUser >= 0 {
		builder = builder.Set("maximum_concurrent_requests_per_user", *maximumConcurrentRequestsPerUser)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	// Execute the statement.
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// Get the request type information.
	requestTypes, err := requestTypesFromRows(rows)
	if err != nil {
		return nil, err
	}

	// We should have a result.
	if len(requestTypes) == 0 {
		return nil, fmt.Errorf("unable to retrieve updated request type information")
	}
	return requestTypes[0], nil
}
