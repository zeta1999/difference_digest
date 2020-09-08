package difference_digest

var (
	hashSeeds         = []uint64{1826998997, 914390139, 207279169, 4179003799, 1648963021}
	multiplier uint64 = 0x555555555555
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
	return (multiplier * (element ^ (element >> 32)) * hashSeeds[index]) % 4294967296
}
