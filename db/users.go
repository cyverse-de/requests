package db

import (
	"context"
	"database/sql"
	"fmt"
)

// GetUserID obtains the internal user ID  for the given username and user domain from the DE database.
func GetUserID(ctx context.Context, tx *sql.Tx, username, userDomain string) (string, error) {
	query := "SELECT id FROM users WHERE username = $1"
	qualifiedUsername := fmt.Sprintf("%s@%s", username, userDomain)

	// Query the database.
	row := tx.QueryRowContext(ctx, query, qualifiedUsername)

	// extract the user ID if at least one row was found.
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
