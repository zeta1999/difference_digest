package difference_digest

import (
	"database/sql"
	"strings"
)

var snowflakeQueries = map[string]string{
	"ibf": `
		SELECT
			F_DD_HASH(idx, %[2]s) %% %[3]d AS cell,
			BIT_XOR_AGG(%[2]s) AS id_sum,
			BIT_XOR_AGG(F_DD_HASH(3 + 0, %[2]s)) AS hash_sum,
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
			F_DD_TRAILING_ZEROS(F_DD_HASH(3 + 1, %[2]s)) AS estimator,
			F_DD_HASH(idx, %[2]s) %% %[3]d AS cell,
			BIT_XOR_AGG(%[2]s) AS id_sum,
			BIT_XOR_AGG(f_dd_hash(3 + 0, %[2]s)) AS hash_sum,
			Count(%[2]s) AS Count
		FROM (
			SELECT 0 AS idx, * FROM %[1]s UNION
			SELECT 1, * FROM %[1]s UNION
			SELECT 2, * FROM %[1]s
		) s
		GROUP BY 1, 2
	`,
}

// SnowflakeSetup creates the necessary functions for running difference_digest SQL queries in Snowflake
func SnowflakeSetup(db *sql.DB) error {
	file, err := Asset("udfs/snowflake.sql")

	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";;")

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

// SnowflakeCleanup drops functions created in SnowflakeSetup()
func SnowflakeCleanup(db *sql.DB) error {
	queries := []string{
		"DROP FUNCTION F_DD_HASH(INT, BIGINT)",
		"DROP FUNCTION F_DD_TRAILING_ZEROS(FLOAT)",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
