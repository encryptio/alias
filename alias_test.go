// Copyright (c) 2012, Jack Christopher Kastorff
// All rights reserved.
// BSD Licensed, see LICENSE for details.

package alias

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	stat "github.com/ematvey/gostat"
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
		return
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

func TestTail(t *testing.T) {
	const size = 33294320
	const half = size / 2
	const tries = 1000000
	const alpha = 0.05
	dist := make([]float64, size)
	for i := range dist {
		dist[i] = 1
	}
	a, err := New(dist)
	if err != nil {
		t.Fatalf("Got an error during creation:", err)
	}
	rng := rand.New(rand.NewSource(42))
	var k int64 // [0,half) (one half)
	for i := 0; i < tries; i++ {
		if a.Gen(rng) < half {
			k++
		}
	}
	// Expected probability of getting k <= k_observed if p == 0.5.
	p := stat.Binomial_CDF_At(0.5, tries, k)
	if p < alpha/2 || p > (1-alpha/2) {
		t.Errorf("The distribution is biased. %d of %d results were in the first half. Binomial_CDF = %f", k, tries, p)
	}
}

func TestBalanceInsideBucket(t *testing.T) {
	const size = 33294320
	//const size = 8
	const half = size / 2
	const tries = 1000000
	const alpha = 0.05
	dist := make([]float64, size)
	for i := range dist {
		if i < half {
			dist[i] = 1
		} else {
			dist[i] = 3
		}
	}
	a, err := New(dist)
	if err != nil {
		t.Fatalf("Got an error during creation:", err)
	}
	rng := rand.New(rand.NewSource(421))
	var k int64 // [0,half) (one half)
	for i := 0; i < tries; i++ {
		if a.Gen(rng) < half {
			k++
		}
	}
	// Expected probability of getting k <= k_observed if p == 0.5.
	p := stat.Binomial_CDF_At(0.25, tries, k)
	if p < alpha/2 || p > (1-alpha/2) {
		t.Errorf("The distribution is biased. %d of %d results were in the first half. Binomial_CDF = %f", k, tries, p)
	}
}

func TestMarshalBinary(t *testing.T) {
	distributions := [][]float64{
		{1},
		{1, 1},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 1000},
	}
	for _, distribution := range distributions {
		a, err := New(distribution)
		if err != nil {
			t.Fatalf("Couldn't create alias: %v", err)
		}

		data, err := a.MarshalBinary()
		if err != nil {
			t.Fatalf("Couldn't MarshalBinary: %v", err)
		}

		a2 := &Alias{}
		err = a2.UnmarshalBinary(data)
		if err != nil {
			t.Fatalf("Couldn't UnmarshalBinary: %v", err)
		}

		if !reflect.DeepEqual(a, a2) {
			t.Fatalf("Unmarshalled version was not the same as original")
		}
	}
}
