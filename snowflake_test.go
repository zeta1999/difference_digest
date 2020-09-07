// +build snowflake

package difference_digest_test

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"testing"

	_ "github.com/snowflakedb/gosnowflake"

	"github.com/hundredwatt/difference_digest"
	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"
)

func TestSnowflake(t *testing.T) {
	db := connectSnowflake()
	defer func() {
		err := difference_digest.SnowflakeCleanup(db)
		if err != nil {
			fmt.Println(err)
		}
		db.Close()
	}()

	difference_digest.DatabaseType = difference_digest.Snowflake

	db.Exec("CREATE TEMP TABLE sourcethings (id bigint)")
	db.Exec("CREATE TEMP TABLE sinkthings (id bigint)")
	db.Exec("INSERT INTO sourcethings SELECT SEQ4() + 1 FROM table(generator(rowcount => 9900)) v") // 100 items missing
	db.Exec("INSERT INTO sinkthings SELECT SEQ4() + 1 FROM table(generator(rowcount => 10000)) v")

	// 1. Use the Strata Estimator to get the approximate number of differences between the two tables
	sourceEstimator, err := difference_digest.EncodeEstimatorDB(db, "sourcethings", "id")
	assert.NoError(t, err)
	sinkEstimator, err := difference_digest.EncodeEstimatorDB(db, "sinkthings", "id")
	assert.NoError(t, err)
	estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

	// 2. Get an IBF of the appropriate size from each source
	alpha := 5.0
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	sourceIBF, _ := difference_digest.EncodeIBFDB(cells, db, "sourcethings", "id")
	sinkIBF, _ := difference_digest.EncodeIBFDB(cells, db, "sinkthings", "id")

	// 3. Compute the difference of the IBFs
	diff := sinkIBF.Subtract(sourceIBF)
	sinkWithoutSource, sourceWithoutSink, ok := diff.Decode()
	assert.True(t, ok)
	assert.Len(t, sinkWithoutSource, 100)
	assert.Empty(t, sourceWithoutSink)
}

func connectSnowflake() *sql.DB {
	db, err := dburl.Open(os.ExpandEnv("snowflake://$TEST_SNOWFLAKE_USER:$TEST_SNOWFLAKE_PASSWORD@$TEST_SNOWFLAKE_HOST/$TEST_SNOWFLAKE_DBNAME"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE SCHEMA PUBLIC")
	if err != nil {
		panic(err)
	}

	err = difference_digest.SnowflakeSetup(db)
	if err != nil {
		panic(err)
	}

	return db
}
