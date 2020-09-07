package difference_digest_test

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	_ "github.com/lib/pq"

	"github.com/hundredwatt/difference_digest"
)

func Example() {
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

	// 1. Compute an estimator for sink
	sinkEstimator := difference_digest.NewStrataEstimator()
	for k := range sinkSet {
		sinkEstimator.Add(k)
	}

	// 2. Compute an estimator for source
	sourceEstimator := difference_digest.NewStrataEstimator()
	for k := range sourceSet {
		sourceEstimator.Add(k)
	}

	// 3. Caluclate the estimated number of differences between the two sets
	estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

	// 4. Compute an IBF of the appropriate size for the source set
	alpha := 5.0
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	sourceIBF := difference_digest.NewIBF(cells)
	for k := range sourceSet {
		sourceIBF.Add(k)
	}

	// 5. Compute an IBF of the appropriate size for the sink set
	sinkIBF := difference_digest.NewIBF(sourceIBF.Size)
	for k := range sinkSet {
		sinkIBF.Add(k)
	}

	// 6. Subtract the two IBFS and Decode the result to find the differences
	diff := sinkIBF.Subtract(sourceIBF)

	sinkWithoutSource, _, ok := diff.Decode()
	if !ok {
		fmt.Println("Invertible Bloom Filter failed to decode, please try again")
	} else {
		fmt.Printf("%d elements found in sink that are not in source", len(sinkWithoutSource))
		// Output: 42 elements found in sink that are not in source
	}
}
func Example_dbs() {
	db := connectDB()
	defer db.Close()

	db.Exec("CREATE TEMP TABLE sourcethings (id bigint)")
	db.Exec("CREATE TEMP TABLE sinkthings (id bigint)")
	db.Exec("INSERT INTO sourcethings (id) SELECT * from generate_series(1,9900)") // 100 items missing
	db.Exec("INSERT INTO sinkthings (id) SELECT * from generate_series(1,10000)")

	// 1. Use the Strata Estimator to get the approximate number of differences between the two tables
	sourceEstimator, _ := difference_digest.EncodeEstimatorDB(db, "sourcethings", "id")
	sinkEstimator, _ := difference_digest.EncodeEstimatorDB(db, "sinkthings", "id")
	estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

	// 2. Get an IBF of the appropriate size from each source
	alpha := 5.0
	cells := int(math.Ceil(float64(estimatedDeletes) * float64(alpha)))
	sourceIBF, _ := difference_digest.EncodeIBFDB(cells, db, "sourcethings", "id")
	sinkIBF, _ := difference_digest.EncodeIBFDB(cells, db, "sinkthings", "id")

	// 3. Compute the difference of the IBFs
	diff := sinkIBF.Subtract(sourceIBF)
	sinkWithoutSource, _, ok := diff.Decode()
	if !ok {
		fmt.Println("Invertible Bloom Filter failed to decode, please try again")
	} else {
		fmt.Printf("%d elements found in sink that are not in source", len(sinkWithoutSource))
		// Output: 100 elements found in sink that are not in source
	}
}
