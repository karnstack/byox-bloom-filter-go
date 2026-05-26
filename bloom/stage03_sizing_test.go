// Tests for Stage 3: optimal sizing math.
//
// Read the stage on karnstack.com/build/bloom-filter/03-optimal-sizing-math
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section.

package bloom_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage03_OptimalSizeAt10kKeysOnePercent checks the canonical case
// (n=10000, p=0.01). Textbook closed-form gives m=95850.6 and k=6.93;
// implementations should round m up to an integer and k to the nearest
// integer. We allow 5% slack on m and +-1 on k to absorb rounding
// strategy variation.
func TestStage03_OptimalSizeAt10kKeysOnePercent(t *testing.T) {
	const (
		n           uint64  = 10000
		p           float64 = 0.01
		expectedM   uint64  = 95850
		expectedK   int     = 7
		mSlack              = 0.05
		kTolerance          = 1
	)
	m, k := bloom.OptimalSize(n, p)
	if m == 0 {
		t.Fatal("OptimalSize returned m=0; check that the formula uses ln(p), not log10(p)")
	}
	if k == 0 {
		t.Fatal("OptimalSize returned k=0; check that you round to the nearest integer, not floor")
	}
	mDelta := math.Abs(float64(int64(m)-int64(expectedM))) / float64(expectedM)
	if mDelta > mSlack {
		t.Fatalf("m = %d, want within %.0f%% of %d (delta %.2f%%)",
			m, mSlack*100, expectedM, mDelta*100)
	}
	if k < expectedK-kTolerance || k > expectedK+kTolerance {
		t.Fatalf("k = %d, want %d +- %d", k, expectedK, kTolerance)
	}
}

// TestStage03_OptimalSizeAtVariedRates spot-checks across a small grid of
// (n, p) values. The numbers come from the closed-form. Implementations
// that swap natural log for base-10 log here are caught by the m delta.
func TestStage03_OptimalSizeAtVariedRates(t *testing.T) {
	cases := []struct {
		n         uint64
		p         float64
		wantBitsPerKey float64 // m/n, derived from -ln(p)/(ln 2)^2
	}{
		{1000, 0.01, 9.585},
		{1000, 0.001, 14.377},
		{1000, 0.0001, 19.170},
		{100000, 0.01, 9.585},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("n=%d_p=%g", c.n, c.p), func(t *testing.T) {
			m, _ := bloom.OptimalSize(c.n, c.p)
			gotBitsPerKey := float64(m) / float64(c.n)
			delta := math.Abs(gotBitsPerKey-c.wantBitsPerKey) / c.wantBitsPerKey
			if delta > 0.05 {
				t.Errorf("bits/key = %.3f, want %.3f (delta %.2f%%)",
					gotBitsPerKey, c.wantBitsPerKey, delta*100)
			}
		})
	}
}

// TestStage03_OptimalSizeProducesUsableFilter wires the output back into
// NewWithK and confirms the resulting filter sits at or below the target
// false-positive rate (within 2x slack to absorb finite-trial noise).
func TestStage03_OptimalSizeProducesUsableFilter(t *testing.T) {
	const (
		n      uint64  = 10000
		p      float64 = 0.01
		trials         = 100000
	)
	m, k := bloom.OptimalSize(n, p)
	f := bloom.NewWithK(m, uint64(k))
	for i := uint64(0); i < n; i++ {
		f.Add([]byte(fmt.Sprintf("added-%d", i)))
	}
	fps := 0
	for i := 0; i < trials; i++ {
		if f.Test([]byte(fmt.Sprintf("query-%d", i))) {
			fps++
		}
	}
	observed := float64(fps) / float64(trials)
	if observed > 2*p {
		t.Fatalf("FP rate %.4f exceeds 2x target %.4f", observed, p)
	}
}
