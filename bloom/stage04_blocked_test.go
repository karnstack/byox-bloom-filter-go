// Tests for Stage 4: cache-line-blocked layout.
//
// Read the stage on karnstack.com/build/bloom-filter/04-cache-line-blocked
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section.

package bloom_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage04_NewBlockedAllocates verifies the constructor exists, takes
// (m, k), and reports rounded-up capacity at least m. Catches a missing
// NewBlocked or a flat-layout fallback that ignores k.
func TestStage04_NewBlockedAllocates(t *testing.T) {
	f := bloom.NewBlocked(95850, 7)
	if f == nil {
		t.Fatal("NewBlocked(95850, 7) returned nil")
	}
	if got := f.M(); got < 95850 {
		t.Fatalf("M() = %d, want >= 95850", got)
	}
	if got := f.K(); got != 7 {
		t.Fatalf("K() = %d, want 7", got)
	}
}

// TestStage04_AddedKeysArePresent verifies the basic Add/Test round-trip
// on the blocked layout. All added keys must Test true. Catches off-by-
// one within a block, or probes accidentally crossing block boundaries
// during Test.
func TestStage04_AddedKeysArePresent(t *testing.T) {
	f := bloom.NewBlocked(95850, 7)
	keys := [][]byte{
		[]byte("alice"), []byte("bob"), []byte("carol"),
		[]byte("dan"), []byte("eve"), []byte("frank"),
	}
	for _, k := range keys {
		f.Add(k)
	}
	for _, k := range keys {
		if !f.Test(k) {
			t.Errorf("Test(%q) = false, want true after Add on blocked filter", k)
		}
	}
}

// TestStage04_FPRateBelowTheoreticalBound is the moat. The blocked layout
// trades a small FP-rate cost for cache locality (Putze-Sanders-Singler
// 2007); the bound is 2x the standard theoretical (vs 2x in stage 2 for
// the same m and k). A broken block-selection scheme blows past that.
func TestStage04_FPRateBelowTheoreticalBound(t *testing.T) {
	const (
		m      uint64 = 95850
		k      uint64 = 7
		n             = 10000
		trials        = 100000
	)
	f := bloom.NewBlocked(m, k)
	for i := 0; i < n; i++ {
		f.Add([]byte(fmt.Sprintf("added-%d", i)))
	}
	fps := 0
	for i := 0; i < trials; i++ {
		if f.Test([]byte(fmt.Sprintf("query-%d", i))) {
			fps++
		}
	}
	observed := float64(fps) / float64(trials)
	theoretical := math.Pow(
		1-math.Exp(-float64(k*n)/float64(m)),
		float64(k),
	)
	bound := 2 * theoretical
	if observed > bound {
		t.Fatalf(
			"blocked FP rate %.4f exceeds 2x theoretical %.4f (bound %.4f)",
			observed, theoretical, bound,
		)
	}
}
