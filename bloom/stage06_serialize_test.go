// Tests for Stage 6: serialize, deserialize, and saturation estimation.
//
// Read the stage on karnstack.com/build/bloom-filter/06-serialize-and-saturation
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section.

package bloom_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage06_RoundTripPreservesMembership verifies marshal followed by
// unmarshal yields a filter that returns the same Test answer for every
// key. Catches a serialization format that drops bits, mis-orders bytes,
// or fails to round-trip m and k.
func TestStage06_RoundTripPreservesMembership(t *testing.T) {
	src := bloom.NewWithK(95850, 7)
	const n = 5000
	for i := 0; i < n; i++ {
		src.Add([]byte(fmt.Sprintf("key-%d", i)))
	}

	data, err := src.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("MarshalBinary returned empty bytes")
	}

	dst := bloom.NewWithK(1, 1) // dummy; UnmarshalBinary rewrites state
	if err := dst.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if dst.M() != src.M() {
		t.Fatalf("M() mismatch after round-trip: dst=%d src=%d", dst.M(), src.M())
	}
	if dst.K() != src.K() {
		t.Fatalf("K() mismatch after round-trip: dst=%d src=%d", dst.K(), src.K())
	}

	for i := 0; i < n; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		if src.Test(key) != dst.Test(key) {
			t.Errorf("Test(%q) differs after round-trip", key)
		}
	}
}

// TestStage06_SerializationIsDeterministic verifies marshaling the same
// state twice produces byte-equal output. Catches a header that captures
// process-time, map iteration order, or other non-determinism.
func TestStage06_SerializationIsDeterministic(t *testing.T) {
	f := bloom.NewWithK(1024, 4)
	for _, k := range [][]byte{[]byte("a"), []byte("b"), []byte("c")} {
		f.Add(k)
	}
	a, err := f.MarshalBinary()
	if err != nil {
		t.Fatalf("first MarshalBinary: %v", err)
	}
	b, err := f.MarshalBinary()
	if err != nil {
		t.Fatalf("second MarshalBinary: %v", err)
	}
	if len(a) != len(b) {
		t.Fatalf("len mismatch: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("byte %d differs: %02x vs %02x", i, a[i], b[i])
		}
	}
}

// TestStage06_SaturationTracksTheoretical verifies the saturation
// estimate is within 5% of the theoretical 1 - e^(-kn/m). Catches a
// popcount-on-bytes implementation that mishandles uint64 word boundaries.
func TestStage06_SaturationTracksTheoretical(t *testing.T) {
	const (
		m uint64 = 95850
		k uint64 = 7
		n        = 10000
	)
	f := bloom.NewWithK(m, k)
	for i := 0; i < n; i++ {
		f.Add([]byte(fmt.Sprintf("key-%d", i)))
	}
	observed := f.Saturation()
	theoretical := 1 - math.Exp(-float64(k*n)/float64(m))
	delta := math.Abs(observed-theoretical) / theoretical
	if delta > 0.05 {
		t.Fatalf("saturation = %.4f, theoretical %.4f (delta %.2f%%, want <= 5%%)",
			observed, theoretical, delta*100)
	}
}
