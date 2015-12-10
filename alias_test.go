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
const errorBound = 0.01

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

func TestBasicCorrectness(t *testing.T) {
	testDistribution(t, []float64{1, 1}, 1)
	testDistribution(t, []float64{1, 2, 3}, 2)
	testDistribution(t, []float64{9, 8, 1, 4, 2}, 5)
	testDistribution(t, []float64{1000, 1, 3, 10}, 39)
	testDistribution(t, []float64{1000, 1, 3, 10}, 61)
}

func benchGen(b *testing.B, size int) {
	b.StopTimer()

	arr := make([]float64, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Float64()
	}

	a, err := New(arr)
	if err != nil {
		b.Error("Got an error during creation:", err)
	}

	rng := rand.New(rand.NewSource(99))

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a.Gen(rng)
	}
}

func BenchmarkGen5(b *testing.B) {
	benchGen(b, 5)
}

func BenchmarkGen50(b *testing.B) {
	benchGen(b, 50)
}

func BenchmarkGen500(b *testing.B) {
	benchGen(b, 500)
}

func BenchmarkGen5000(b *testing.B) {
	benchGen(b, 5000)
}

func BenchmarkGen50000(b *testing.B) {
	benchGen(b, 50000)
}

func benchCreationSize(b *testing.B, size int) {
	b.StopTimer()

	arr := make([]float64, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Float64()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		New(arr)
	}
}

func BenchmarkCreate5(b *testing.B) {
	benchCreationSize(b, 5)
}

func BenchmarkCreate50(b *testing.B) {
	benchCreationSize(b, 50)
}

func BenchmarkCreate500(b *testing.B) {
	benchCreationSize(b, 500)
}

func BenchmarkCreate5000(b *testing.B) {
	benchCreationSize(b, 5000)
}

func BenchmarkCreate50000(b *testing.B) {
	benchCreationSize(b, 50000)
}
