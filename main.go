package main

import (
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

// Based on https://www.ics.uci.edu/~eppstein/pubs/EppGooUye-SIGCOMM-11.pdf

type Bitmap [8]byte

func (bits *Bitmap) IsSet(i int) bool { i -= 1; return bits[i/8]&(1<<uint(7-i%8)) != 0 }
func (bits *Bitmap) Set(i int)        { i -= 1; bits[i/8] |= 1 << uint(7-i%8) }
func (bits *Bitmap) Clear(i int)      { i -= 1; bits[i/8] &^= 1 << uint(7-i%8) }

func (bits *Bitmap) Sets(xs ...int) {
	for _, x := range xs {
		bits.Set(x)
	}
}

func (bits *Bitmap) XOR(other *Bitmap) {
	for i := 1; i <= 64; i++ {
		if bits.IsSet(i) != other.IsSet(i) {
			bits.Set(i)
		} else {
			bits.Clear(i)
		}
	}
}

func (bits *Bitmap) Zero() bool {
	for i := 1; i <= 64; i++ {
		if bits.IsSet(i) {
			return false
		}
	}
	return true
}

func ToBitmap(n uint64) *Bitmap {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, n)

	bits := Bitmap{}

	for i := 0; i < 64; i++ {
		if bs[i/8]&(1<<uint(7-i%8)) != 0 {
			bits.Set(i + 1)
		}
	}
	return &bits
}

func (bits Bitmap) String() string { return hex.EncodeToString(bits[:]) }
func (bits Bitmap) Uint64() uint64 { return binary.LittleEndian.Uint64(bits[:]) }

var (
	min = 0
	max = 100

	hashSeeds        = []uint64{18269989962351869307, 9143901319630896501, 2072764263930962169, 417226483919003799, 16485935163296413021}
	mask      uint64 = 0x5555555555555555
	// prime1    uint64 = 2147483647
	// prime2    uint64 = 2654435761
	// prime3    uint64 = 4638325457
	// offset    uint64 = 4869563535558310077
)

func hash(index int, element uint64) uint64 {
	// bs := make([]byte, 8)
	// binary.LittleEndian.PutUint64(bs, element)
	// hashfn := murmur3.New64WithSeed(uint32(hashSeeds[index] >> 32))
	// hashfn.Write(bs)
	// return hashfn.Sum64()

	// h1 := ((hashSeeds[0] ^ (element)) + offset) * (prime1 * prime2)
	// h2 := ((hashSeeds[1] ^ (element)) + offset) * (prime1 * prime3)
	// return h1 + uint64(index+1)*h2

	// inner := mask * (element ^ (element >> 32))
	// outer := hashSeeds[index] * (inner ^ (inner >> 32))
	// return outer

	return mask * (element ^ (element >> 32)) * hashSeeds[index]

	// bs := make([]byte, 9)
	// binary.LittleEndian.PutUint64(bs, element)
	// bs[8] = byte(index)
	// h := md5.Sum(bs)
	// return binary.LittleEndian.Uint64((h[:8]))
}

func hashes(element uint64) []uint64 {
	indices := make([]uint64, 3)
	for i := 0; i < 3; i++ {
		indices[i] = hash(i, element)
	}
	return indices
}

func hashsum(element uint64) uint64 {
	return hash(3+0, element)
}

func hashestimator(element uint64) int {
	h := hash(3+1, element)
	count := 0

	if h == 0 {
		return count
	}

	for {
		if (h & 1) == 1 {
			break
		}

		h = h >> 1
		count++
	}

	return count
}

type ibf []ibfelement

type ibfelement struct {
	idSum   Bitmap
	hashSum Bitmap
	count   int64
}

func (bj *ibfelement) Pure() bool {
	return (bj.count == 1 || bj.count == -1) && bj.hashSum.Uint64() == hashsum(bj.idSum.Uint64())
}

func EncodeIBF(size int, set map[uint64]bool) *ibf {
	b := make(ibf, size)
	for s := range set {
		b.Add(s)
	}

	return &b
}

func EncodeIBFDB(size int, db *sql.DB, table string, column string) (*ibf, error) {
	query := `
	SELECT 
		f_hash(idx, %[2]s) %% %[3]d AS cell, 
		bit_xor(%[2]s::bigint) AS id_sum, 
		bit_xor_numeric(f_hash(3 + 0, %[2]s)) AS hash_sum,  
		COUNT(id) AS count
	FROM (
		SELECT 0 AS idx, * FROM %[1]s UNION 
		SELECT 1, * FROM %[1]s UNION
		SELECT 2, * FROM %[1]s 
	) s
	GROUP BY 1
	ORDER BY 1;
`

	rows, err := db.Query(fmt.Sprintf(query, table, column, size))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	b := make(ibf, size)

	for rows.Next() {
		var (
			cell           int
			idSum, hashSum uint64
			count          int64
		)

		err := rows.Scan(&cell, &idSum, &hashSum, &count)
		if err != nil {
			return nil, err
		}

		idBitmap := ToBitmap(idSum)
		hashBitmap := ToBitmap(hashSum)

		el := ibfelement{
			idSum:   *idBitmap,
			hashSum: *hashBitmap,
			count:   count,
		}
		b[cell] = el
	}

	return &b, nil
}

