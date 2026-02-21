package bloom

import (
	"testing"
	"time"
)

// Benchmark the Add operation
func BenchmarkBloomFilterAdd(b *testing.B) {
	// Initialize a filter expecting 1 million items with a 1% false positive rate
	bf := New(1_000_000, 0.01)
	item := []byte("golang-performance-testing")

	b.ResetTimer() // Reset timer so setup isn't counted

	for i := 0; i < b.N; i++ {
		bf.Add(item)
	}
}

// Benchmark the Contains operation (testing an existing item)
func BenchmarkBloomFilterContains(b *testing.B) {
	bf := New(1_000_000, 0.01)
	item := []byte("golang-performance-testing")
	bf.Add(item)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bf.Contains(item)
	}
}

// Simulate an expensive operation, like checking a database on disk.
func expensiveDatabaseLookup(key string) bool {
	time.Sleep(50 * time.Microsecond) // Simulate 50µs latency
	return false                      // Simulate a "cache miss" / not found
}

// Benchmark: Looking up a missing item WITHOUT a Bloom filter
func BenchmarkWithoutBloomFilter(b *testing.B) {
	key := "missing-user-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Every single check hits the slow database directly
		_ = expensiveDatabaseLookup(key)
	}
}

// Benchmark: Looking up a missing item WITH a Bloom filter
func BenchmarkWithBloomFilter(b *testing.B) {
	bf := New(1_000_000, 0.01)
	key := []byte("missing-user-id")

	// We intentionally do NOT add the key to the filter,
	// simulating a true negative.

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Fast path: Bloom filter checks first.
		// Since it returns false, we completely skip the database!
		if bf.Contains(key) {
			_ = expensiveDatabaseLookup(string(key))
		}
	}
}
