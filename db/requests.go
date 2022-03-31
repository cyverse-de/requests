package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/cyverse-de/requests/model"

	sq "github.com/Masterminds/squirrel"

	"github.com/pkg/errors"
)

// CountRequestsOfType counts the number of requests of the given type that have been submitted by the given user.
func CountRequestsOfType(ctx context.Context, tx *sql.Tx, userID, requestTypeID string) (int32, error) {

	// Prepare the query.
	query, args, err := psql.Select("count(*)").
		From("requests").
		Where(sq.Eq{"requesting_user_id": userID}).
		Where(sq.Eq{"request_type_id": requestTypeID}).
		ToSql()
	if err != nil {
		return 0, err
	}

	// Query the database and extract the count.
	var count int32
	row := tx.QueryRowContext(ctx, query, args...)
	err = row.Scan(&count)
	return count, err
}

// CountActiveRequestsOfType counts the number of active requests of the given type that have been submitted by the
// given user.
func CountActiveRequestsOfType(ctx context.Context, tx *sql.Tx, userID, requestTypeID string) (int32, error) {

	// Prepare the primary query.
	subquery := psql.Select().
		Column("r.id").
		Column("last(rsc.name) OVER w AS status").
		From("requests r").
		Join("request_updates ru ON r.id = ru.request_id").
		Join("request_status_codes rsc ON ru.request_status_code_id = rsc.id").
		Where(sq.Eq{"requesting_user_id": userID}).
		Where(sq.Eq{"request_type_id": requestTypeID}).
		Suffix("WINDOW w AS (" +
			"PARTITION BY ru.request_id " +
			"ORDER BY ru.created_date " +
			"RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)")

	// Prepare the query.
	query, args, err := psql.Select("count(*)").
		FromSelect(subquery, "subquery").
		Where(sq.NotEq{"status": []string{"approved", "rejected"}}).
		ToSql()

	// Query the database and extract the count.
	var count int32
	row := tx.QueryRowContext(ctx, query, args...)
	err = row.Scan(&count)
	return count, err
}

// AddRequest adds a new request to the database.
func AddRequest(ctx context.Context, tx *sql.Tx, userID, requestTypeID string, details interface{}) (string, error) {
	query := `INSERT INTO requests (request_type_id, requesting_user_id, details)
			  VALUES ($1, $2, CAST($3 AS json))
			  RETURNING id`

	// Encode the request details.
	encodedDetails, err := json.Marshal(details)
	if err != nil {
		return "", errors.Wrap(err, "unable to JSON encode the request details")
	}

	// Insert the new request.
	row := tx.QueryRowContext(ctx, query, requestTypeID, userID, encodedDetails)

	// Extract the request ID.
	var requestID string
	err = row.Scan(&requestID)
	if err != nil {
		return "", err
	}

	return requestID, nil
}

// RequestListingOptions represents options that can be used to filter a request listing.
type RequestListingOptions struct {
	IncludeCompletedRequests bool
	RequestType              string
	RequestingUser           string
}

// GetRequestListing obtains a list of requests from the database.
func GetRequestListing(ctx context.Context, tx *sql.Tx, options *RequestListingOptions) ([]*model.RequestSummary, error) {

	// Prepare the primary listing query as a subquery.
	subquery := psql.Select().Distinct().
		Column("r.id").
		Column("regexp_replace(u.username, '@.*', '') AS username").
		Column("rt.name AS request_type").
		Column("first(ru.created_date) OVER w AS created_date").
		Column("last(rsc.name) OVER w AS status").
		Column("last(rsc.display_name) OVER w AS display_status").
		Column("last(ru.created_date) OVER w AS updated_date").
		Column("CAST(r.details AS text) AS details").
		From("requests r").
		Join("users u ON r.requesting_user_id = u.id").
		Join("request_types rt ON r.request_type_id = rt.id").
		Join("request_updates ru ON r.id = ru.request_id").
		Join("request_status_codes rsc ON ru.request_status_code_id = rsc.id").
		Suffix("WINDOW w AS (" +
			"PARTITION BY ru.request_id " +
			"ORDER BY ru.created_date " +
			"RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)")

	// Prepare the base query.
	base := psql.Select().
		Column("id").
		Column("username").
		Column("request_type").
		Column("created_date").
		Column("display_status").
		Column("updated_date").
		Column("details").
		FromSelect(subquery, "subquery").
		OrderBy("created_date")

	// Add the filter to omit completed requests if we're not supposed to include them in the listing.
	if !options.IncludeCompletedRequests {
		base = base.Where(sq.NotEq{"status": []string{"approved", "rejected"}})
	}

	// Add the filter to limit the listing to requests of a given type if applicable.
	if options.RequestType != "" {
		base = base.Where(sq.Eq{"request_type": options.RequestType})
	}

	// Add the filter to limit the listing to requests submitted by a user if applicable.
	if options.RequestingUser != "" {
		base = base.Where(sq.Eq{"username": options.RequestingUser})
	}

	// Build the query.
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}

	// Query the database.
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the listing from the result set.
	listing := make([]*model.RequestSummary, 0)
	for rows.Next() {
		var request model.RequestSummary
		var requestDetails string
		err = rows.Scan(
			&request.ID,
			&request.RequestingUser,
			&request.RequestType,
			&request.CreatedDate,
			&request.Status,
			&request.UpdatedDate,
			&requestDetails,
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(requestDetails), &request.Details)
		if err != nil {
			return nil, err
		}
		listing = append(listing, &request)
	}

	return listing, nil
}

// GetRequestDetails looks up the details of a request.
func GetRequestDetails(ctx context.Context, tx *sql.Tx, id string) (*model.RequestDetails, error) {
	query := `SELECT r.id, regexp_replace(u.username, '@.*', ''), rt.name, r.details
			  FROM requests r
			  JOIN users u ON r.requesting_user_id = u.id
			  JOIN request_types rt ON r.request_type_id = rt.id
			  WHERE r.id = $1`

	// Query the database.
	row := tx.QueryRowContext(ctx, query, id)

	// Extract the request details.
	var rd model.RequestDetails
	var rdDetails string
	err := row.Scan(&rd.ID, &rd.RequestingUser, &rd.RequestType, &rdDetails)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	err = json.Unmarshal([]byte(rdDetails), &rd.Details)
	if err != nil {
		return nil, err
	}

	// Add status information to the request details.
	rd.Updates, err = GetRequestStatusUpdates(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
