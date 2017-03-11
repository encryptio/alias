// Copyright (c) 2012-2015, Jack Christopher Kastorff
// All rights reserved.
// BSD Licensed, see LICENSE for details.

// The alias package picks items from a discrete distribution
// efficiently using the alias method.
package alias

import (
	"encoding/binary"
	"errors"
	"math/rand"
)

type Alias struct {
	table []ipiece
	maxRi uint32
	maxRj uint32
	avgP  uint32
	dummy uint32
}

type fpiece struct {
	prob  float64
	alias uint32
}

type ipiece struct {
	prob  uint32 // [0,2^31)
	alias uint32
}

func calcMax(n uint32) uint32 {
	// Taken from math/rand.Rand.Int31n source.
	return (1 << 31) - 1 - (1<<31)%uint32(n)
}

// checkAvgP checks assumption that piece with prob = al.avgP - 1 exists.
// This assumption is used in UnmarshalBinary.
func checkAvgP(al *Alias) {
	for _, p := range al.table {
		if p.prob == al.avgP-1 {
			return
		}
	}
	panic("Internal error: no piece with prob = al.avgP-1 found.")
}

// Create a new alias object.
// For example,
//   var v = alias.New([]float64{8,10,2})
// creates an alias that returns 0 40% of the time, 1 50% of the time, and
// 2 10% of the time.
func New(prob []float64) (*Alias, error) {

	// This implementation is based on
	// http://www.keithschwarz.com/darts-dice-coins/

	n := len(prob)

	if n < 1 {
		return nil, errors.New("too few probabilities")
	}

	if int(uint32(n)) != n {
		return nil, errors.New("too many probabilities")
	}

	total := float64(0)
	for _, v := range prob {
		if v <= 0 {
			return nil, errors.New("a probability is non-positive")
		}
		total += v
	}

	var al Alias
	al.table = make([]ipiece, n)

	// Michael Vose's algorithm

	// "small" stack grows from the bottom of this array
	// "large" stack from the top
	twins := make([]fpiece, n)

	smTop := -1
	lgBot := n

	// invariant: smTop < lgBot, that is, the twin stacks don't collide

	mult := float64(n) / total
	for i, p := range prob {
		p = p * mult

		// push large items (>=1 probability) into the large stack
		// others in the small stack
		if p >= 1 {
			lgBot--
			twins[lgBot] = fpiece{p, uint32(i)}
		} else {
			smTop++
			twins[smTop] = fpiece{p, uint32(i)}
		}
	}

	if lgBot != smTop+1 {
		panic("alias.New: internal error")
	}

	for smTop >= 0 && lgBot < n {
		// pair off a small and large block, taking the chunk from the large block that's wanted
		l := twins[smTop]
		smTop--

		g := twins[lgBot]
		lgBot++

		al.table[l.alias].prob = uint32(l.prob * (1<<31 - 1))
		al.table[l.alias].alias = g.alias

		g.prob = (g.prob + l.prob) - 1

		// put the rest of the large block back in a list
		if g.prob < 1 {
			smTop++
			twins[smTop] = g
		} else {
			lgBot--
			twins[lgBot] = g
		}
	}

	// clear out any remaining blocks
	for i := n - 1; i >= lgBot; i-- {
		al.table[twins[i].alias].prob = 1<<31 - 1
	}

	// there shouldn't be anything here, but sometimes floating point
	// errors send a probability just under 1.
	for i := 0; i <= smTop; i++ {
		al.table[twins[i].alias].prob = 1<<31 - 1
	}

	al.avgP = (1 << 31)
	al.maxRi = calcMax(uint32(len(al.table)))
	al.maxRj = calcMax(al.avgP)
	al.dummy = uint32(len(al.table))

	checkAvgP(&al)

	return &al, nil
}

