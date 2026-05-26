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
	// k is the number of hash functions per probe. Stage 1 sets k=1
	// via New. Stage 2 uses NewWithK for arbitrary k via the
	// Kirsch-Mitzenmacher construction.
	k uint64
	// blocked toggles the stage-4 cache-line-blocked layout. When true,
	// Add/Test pick a 512-bit (8 x uint64) block once via h_a and place
	// all k probes inside that block.
	blocked bool
}

// New returns a filter with at least m bits of capacity, rounded up to the
// nearest multiple of 64. k defaults to 1 (a single hash function).
//
// Stage 1: implement this so the bit array can hold m bits.
func New(m uint64) *Filter {
	// TODO(stage1): allocate the bit array, set f.m to the rounded-up size.
	return nil
}

// NewWithK returns a filter sized like New but using k hash functions per
// Add/Test via the Kirsch-Mitzenmacher two-hash construction
// h_i(key) = h_a(key) + i * h_b(key) for i in [0, k).
//
// Stage 2: implement this.
func NewWithK(m, k uint64) *Filter {
	// TODO(stage2): build on New + remember k for Add/Test.
	return nil
}

// Add inserts the key into the filter.
//
// Stage 1: hash key, map to a bit position in [0, m), set the bit.
// Stage 2: extend to k bit positions via Kirsch-Mitzenmacher.
func (f *Filter) Add(key []byte) {
	// TODO(stage1/stage2): set k bit positions driven by k hashes of key.
}

// Test reports whether the key may be in the filter. Returns false only
// if the key was definitely not added.
//
// Stage 1: hash key, map to a bit position in [0, m), check the bit.
// Stage 2: extend to k bit positions, ALL must be set.
func (f *Filter) Test(key []byte) bool {
	// TODO(stage1/stage2): check k bit positions; all must be set.
	return false
}

// M returns the actual bit capacity (after rounding up).
func (f *Filter) M() uint64 {
	return f.m
}

// K returns the number of hash functions per probe.
func (f *Filter) K() uint64 {
	return f.k
}

// OptimalSize returns the (m, k) that minimize the false-positive rate
// for n keys at target rate p, using the closed-form bounds:
//
//	m = ceil(-n * ln(p) / (ln 2)^2)
//	k = round((m / n) * ln 2)
//
// Stage 3: implement this.
func OptimalSize(n uint64, p float64) (m uint64, k int) {
	// TODO(stage3).
	return 0, 0
}

// NewBlocked returns a filter sized like NewWithK but using the
// cache-line-blocked layout from Putze-Sanders-Singler 2007: a primary
// hash picks one 512-bit (eight uint64) block, then all k probes live
// inside that block. Same Add/Test signatures as Filter.
//
// Stage 4: implement this.
func NewBlocked(m, k uint64) *Filter {
	// TODO(stage4): build a blocked-layout filter.
	return nil
}

// MarshalBinary serializes the filter's m, k, and bit array into a
// byte slice. Use a fixed little-endian header so two implementations of
// the same stage produce byte-equal output for the same state.
//
// Suggested format:
//
//	bytes 0..3:  magic "BLM1"
//	bytes 4..11: m (uint64 LE)
//	bytes 12..19: k (uint64 LE)
//	byte  20:    flags (bit 0 = blocked layout)
//	byte  21..: bit array, len = m/8 bytes, little-endian within each word
//
// Stage 6: implement this.
func (f *Filter) MarshalBinary() ([]byte, error) {
	// TODO(stage6).
	return nil, nil
}

// UnmarshalBinary reverses MarshalBinary. After this call, f.M(), f.K()
// and Test results match the source filter.
//
// Stage 6: implement this.
func (f *Filter) UnmarshalBinary(data []byte) error {
	// TODO(stage6).
	return nil
}

// Saturation returns the fraction of bits currently set in the bit
// array. Useful for monitoring fill in production. The theoretical
// expectation is 1 - exp(-k*n/m).
//
// Stage 6: implement this.
func (f *Filter) Saturation() float64 {
	// TODO(stage6).
	return 0
}
