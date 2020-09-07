package difference_digest

import (
	"database/sql"
	"strings"
)

var postgresQueries = map[string]string{
	"ibf": `
		SELECT
			pg_temp.f_hash(idx, %[2]s) %% %[3]d AS cell,
			pg_temp.f_bit_xor(%[2]s::bigint) AS id_sum,
			pg_temp.f_bit_xor_numeric(pg_temp.f_hash(3 + 0, %[2]s)) AS hash_sum,
			Count(%[2]s) AS Count
		FROM (
			SELECT 0 AS idx, * FROM %[1]s UNION
			SELECT 1, * FROM %[1]s UNION
			SELECT 2, * FROM %[1]s
		) s
		GROUP BY 1
	`,
	"strata_estimator": `
		SELECT
			pg_temp.f_trailing_zeros(pg_temp.f_hash(3 + 1, %[2]s)) AS estimator,
			pg_temp.f_hash(idx, %[2]s) %% 80 AS cell,
			pg_temp.f_bit_xor(%[2]s::bigint) AS id_sum,
			pg_temp.f_bit_xor_numeric(pg_temp.f_hash(3 + 0, %[2]s)) AS hash_sum,
			Count(%[2]s) AS Count
		FROM (
			SELECT 0 AS idx, * FROM %[1]s UNION
			SELECT 1, * FROM %[1]s UNION
			SELECT 2, * FROM %[1]s
		) s
		GROUP BY 1, 2
	`,
}

// PostgresSetup creates the necessary aggregates and functions for running difference_digest SQL queries in PostgreSQL
func PostgresSetup(db *sql.DB) error {
	file, err := Asset("difference_digest_postgres.sql")

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
