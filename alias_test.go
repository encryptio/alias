// Copyright (c) 2012, Jack Christopher Kastorff
// All rights reserved.
// BSD Licensed, see LICENSE for details.

package alias

import (
	"math"
	"math/rand"
	"testing"
)

const distributionCount = 1000000
const errorBound = 0.001

func testDistribution(t *testing.T, dist []float64, seed int64) {
	sum := float64(0)
	for i := 0; i < len(dist); i++ {
		sum += dist[i]
	}

	a, err := New(dist)
	if err != nil {
		t.Error("Got an error during creation:", err)
	}

	rng := rand.New(rand.NewSource(seed))

	counts := make([]int64, len(dist))
	for i := 0; i < distributionCount; i++ {
		counts[a.Gen(rng)]++
	}

	for i := 0; i < len(dist); i++ {
		p := float64(counts[i]) / distributionCount
		if math.Abs(p-dist[i]/sum) > errorBound {
			t.Error("Distribution did not match, seed", seed, "- got ", p, "expected", dist[i]/sum)
		}
	}
}

func TestDistribution(t *testing.T) {
	testDistribution(t, []float64{1, 1}, 1)
	testDistribution(t, []float64{1, 2, 3}, 2)
	testDistribution(t, []float64{9, 8, 1, 4, 2}, 5)
	testDistribution(t, []float64{1000, 1, 3, 10}, 39)
	testDistribution(t, []float64{1000, 1, 3, 10}, 61)
}
