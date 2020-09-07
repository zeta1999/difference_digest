package difference_digest

import (
	"encoding/binary"
	"encoding/hex"
)

type databaseDriver int

const (
	// PostgreSQL represents the PostgreSQL database
	PostgreSQL databaseDriver = iota
)

// DatabaseType indicates which type of database to use
var DatabaseType databaseDriver = PostgreSQL

func query(name string) string {
	var queries map[string]string
	switch DatabaseType {
	case PostgreSQL:
		queries = postgresQueries
	}

	return queries[name]
}

// Todo: Private

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

func (bits *Bitmap) IsZero() bool {
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
