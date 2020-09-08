package difference_digest

import (
	"fmt"
	"testing"

	_ "github.com/snowflakedb/gosnowflake"
)

func xTestHashSpread(t *testing.T) {
	n := uint64(10)
	bins := make(map[int]int)
	for i := 0; i < 4294967296; i++ {
		bins[int(hash(0, uint64(i))%n)]++
	}

	fmt.Println("Modulus bins: (should be uniformly distributed)")
	fmt.Println(bins)
}

func xTestStrataEstimatorSpread(t *testing.T) {
	bins := make(map[int]int)
	for i := 0; i < 4294967296; i++ {
		bins[estimatorHash(uint64(i))]++
	}

	fmt.Println("Estimator bins: (should be logarithmically distributed)")
	fmt.Println(bins)
}
