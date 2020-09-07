package difference_digest

var (
	hashSeeds         = []uint64{18269989962351869307, 9143901319630896501, 2072764263930962169, 417226483919003799, 16485935163296413021}
	multiplier uint64 = 0x5555555555555555
)

func indiciesHashes(element uint64) []uint64 {
	indices := make([]uint64, 3)
	for i := 0; i < 3; i++ {
		indices[i] = hash(i, element)
	}
	return indices
}

func checkSumHash(element uint64) uint64 {
	return hash(3+0, element)
}

func estimatorHash(element uint64) int {
	h := hash(3+1, element)
	count := 0

	if h == 0 {
		return count
	}

	for (h & 1) == 0 {
		h = h >> 1
		count++
	}

	return count
}

func hash(index int, element uint64) uint64 {
	return multiplier * (element ^ (element >> 32)) * hashSeeds[index]
}
