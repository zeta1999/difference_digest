package difference_digest

import (
	"database/sql"
	"fmt"
	"math"
)

const (
	stratumCount = 64
	cellsCount   = 80
)

// StrataEstimator is a data structure used to estimate the number of differences between 2 sets probablistically
type StrataEstimator struct {
	Stratum []InvertibleBloomFilter
}

// NewStrataEstimator initalizes a new StrataEstimator
func NewStrataEstimator() *StrataEstimator {
	se := StrataEstimator{
		Stratum: make([]InvertibleBloomFilter, stratumCount),
	}

	for i := range se.Stratum {
		se.Stratum[i] = *NewIBF(cellsCount)
	}

	return &se
}

// Add adds an element to the StrataEstimator
func (se *StrataEstimator) Add(element uint64) {
	j := estimatorHash(element)
	se.Stratum[j].Add(element)
}

// EstimateDifference returns the estimated number of differences between the receiver and a 2nd Strata Estimator
func (se *StrataEstimator) EstimateDifference(se2 *StrataEstimator) uint64 {
	var Count uint64 = 0

	for i := 63; i >= 0; i-- {
		diff := se.Stratum[i].Subtract(&se2.Stratum[i])
		aWb, _, ok := diff.Decode()
		if ok {
			Count += uint64(len(aWb))
		} else {
			return uint64(math.Pow(2.0, float64(i+1))) * (Count + 1)
		}
	}

	return Count
}

// EncodeEstimatorDB queries a PostgreSQL database and returns a StrataEstimator for the specified table and column
func EncodeEstimatorDB(db *sql.DB, table string, column string) (*StrataEstimator, error) {
	rows, err := db.Query(fmt.Sprintf(query("strata_estimator"), table, column, cellsCount))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	estimator := NewStrataEstimator()

	for rows.Next() {
		var (
			strata, cell   int
			IDSum, HashSum uint64
			Count          int64
		)

		err := rows.Scan(&strata, &cell, &IDSum, &HashSum, &Count)
		if err != nil {
			return nil, err
		}

		idBitmap := ToBitmap(IDSum)
		hashBitmap := ToBitmap(HashSum)

		el := IBFCell{
			IDSum:   *idBitmap,
			HashSum: *hashBitmap,
			Count:   Count,
		}
		estimator.Stratum[strata].Cells[cell] = el
	}

	return estimator, nil
}
