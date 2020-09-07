package difference_digest_test

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"

	"github.com/hundredwatt/difference_digest"
)

func TestIBFCell(t *testing.T) {
	s1 := uint64(12345)
	s2 := uint64(98765)
	cell := difference_digest.IBFCell{}

	assert.True(t, cell.IsZero())
	assert.False(t, cell.IsPure())

	cell.Insert(s1)

	assert.False(t, cell.IsZero())
	assert.True(t, cell.IsPure())

	cell.Insert(s2)

	assert.False(t, cell.IsZero())
	assert.False(t, cell.IsPure())

	cell2 := difference_digest.IBFCell{}
	cell2.Insert(s2)
	cell.Subtract(&cell2)

	assert.False(t, cell.IsZero())
	assert.True(t, cell.IsPure())
}

func TestDBComputeDifference(t *testing.T) {
	db := connectDB()
	defer db.Close()

	_, err := db.Exec("CREATE TEMP TABLE sourcethings (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("CREATE TEMP TABLE sinkthings (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO sourcethings (id) SELECT * from generate_series(1,9900)") // 100 items missing
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO sinkthings (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	// 1. Use the Strata Estimator to get the approximate number of differences
	sourceEstimator, err := difference_digest.EncodeEstimatorDB(db, "sourcethings", "id")
	assert.NoError(t, err)
	sinkEstimator, err := difference_digest.EncodeEstimatorDB(db, "sinkthings", "id")
	assert.NoError(t, err)

	estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

	assert.Less(t, int(estimatedDeletes), 150)
	assert.Greater(t, int(estimatedDeletes), 50)

	// 2. Get an IBF of the appropriate size from each source
	alpha := 5
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	sourceIBF, err := difference_digest.EncodeIBFDB(cells, db, "sourcethings", "id")
	assert.NoError(t, err)
	sinkIBF, err := difference_digest.EncodeIBFDB(cells, db, "sinkthings", "id")
	assert.NoError(t, err)

	// 3. Compute the difference of the IBFs
	diff := sinkIBF.Subtract(sourceIBF)

	sinkWithoutSource, sourceWithoutSink, ok := diff.Decode()
	assert.True(t, ok)
	assert.Len(t, sinkWithoutSource, 100)
	assert.Empty(t, sourceWithoutSink)

	assert.Contains(t, sinkWithoutSource, uint64(9901))
	assert.Contains(t, sinkWithoutSource, uint64(9997))
	assert.NotContains(t, sinkWithoutSource, uint64(2))
}

func TestDBInMemoryEstimatorEquivalence(t *testing.T) {
	db := connectDB()
	defer db.Close()

	_, err := db.Exec("CREATE TEMP TABLE things (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO things (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	dbEstimator, err := difference_digest.EncodeEstimatorDB(db, "things", "id")
	assert.NoError(t, err)

	set := makeSet(1, 10000)
	inMemoryEstimator := difference_digest.NewStrataEstimator()
	for k := range set {
		inMemoryEstimator.Add(k)
	}

	assert.Equal(t, 0, int(inMemoryEstimator.EstimateDifference(dbEstimator)))

	for i := range inMemoryEstimator.Stratum {
		diff := dbEstimator.Stratum[i].Subtract(&inMemoryEstimator.Stratum[i])
		aWithoutB, bWithoutA, ok := diff.Decode()
		assert.True(t, ok)
		assert.Empty(t, aWithoutB)
		assert.Empty(t, bWithoutA)
	}
}

func TestDBInMemoryIBFEquivalence(t *testing.T) {
	db := connectDB()
	defer db.Close()

	_, err := db.Exec("CREATE TEMP TABLE things (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO things (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	dbIBF, err := difference_digest.EncodeIBFDB(10, db, "things", "id")
	assert.NoError(t, err)

	set := makeSet(1, 10000)
	setIBF := difference_digest.NewIBF(10)
	for k := range set {
		setIBF.Add(k)
	}

	diff := dbIBF.Subtract(setIBF)

	aWithoutB, bWithoutA, ok := diff.Decode()
	assert.True(t, ok)
	assert.Empty(t, aWithoutB)
	assert.Empty(t, bWithoutA)
}

func TestInMemory(t *testing.T) {
	// TODO: make an example
	setSize := uint64(10 * 1000)
	deleteCount := 42

	sourceSet := makeSet(0, setSize)
	sinkSet := makeSet(0, setSize)

	deletes := make(map[uint64]bool)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < deleteCount; i++ {
		d := uint64(r1.Int63n(int64(setSize)))
		deletes[d] = true
		delete(sourceSet, d)
	}

	// Sink:
	// Sink computes an estimator and sends it to Sink
	sinkEstimator := difference_digest.NewStrataEstimator()
	for k := range sinkSet {
		sinkEstimator.Add(k)
	}

	// Source:
	// Source computes its own estimator, and then decodes it with sink's
	sourceEstimator := difference_digest.NewStrataEstimator()
	for k := range sourceSet {
		sourceEstimator.Add(k)
	}

	estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

	alpha := 5
	// TODO: function
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	// Source computes its IBF and sends it to sink
	sourceIBF := difference_digest.NewIBF(cells)
	for k := range sourceSet {
		sourceIBF.Add(k)
	}

	// Sink:
	// Sink computes its IBF
	sinkIBF := difference_digest.NewIBF(sourceIBF.Size)
	for k := range sinkSet {
		sinkIBF.Add(k)
	}

	// Sink subtracts source's IBF from it's own
	diff := sinkIBF.Subtract(sourceIBF)

	sinkWithoutSource, _, ok := diff.Decode()
	assert.True(t, ok)
	assert.Len(t, sinkWithoutSource, deleteCount)

	// Sink now knows which elements to delete!
}

func BenchmarkEstimator(b *testing.B) {
	setSize := uint64(1 * 1000 * 1000)
	set := makeSet(0, setSize)
	se := difference_digest.NewStrataEstimator()
	for k := range set {
		se.Add(k)
	}
}

func BenchmarkIBF(b *testing.B) {
	setSize := uint64(1 * 1000 * 1000)
	deleteCount := 100 * 1000

	sourceSet := makeSet(0, setSize)

	alpha := 5
	cells := deleteCount * alpha

	ibf := difference_digest.NewIBF(cells)
	for k := range sourceSet {
		ibf.Add(k)
	}
}

func TestSignificantInts(t *testing.T) {
	testSetWithInt(t, 0, true)
	testSetWithInt(t, 1, true)
	testSetWithInt(t, int64(9223372036854775807-1), true)
	testSetWithInt(t, uint64(18446744073709551615-1), false) // 2 ** 64 - 1 passes
	// testSetWithInt(t, uint64(18446744073709551615), false) // 2 ** 64 fails
	// testSetWithInt(t, -1, true) // Negative numbers are not supported
}

func testSetWithInt(t *testing.T, i interface{}, postgresSupproted bool) {
	t.Run(fmt.Sprintf("Test Int Set with %d", i), func(t *testing.T) {
		db := connectDB()
		defer db.Close()

		_, err := db.Exec("CREATE TEMP TABLE intset (id bigint)")
		assert.NoError(t, err)
		_, err = db.Exec(fmt.Sprintf("INSERT INTO intset (id) VALUES ( %d )", i))
		if postgresSupproted {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
		defer db.Exec("DROP TABLE intset")

		dbIBF, err := difference_digest.EncodeIBFDB(4, db, "intset", "id")
		assert.NoError(t, err)

		set := make(map[uint64]bool)
		switch i.(type) {
		case int:
			set[uint64(i.(int))] = true
		case int64:
			set[uint64(i.(int64))] = true
		case uint64:
			set[i.(uint64)] = true
		default:
			fmt.Printf("%T\n", i)
			panic("AHH")
		}
		setIBF := difference_digest.NewIBF(4)
		for k := range set {
			setIBF.Add(k)
		}

		diff := setIBF.Subtract(dbIBF)

		setWithoutDb, dbWithoutSet, ok := diff.Decode()
		assert.True(t, ok)
		if postgresSupproted {
			assert.Empty(t, setWithoutDb)
		} else {
			assert.Len(t, setWithoutDb, 1)
		}
		assert.Empty(t, dbWithoutSet)
	})
}

func connectDB() *sql.DB {
	// TODO: use docker-compose
	db, err := dburl.Open("postgres://localhost/jason?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = difference_digest.PostgresSetup(db)
	if err != nil {
		panic(err)
	}

	return db
}

func makeSet(min, max uint64) map[uint64]bool {
	a := make(map[uint64]bool)
	i := min
	for {
		if i > max {
			break
		}

		a[uint64(i)] = true
		i++
	}
	return a
}
