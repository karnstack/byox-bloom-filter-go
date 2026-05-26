// Tests for Stage 2: multiple hashes via the Kirsch-Mitzenmacher
// construction.
//
// Read the stage on karnstack.com/build/bloom-filter/02-multiple-hashes
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section.

package bloom_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage02_NewWithKAllocates verifies the multi-hash constructor
// exists and that the rounded-up capacity invariant still holds. Catches
// a missing NewWithK or one that drops the k parameter on the floor.
func TestStage02_NewWithKAllocates(t *testing.T) {
	f := bloom.NewWithK(1024, 3)
	if f == nil {
		t.Fatal("NewWithK(1024, 3) returned nil")
	}
	if got := f.M(); got < 1024 {
		t.Fatalf("M() = %d, want >= 1024", got)
	}
}

// TestStage02_AddedKeysArePresent verifies Add/Test still work end to end
// once Add is using k hashes. Failure here means the multi-hash
// construction is mis-indexing (off-by-one on i, wrong modulo, etc.).
func TestStage02_AddedKeysArePresent(t *testing.T) {
	f := bloom.NewWithK(95850, 7)
	keys := [][]byte{
		[]byte("alice"),
		[]byte("bob"),
		[]byte("carol"),
		[]byte("dan"),
		[]byte("eve"),
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

// TestStage02_FPRateBelowTheoreticalBound is the moat. Inserts n=10000
// keys into an m=95850 / k=7 filter, queries 100k unseen keys, and
// verifies the observed false-positive rate stays within 2x of the
// theoretical (1 - e^-kn/m)^k bound. Slack absorbs hash quality
// variation, but a broken Kirsch-Mitzenmacher derivation will blow well
// past 2x.
func TestStage02_FPRateBelowTheoreticalBound(t *testing.T) {
	const (
		m      uint64 = 95850
		k      uint64 = 7
		n             = 10000
		trials        = 100000
	)
	f := bloom.NewWithK(m, k)
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
			"FP rate %.4f exceeds 2x theoretical %.4f (bound %.4f)",
			observed, theoretical, bound,
		)
	}
}
