// Package bloom contains your Bloom filter implementation.
//
// The interface below is the contract karnstack tests target. Do not rename
// the exported methods. Implementation lives in this file across six stages:
//
//	Stage 1: bit array storage and a single hash function.
//	Stage 2: multiple hash functions via the Kirsch-Mitzenmacher construction.
//	Stage 3: optimal sizing helpers (m and k from a target false-positive rate).
//	Stage 4: cache-line-blocked layout.
//	Stage 5: concurrent-safe Add via atomic OR.
//	Stage 6: serialize, deserialize, and saturation estimation.
//
// Read the stage on karnstack.com/build/bloom-filter before implementing.
package bloom

// Filter is a Bloom filter.
type Filter struct {
	// bits stores the bit array packed into 64-bit words.
	bits []uint64
	// m is the actual bit capacity (multiple of 64).
	m uint64
}

// New returns a filter with at least m bits of capacity, rounded up to the
// nearest multiple of 64.
//
// Stage 1: implement this so the bit array can hold m bits.
func New(m uint64) *Filter {
	// TODO(stage1): allocate the bit array, set f.m to the rounded-up size.
	return nil
}

// Add inserts the key into the filter.
//
// Stage 1: hash key, map to a bit position in [0, m), set the bit.
func (f *Filter) Add(key []byte) {
	// TODO(stage1): set one bit driven by hash(key) % m.
}

// Test reports whether the key may be in the filter. Returns false only
// if the key was definitely not added.
//
// Stage 1: hash key, map to a bit position in [0, m), check the bit.
func (f *Filter) Test(key []byte) bool {
	// TODO(stage1): check one bit driven by hash(key) % m.
	return false
}

// M returns the actual bit capacity (after rounding up).
func (f *Filter) M() uint64 {
	return f.m
}
