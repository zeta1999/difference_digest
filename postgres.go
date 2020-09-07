package difference_digest

import (
	"database/sql"
	"strings"
)

// PostgresSetup creates the necessary aggregates and functions for running difference_digest SQL queries
func PostgresSetup(db *sql.DB) error {
	file, err := Asset("difference_digest.sql")

	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		if strings.TrimSpace(request) == "" {
			continue
		}
		_, err := db.Exec(request)
		if err != nil {
			return err
		}
	}

	return nil
}
