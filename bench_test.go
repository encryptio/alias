// Copyright (c) 2012, Jack Christopher Kastorff
// All rights reserved.
// BSD Licensed, see LICENSE for details.

package alias

import (
	"math/rand"
	"testing"
)

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

func benchGenInt(b *testing.B, size int, maxP int32) {
	b.StopTimer()
	arr := make([]int32, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Int31n(maxP)
	}
	a, err := NewInt(arr)
	if err != nil {
		b.Error("Got an error during creation:", err)
	}
	rng := rand.New(rand.NewSource(99))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.Gen(rng)
	}
}

func BenchmarkGenIntLargeP5(b *testing.B) {
	benchGenInt(b, 5, (1<<31)-1)
}

func BenchmarkGenIntLargeP50(b *testing.B) {
	benchGenInt(b, 50, (1<<31)-1)
}

func BenchmarkGenIntLargeP500(b *testing.B) {
	benchGenInt(b, 500, (1<<31)-1)
}

func BenchmarkGenIntLargeP5000(b *testing.B) {
	benchGenInt(b, 5000, (1<<31)-1)
}

func BenchmarkGenIntLargeP50000(b *testing.B) {
	benchGenInt(b, 50000, (1<<31)-1)
}

func BenchmarkGenIntSmallP5(b *testing.B) {
	benchGenInt(b, 5, 1<<29)
}

func BenchmarkGenIntSmallP50(b *testing.B) {
	benchGenInt(b, 50, 1<<29)
}

func BenchmarkGenIntSmallP500(b *testing.B) {
	benchGenInt(b, 500, 1<<29)
}

func BenchmarkGenIntSmallP5000(b *testing.B) {
	benchGenInt(b, 5000, 1<<29)
}

func BenchmarkGenIntSmallP50000(b *testing.B) {
	benchGenInt(b, 50000, 1<<29)
}

func benchCreationSizeInt(b *testing.B, size int) {
	b.StopTimer()
	arr := make([]int32, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Int31()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewInt(arr)
	}
}

func BenchmarkCreateInt5(b *testing.B) {
	benchCreationSizeInt(b, 5)
}

func BenchmarkCreateInt50(b *testing.B) {
	benchCreationSizeInt(b, 50)
}

func BenchmarkCreateInt500(b *testing.B) {
	benchCreationSizeInt(b, 500)
}

func BenchmarkCreateInt5000(b *testing.B) {
	benchCreationSizeInt(b, 5000)
}

func BenchmarkCreateInt50000(b *testing.B) {
	benchCreationSizeInt(b, 50000)
}
