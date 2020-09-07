# difference_digest

[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hundredwatt/difference_digest)

difference_digest is a [Go](https://golang.org) package for comparing two sets of data on separate hosts or in separate databases with minimal data transfer. This implementation is based on ["What’s the Difference?
Efficient Set Reconciliation without Prior Context"](https://www.ics.uci.edu/~eppstein/pubs/EppGooUye-SIGCOMM-11.pdf).


Currently only driver support for the PostgreSQL database has been implemented

## Installation
```bash
$ go get github.com/hundredwatt/difference_digest
```

```go
import "github.com/hundredwatt/difference_digest"
```

## Example

```go
package difference_digest_test

import (
  "fmt"
  "math"
  "math/rand"
  "time"

  _ "github.com/lib/pq"

  "github.com/hundredwatt/difference_digest"
)

var (
  db *db.SQL
)

func Example_dbs() {
  db.Exec("CREATE TEMP TABLE sourcethings (id bigint)")
  db.Exec("CREATE TEMP TABLE sinkthings (id bigint)")
  db.Exec("INSERT INTO sourcethings (id) SELECT * from generate_series(1,9999900)") // 100 items missing out of 10,000,000
  db.Exec("INSERT INTO sinkthings (id) SELECT * from generate_series(1,10000000)")

  // 1. Use the Strata Estimator to get the approximate number of differences between the two tables
  sourceEstimator, _ := difference_digest.EncodeEstimatorDB(db, "sourcethings", "id")
  sinkEstimator, _ := difference_digest.EncodeEstimatorDB(db, "sinkthings", "id")
  estimatedDeletes := sinkEstimator.EstimateDifference(sourceEstimator)

  // 2. Compute an Invertible Bloom Filter (of the appropriate size) from each source
  alpha := 3.0
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
```

In this example (if the two tables had been on separate, remote databases), downloading both sets of ~10,000,000 rows to a single host to do a row-by-row comparison would have required transferring 153MB of data.

Using difference digests, the amount of data transferred is proportional to the number of differences. In this case, with 100 differences, this comparison would require transferring only 254KB of data.

In this example, difference digests reduced data transfer by 99.9%!

Difference digests are optimal for reconciling data sets synced between different hosts when you expect only a small number of rows to be missing.

## Documentation

GoDocs [https://godoc.org/github.com/hundredwatt/difference_digest](https://godoc.org/github.com/hundredwatt/difference_digest)

## Testing

```
docker-compose up
go test -v
```
