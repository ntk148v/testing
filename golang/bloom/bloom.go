package bloom

import (
	"hash/fnv"
	"math"
	"sync"
)

// BloomFilter represents the probabilistic data structure
type BloomFilter struct {
	mu     sync.RWMutex // Read-write mutex for concurrent access
	bitset []uint64     // Using uint64 chunks for memory density
	m      uint64       // Total number of bits
	k      uint64       // Number of hash functions
}

// New creates a BloomFilter optimized for 'n' expected items and 'p' false positive rate
func New(n uint64, p float64) *BloomFilter {
	// Calculate optimal m and k
	mFloat := -(float64(n) * math.Log(p)) / (math.Pow(math.Log(2), 2))
	kFloat := (mFloat / float64(n)) * math.Log(2)

	m := uint64(math.Ceil(mFloat))
	k := uint64(math.Ceil(kFloat))

	// Ensure k is at least 1
	if k < 1 {
		k = 1
	}

	return &BloomFilter{
		bitset: make([]uint64, (m+63)/64), // ceiling division by 64
		m:      m,
		k:      k,
	}
}

// getHashes calculates the base h1 and h2 using FNV-1a 64-bit hash
func (bf *BloomFilter) getHashes(data []byte) (uint32, uint32) {
	h := fnv.New64a()
	h.Write(data)
	hash64 := h.Sum64()

	h1 := uint32(hash64 & ((1 << 32) - 1)) // Lower 32 bits
	h2 := uint32(hash64 >> 32)             // Upper 32 bits

	return h1, h2
}

// Add inserts data into the BloomFilter
func (bf *BloomFilter) Add(data []byte) {
	h1, h2 := bf.getHashes(data)

	// Lock for writing: Only one goroutine can add at a time
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := uint64(0); i < bf.k; i++ {
		// Kirsch-Mitzenmatcher optimization
		idx := (uint64(h1) + i*uint64(h2)) % bf.m

		// Set the bit at idx
		bf.bitset[idx/64] |= 1 << (idx % 64)
	}
}

// Contains checks if data might be in the Bloom filter.
func (bf *BloomFilter) Contains(data []byte) bool {
	h1, h2 := bf.getHashes(data)

	// Lock for reading: Multiple goroutines can read concurrently
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	for i := uint64(0); i < bf.k; i++ {
		idx := (uint64(h1) + i*uint64(h2)) % bf.m

		// Check if the bit at idx is 0
		if (bf.bitset[idx/64] & (1 << (idx % 64))) == 0 {
			return false // Definitely not in the set
		}
	}

	return true // Probably in the set
}