func (b *ibf) Add(s uint64) {
	for _, h := range hashes(s) {
		j := h % uint64(len(*b))
		(*b)[j].idSum.XOR(ToBitmap(s))
		(*b)[j].hashSum.XOR(ToBitmap(hashsum(s)))
		(*b)[j].count++
	}
}

func (f *ibf) Subtract(other *ibf) *ibf {
	result := make(ibf, len(*f))
	copy(result, *f)

	for j := 0; j < len(*f); j++ {
		result[j].idSum.XOR(&(*other)[j].idSum)
		result[j].hashSum.XOR(&(*other)[j].hashSum)
		result[j].count -= (*other)[j].count
	}

	return &result
}

func (f *ibf) Decode() (aWithoutB []uint64, bWithoutA []uint64, ok bool) {
	pureList := make([]int, 0)
	b := *f

	for j := 0; j < len(*f); j++ {
		if b[j].Pure() {
			pureList = append(pureList, j)
		}
	}

	for {
		n := len(pureList) - 1
		if n == -1 {
			for j := 0; j < len(*f); j++ {
				if b[j].Pure() {
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

		if !b[j].Pure() {
			continue
		}

		s := b[j].idSum
		c := b[j].count

		if c > 0 {
			aWithoutB = append(aWithoutB, s.Uint64())
		} else {
			bWithoutA = append(bWithoutA, s.Uint64())
		}
		for _, h := range hashes(s.Uint64()) {
			j2 := h % uint64(len(b))
			b[j2].idSum.XOR(&s)
			b[j2].hashSum.XOR(ToBitmap(hashsum(s.Uint64())))
			b[j2].count -= c
		}
	}
	for j := 0; j < len(*f); j++ {
		if !b[j].idSum.Zero() || !b[j].hashSum.Zero() || b[j].count != 0 {
			ok = false
			return
		}
	}

	ok = true
	return
}

type strataEstimator []ibf

func EncodeEstimator(set map[uint64]bool) *strataEstimator {
	estimator := make(strataEstimator, 64)
	for i := range estimator {
		estimator[i] = make(ibf, 80)
	}

	for s := range set {
		j := hashestimator(s)
		estimator[j].Add(s)
	}

	return &estimator
}

func EncodeEstimatorDB(db *sql.DB, table string, column string) (*strataEstimator, error) {
	query := `
		SELECT 
			f_trailing_zeros(f_hash(3 + 1, %[2]s)) AS estimator, 
			f_hash(idx, %[2]s) %% 80 AS cell, 
			bit_xor(%[2]s::bigint) AS id_sum, 
			bit_xor_numeric(f_hash(3 + 0, %[2]s)) AS hash_sum,  
			COUNT(id) AS count
		FROM (
			SELECT 0 AS idx, * FROM %[1]s UNION 
			SELECT 1, * FROM %[1]s UNION
			SELECT 2, * FROM %[1]s 
		) s
		GROUP BY 1, 2 
		ORDER BY 1, 2;
	`

	rows, err := db.Query(fmt.Sprintf(query, table, column))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	estimator := make(strataEstimator, 64)
	for i := range estimator {
		estimator[i] = make(ibf, 80)
	}

	for rows.Next() {
		var (
			strata, cell   int
			idSum, hashSum uint64
			count          int64
		)

		err := rows.Scan(&strata, &cell, &idSum, &hashSum, &count)
		if err != nil {
			return nil, err
		}

		idBitmap := ToBitmap(idSum)
		hashBitmap := ToBitmap(hashSum)

		el := ibfelement{
			idSum:   *idBitmap,
			hashSum: *hashBitmap,
			count:   count,
		}
		estimator[strata][cell] = el
	}

	return &estimator, nil
}

func (estimator *strataEstimator) Decode(otherEstmator *strataEstimator) uint64 {
	var count uint64 = 0

	for i := 63; i >= 0; i-- {
		diff := (*estimator)[i].Subtract(&(*otherEstmator)[i])
		aWb, _, ok := diff.Decode()
		if ok {
			count += uint64(len(aWb))
		} else {
			return uint64(math.Pow(2.0, float64(i+1))) * (count + 1)
		}
	}

	return count
}