// Create a new alias object with integer weights.
// For example,
//   var v = alias.NewInt([]int32{8,10,2})
// creates an alias that returns 0 40% of the time, 1 50% of the time, and
// 2 10% of the time.
func NewInt(prob []int32) (*Alias, error) {

	// This implementation is based on
	// http://www.keithschwarz.com/darts-dice-coins/

	n := len(prob)

	if n < 1 {
		return nil, errors.New("too few probabilities")
	}

	if int(uint32(n)) != n {
		return nil, errors.New("too many probabilities")
	}

	total := uint64(0)
	for _, v := range prob {
		if v <= 0 {
			return nil, errors.New("a probability is non-positive")
		}
		total += uint64(v)
	}

	var al Alias
	al.dummy = uint32(n)

	avg := uint32(total / uint64(n))
	if uint64(avg)*uint64(n) < total {
		// Really add dummy. Its index is already set to al.dummy.
		avg += 1
		n += 1
		dummy := uint64(avg)*uint64(n) - total
		total += dummy
		prob = append(prob, int32(dummy))
	}

	// Michael Vose's algorithm

	// "small" stack grows from the bottom of this array
	// "large" stack from the top
	twins := make([]ipiece, n)

	smTop := -1
	lgBot := n

	// invariant: smTop < lgBot, that is, the twin stacks don't collide

	for i, p := range prob {

		// push large items (>=1 probability) into the large stack
		// others in the small stack
		if uint32(p) >= avg {
			lgBot--
			twins[lgBot] = ipiece{uint32(p), uint32(i)}
		} else {
			smTop++
			twins[smTop] = ipiece{uint32(p), uint32(i)}
		}
	}

	if lgBot != smTop+1 {
		panic("alias.New: internal error")
	}

	al.table = make([]ipiece, n)

	for smTop >= 0 && lgBot < n {
		// pair off a small and large block, taking the chunk from the large block that's wanted
		l := twins[smTop]
		smTop--

		g := twins[lgBot]
		lgBot++

		al.table[l.alias].prob = l.prob - 1
		al.table[l.alias].alias = g.alias

		g.prob = (g.prob + l.prob) - avg

		// put the rest of the large block back in a list
		if g.prob < avg {
			smTop++
			twins[smTop] = g
		} else {
			lgBot--
			twins[lgBot] = g
		}
	}

	// clear out any remaining blocks
	for i := n - 1; i >= lgBot; i-- {
		al.table[twins[i].alias].prob = avg - 1
	}

	if smTop != -1 {
		panic("alias.NewInt: internal error")
	}

	al.avgP = avg
	al.maxRi = calcMax(uint32(n))
	al.maxRj = calcMax(al.avgP)

	checkAvgP(&al)

	return &al, nil
}

// Generates a random number according to the distribution using the rng passed.
func (al *Alias) Gen(rng *rand.Rand) uint32 {
begin:
	r := rng.Int63()
	ri := uint32(r & (1<<31 - 1))
	rj := uint32((r >> 31) & (1<<31 - 1))
	if ri > al.maxRi || rj > al.maxRj {
		goto begin
	}
	w := ri % uint32(len(al.table))
	x := rj % al.avgP
	if x > al.table[w].prob {
		w = al.table[w].alias
	}
	if w == al.dummy {
		goto begin
	}
	return w
}

// MarshalBinary implements encoding.BinaryMarshaller.
func (al *Alias) MarshalBinary() ([]byte, error) {
	out := make([]byte, len(al.table)*8, len(al.table)*8+4)
	for i, piece := range al.table {
		bin := out[i*8 : 8+i*8]
		binary.LittleEndian.PutUint32(bin[0:4], piece.prob)
		binary.LittleEndian.PutUint32(bin[4:8], piece.alias)
	}
	if al.dummy != uint32(len(al.table)) {
		dummy := make([]byte, 4)
		binary.LittleEndian.PutUint32(dummy, al.dummy)
		out = append(out, dummy...)
	}
	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaller.
func (al *Alias) UnmarshalBinary(p []byte) error {
	if len(p)%4 != 0 {
		return errors.New("bad data length")
	}

	if int(uint32(len(p)/8)) != len(p)/8 {
		return errors.New("data too large")
	}

	al.table = make([]ipiece, (len(p))/8)
	for i := range al.table {
		bin := p[i*8 : 8+i*8]
		prob := binary.LittleEndian.Uint32(bin[0:4])
		alias := binary.LittleEndian.Uint32(bin[4:8])

		if prob >= 1<<31 {
			return errors.New("bad data: probability out of range")
		}
		if alias >= uint32(len(al.table)) {
			return errors.New("bad data: alias target out of range")
		}

		al.table[i].prob = prob
		al.table[i].alias = alias
	}

	// How we calculate avgP... There must be at least one piece with
	// probability of 1.0. TODO: prove it. In our model such a probability
	// is represented as avgP - 1. (-1 comes from the way we compare prob
	// with random.) So avgP = max(prob) + 1.
	maxProb := uint32(0)
	for _, piece := range al.table {
		if piece.prob > maxProb {
			maxProb = piece.prob
		}
	}
	al.avgP = maxProb + 1

	al.maxRi = calcMax(uint32(len(al.table)))
	al.maxRj = calcMax(al.avgP)

	al.dummy = uint32(len(al.table))
	if len(p)%8 != 0 {
		dummy := p[len(p)-4:]
		al.dummy = binary.LittleEndian.Uint32(dummy)
	}

	return nil
}
