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
}

type fpiece struct {
	prob  float64
	alias uint32
}

type ipiece struct {
	prob  uint32 // [0,2^31)
	alias uint32
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

	return &al, nil
}

// Generates a random number according to the distribution using the rng passed.
func (al *Alias) Gen(rng *rand.Rand) uint32 {
	ri := uint32(rng.Int31())
	w := ri % uint32(len(al.table))
	if ri > al.table[w].prob {
		return al.table[w].alias
	}
	return w
}

// MarshalBinary implements encoding.BinaryMarshaller.
func (al *Alias) MarshalBinary() ([]byte, error) {
	out := make([]byte, len(al.table)*8)
	for i, piece := range al.table {
		bin := out[i*8 : 8+i*8]
		binary.LittleEndian.PutUint32(bin[0:4], piece.prob)
		binary.LittleEndian.PutUint32(bin[4:8], piece.alias)
	}
	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaller.
func (al *Alias) UnmarshalBinary(p []byte) error {
	if len(p)%8 != 0 {
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

	return nil
}
