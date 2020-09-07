package difference_digest

import (
	"database/sql"
	"fmt"
)

// InvertibleBloomFilter is a data structure for compactly storing a recoverable representation of a set
// See: https://www.ics.uci.edu/~eppstein/pubs/EppGooUye-SIGCOMM-11.pdf
type InvertibleBloomFilter struct {
	Cells []IBFCell
	Size  int
}

// IBFCell represents one cell of an Invertible Bloom Filter
type IBFCell struct {
	IDSum   Bitmap
	HashSum Bitmap
	Count   int64
}

// Insert inserts an element into the IBF Cell
func (b *IBFCell) Insert(s uint64) {
	b.IDSum.XOR(ToBitmap(s))
	b.HashSum.XOR(ToBitmap(checkSumHash(s)))
	b.Count++
}

// Subtract the value of another IBF Cell
func (b *IBFCell) Subtract(b2 *IBFCell) {
	b.IDSum.XOR(&b2.IDSum)
	b.HashSum.XOR(&b2.HashSum)
	b.Count -= b2.Count
}

// IsPure returns true when the IBFCell has a Count of 1 or -1 and the HashSum is identical to the IDSum
func (b *IBFCell) IsPure() bool {
	return (b.Count == 1 || b.Count == -1) && b.HashSum.Uint64() == checkSumHash(b.IDSum.Uint64())
}

// IsZero returns true when the IBFCell is empty (all values equal to 0)
func (b *IBFCell) IsZero() bool {
	return b.IDSum.IsZero() && b.HashSum.IsZero() && b.Count == 0
}

// NewIBF initalizes an empty InvertibleBloomFilter
func NewIBF(size int) *InvertibleBloomFilter {
	return &InvertibleBloomFilter{
		Cells: make([]IBFCell, size),
		Size:  size,
	}
}

// EncodeIBFDB encodes an InvertibleBloomFilter from a column in a database table;
// currently, only PostgreSQL is supported and column mast have type BIGINT/INT8
func EncodeIBFDB(size int, db *sql.DB, table string, column string) (*InvertibleBloomFilter, error) {
	rows, err := db.Query(fmt.Sprintf(query("ibf"), table, column, size))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	b := NewIBF(size)

	for rows.Next() {
		var (
			cell           int
			IDSum, HashSum uint64
			Count          int64
		)

		err := rows.Scan(&cell, &IDSum, &HashSum, &Count)
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
		b.Cells[cell] = el
	}

	return b, nil
}

// Add inserts an element into the InvertibleBloomFilter
func (ibf *InvertibleBloomFilter) Add(s uint64) {
	for _, h := range indiciesHashes(s) {
		j := h % uint64(ibf.Size)
		ibf.Cells[j].Insert(s)
	}
}

// Subtract computes the difference between 2 Invertible Bloom Filter
func (ibf *InvertibleBloomFilter) Subtract(ibf2 *InvertibleBloomFilter) *InvertibleBloomFilter {
	difference := NewIBF(ibf.Size)
	copy(difference.Cells, ibf.Cells)

	for j := 0; j < ibf.Size; j++ {
		difference.Cells[j].Subtract(&ibf2.Cells[j])
	}

	return difference
}

// Decode recovers the Cells represented in the Invertible Bloom Filter (use only after performing a Subtract())
func (ibf *InvertibleBloomFilter) Decode() (aWithoutB []uint64, bWithoutA []uint64, ok bool) {
	pureList := make([]int, 0)

	for {
		n := len(pureList) - 1

		if n == -1 {
			for j := 0; j < ibf.Size; j++ {
				if ibf.Cells[j].IsPure() {
					pureList = append(pureList, j)
				}
			}
			if len(pureList) == 0 {
				break
			}
			continue
		}

		j := pureList[n]
		pureList = pureList[:n]

		if !ibf.Cells[j].IsPure() {
			continue
		}

		s := ibf.Cells[j].IDSum
		c := ibf.Cells[j].Count

		if c > 0 {
			aWithoutB = append(aWithoutB, s.Uint64())
		} else {
			bWithoutA = append(bWithoutA, s.Uint64())
		}
		for _, h := range indiciesHashes(s.Uint64()) {
			j2 := h % uint64(ibf.Size)
			ibf.Cells[j2].IDSum.XOR(&s)
			ibf.Cells[j2].HashSum.XOR(ToBitmap(checkSumHash(s.Uint64())))
			ibf.Cells[j2].Count -= c
		}
	}
	for j := 0; j < ibf.Size; j++ {
		if !ibf.Cells[j].IsZero() {
			ok = false
			return
		}
	}

	ok = true
	return
}
