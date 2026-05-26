// Tests for Stage 5: concurrent-safe Add via atomic OR.
//
// Read the stage on karnstack.com/build/bloom-filter/05-concurrent-add
// before debugging failures. Each test catches a specific failure mode
// listed in the "Common Pitfalls" section. The race detector
// (`go test -race`) is the primary safety net for this stage; mise.toml
// enables it by default.

package bloom_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/karnstack/byox-bloom-filter-go/bloom"
)

// TestStage05_ConcurrentAddAllKeysPresent spawns 100 goroutines that each
// Add a distinct slice of 100 keys (10000 total). After they all join,
// every added key must Test true. The race detector picks up unsynced
// writes to the bit array; the assertion picks up dropped Adds.
func TestStage05_ConcurrentAddAllKeysPresent(t *testing.T) {
	const (
		workers     = 100
		keysPerWorker = 100
	)
	f := bloom.NewWithK(95850, 7)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < keysPerWorker; i++ {
				key := []byte(fmt.Sprintf("w%d-k%d", w, i))
				f.Add(key)
			}
		}(w)
	}
	wg.Wait()

	for w := 0; w < workers; w++ {
		for i := 0; i < keysPerWorker; i++ {
			key := []byte(fmt.Sprintf("w%d-k%d", w, i))
			if !f.Test(key) {
				t.Errorf("Test(%q) = false after concurrent Add", key)
			}
		}
	}
}

// TestStage05_ConcurrentAddOnBlockedFilter exercises the same property on
// the blocked layout. Stage 4 introduced a different memory layout; stage
// 5's atomic OR has to work for both. Failure here means stage 5 only
// atomicized the flat path.
func TestStage05_ConcurrentAddOnBlockedFilter(t *testing.T) {
	const workers = 50
	f := bloom.NewBlocked(95850, 7)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				f.Add([]byte(fmt.Sprintf("blocked-w%d-k%d", w, i)))
			}
		}(w)
	}
	wg.Wait()

	for w := 0; w < workers; w++ {
		for i := 0; i < 100; i++ {
			key := []byte(fmt.Sprintf("blocked-w%d-k%d", w, i))
			if !f.Test(key) {
				t.Errorf("Test(%q) = false after concurrent Add on blocked filter", key)
			}
		}
	}
}
