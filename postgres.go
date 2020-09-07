package difference_digest

import (
	"database/sql"
	"io/ioutil"
	"strings"
)

func PostgresSetup(db *sql.DB) error {
	file, err := ioutil.ReadFile("difference_digest.sql")

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
