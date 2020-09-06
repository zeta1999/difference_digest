package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"
)

func TestDBComputeDifference(t *testing.T) {
	db, err := dburl.Open("postgres://localhost/jason?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TEMP TABLE sourcethings (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("CREATE TEMP TABLE sinkthings (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO sourcethings (id) SELECT * from generate_series(1,9900)") // 100 items missing
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO sinkthings (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	// 1. Use the Strata Estimator to get the approximate number of differences
	sourceEstimator, err := EncodeEstimatorDB(db, "sourcethings", "id")
	assert.NoError(t, err)
	sinkEstimator, err := EncodeEstimatorDB(db, "sinkthings", "id")
	assert.NoError(t, err)

	estimatedDeletes := sinkEstimator.Decode(sourceEstimator)

	assert.Less(t, int(estimatedDeletes), 150)
	assert.Greater(t, int(estimatedDeletes), 50)

	// 2. Get an IBF of the appropriate size from each source
	alpha := 5
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	sourceIBF, err := EncodeIBFDB(cells, db, "sourcethings", "id")
	assert.NoError(t, err)
	sinkIBF, err := EncodeIBFDB(cells, db, "sinkthings", "id")
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
	db, err := dburl.Open("postgres://localhost/jason?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TEMP TABLE things (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO things (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	dbEstimator, err := EncodeEstimatorDB(db, "things", "id")
	assert.NoError(t, err)

	set := makeSet(1, 10000)
	inMemoryEstimator := EncodeEstimator(set)

	assert.Equal(t, 0, int(inMemoryEstimator.Decode(dbEstimator)))

	for i := range *inMemoryEstimator {
		diff := (*dbEstimator)[i].Subtract(&(*inMemoryEstimator)[i])
		aWithoutB, bWithoutA, ok := diff.Decode()
		assert.True(t, ok)
		assert.Empty(t, aWithoutB)
		assert.Empty(t, bWithoutA)
	}
}

func TestDBInMemoryIBFEquivalence(t *testing.T) {
	db, err := dburl.Open("postgres://localhost/jason?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TEMP TABLE things (id bigint)")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO things (id) SELECT * from generate_series(1,10000)")
	assert.NoError(t, err)

	dbIBF, err := EncodeIBFDB(10, db, "things", "id")
	assert.NoError(t, err)

	set := makeSet(1, 10000)
	setIBF := EncodeIBF(10, set)

	diff := dbIBF.Subtract(setIBF)

	aWithoutB, bWithoutA, ok := diff.Decode()
	assert.True(t, ok)
	assert.Empty(t, aWithoutB)
	assert.Empty(t, bWithoutA)
}

func TestInMemory(t *testing.T) {

	fmt.Println("Generating data...")
	setSize := uint64(1 * 1000 * 1000)
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
	fmt.Println(len(deletes))

	fmt.Println("Generating estimators...")
	// Sink:
	// Sink computes an estimator and sends it to Sink
	PrintMemUsage()
	sinkEstimator := EncodeEstimator(sinkSet)
	PrintMemUsage()

	// Source:
	// Source computes its own estimator, and then decodes it with sink's
	sourceEstimator := EncodeEstimator(sourceSet)
	PrintMemUsage()

	estimatedDeletes := sinkEstimator.Decode(sourceEstimator)
	fmt.Println(estimatedDeletes)

	fmt.Println("Generating IBFs...")
	PrintMemUsage()
	alpha := 5
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	// Source computes its IBF and sends it to sink
	sourceIBF := EncodeIBF(cells, sourceSet)
	PrintMemUsage()

	// Sink:
	// Sink computes its IBF
	sinkIBF := EncodeIBF(len(*sourceIBF), sinkSet)
	PrintMemUsage()

	// Sink subtracts source's IBF from it's own
	diff := sinkIBF.Subtract(sourceIBF)
	PrintMemUsage()

	sinkWithoutSource, _, ok := diff.Decode()
	if !ok {
		// Retry the process again with different hashses
		fmt.Println("FAIL")
		fmt.Println(len(sinkWithoutSource))
	} else {
		fmt.Println(len(sinkWithoutSource))
	}

	// Sink now knows which elements to delete!
}

func xTestHashSpread(t *testing.T) {
	hashResults := make(map[int]uint64)
	for i := 0; i < 10000; i++ {
		hashResults[i] = hash(0, uint64(i))
	}

	n := 10
	bins := make(map[int]int)
	for _, j := range hashResults {
		bins[int(j%uint64(n))]++
	}

	fmt.Println("Modulus bins: (should be uniformly distributed)")
	fmt.Println(bins)

	sbins := make(map[int]int)
	for i := range hashResults {
		s := hashestimator(uint64(i))
		sbins[s]++
	}

	fmt.Println("Estimator bins: (should be logarithmically distributed")
	fmt.Println(sbins)
}

func BenchmarkEstimator(b *testing.B) {
	setSize := uint64(1 * 1000 * 1000)
	set := makeSet(0, setSize)
	EncodeEstimator(set)
}

func BenchmarkIBF(b *testing.B) {
	setSize := uint64(1 * 1000 * 1000)
	deleteCount := 100 * 1000

	sourceSet := makeSet(0, setSize)

	alpha := 5
	cells := deleteCount * alpha

	EncodeIBF(cells, sourceSet)
}

func TestMain(t *testing.T) {
	// for i := 1; i <= 4; i++ {
	// 	for idx := 0; idx < 3; idx++ {
	// 		fmt.Printf("%d | %d | %d | %d\n", i, idx, hashestimator(uint64(i)), hashes(uint64(i))[idx]%80)
	// 	}
	// }
	set := makeSet(1, 10000)
	estimator := EncodeEstimator(set)
	for i, b := range (*estimator)[0] {
		fmt.Printf("%d %d %d %d\n", i, b.idSum.Uint64(), b.hashSum.Uint64(), b.count)
		if i > 33 {
			break
		}
	}
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
