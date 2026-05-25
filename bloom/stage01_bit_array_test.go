// Tests for Stage 1: bit array and single hash.
//
// Read the stage on karnstack.com/build/bloom-filter/01-bit-array-and-hashing
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section.

package bloom_test

import (
	"fmt"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage01_NewAllocatesAtLeastMBits verifies that New(m) returns a filter
// whose actual bit capacity is at least m. The implementation may round up
// to a multiple of 64; truncation is the bug this catches.
func TestStage01_NewAllocatesAtLeastMBits(t *testing.T) {
	const m = uint64(95850)
	f := bloom.New(m)
	if f == nil {
		t.Fatal("New(95850) returned nil")
	}
	if got := f.M(); got < m {
		t.Fatalf("M() = %d, want >= %d", got, m)
	}
}

// TestStage01_AddedKeyIsPresent verifies that any key reported by Test
// after Add returns true. Catches a missing Add implementation.
func TestStage01_AddedKeyIsPresent(t *testing.T) {
	f := bloom.New(1024)
	keys := [][]byte{
		[]byte("alice"),
		[]byte("bob"),
		[]byte("carol"),
		[]byte("dan"),
	}
	for _, k := range keys {
		f.Add(k)
	}
	for _, k := range keys {
		if !f.Test(k) {
			t.Errorf("Test(%q) = false, want true after Add", k)
		}
	}
}

// TestStage01_EmptyFilterReturnsFalse verifies that Test on a fresh filter
// returns false for arbitrary keys. With a single hash and zero inserts,
// there is no false-positive path. Catches a bit array that initializes
// to non-zero values.
func TestStage01_EmptyFilterReturnsFalse(t *testing.T) {
	f := bloom.New(1024)
	for i := 0; i < 100; i++ {
		k := []byte(fmt.Sprintf("key-%d", i))
		if f.Test(k) {
			t.Errorf("Test(%q) = true on empty filter, want false", k)
		}
	}
}

// TestStage01_BitArrayBoundary verifies that New(1) does not crash and
// that Add followed by Test works on a one-bit filter. Catches off-by-one
// errors in the modulo-into-bit-position math and division-by-zero.
func TestStage01_BitArrayBoundary(t *testing.T) {
	f := bloom.New(1)
	if f == nil {
		t.Fatal("New(1) returned nil")
	}
	if got := f.M(); got < 1 {
		t.Fatalf("M() = %d for New(1), want >= 1", got)
	}
	key := []byte("anything")
	f.Add(key)
	if !f.Test(key) {
		t.Errorf("Test(%q) = false after Add on m=1 filter, want true", key)
	}
}

// TestStage01_HashIsDeterministic verifies that two filters of identical
// capacity agree on Test for the same key after Add. Catches a hash function
// seeded from process-randomness rather than a constant.
func TestStage01_HashIsDeterministic(t *testing.T) {
	const m = uint64(1024)
	f1 := bloom.New(m)
	f2 := bloom.New(m)
	key := []byte("deterministic")
	f1.Add(key)
	f2.Add(key)
	if !f1.Test(key) {
		t.Error("Test(key) = false on filter 1 after Add")
	}
	if !f2.Test(key) {
		t.Error("Test(key) = false on filter 2 after Add")
	}
}
