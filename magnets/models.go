// Package magnets implements Ising spin models
package magnets

import (
	"fmt"
)

//
type Magnet interface {
	TryFlip(uint, float64)
	E() float64
	S() float64 // micro state magnetization
	N() uint    // number of spins
	Set(Spint)
	SetRaw(Spint)
	Raw() Spint
}

// 1D Ising magnet with PBC
type ising1d struct {
	// the size of Spins needs to be longer than N + 2 to deal with PBC
	spins    Spint      // spin configuration
	n        uint       // number of spins
	j        float64    // spin-spin interaction, in unit of kT
	h        float64    // external magnetic field, in unit of kT
	z        float64    // partition function
	masks1   []Spint    // select single spin, i in 0..N-1
	masks2   []Spint    // select 2 neighboring spins (i,i+1)
	masks3   []Spint    // select 3 neighboring spins (i-1, i, i+1), i in 1..N
	mask     Spint      // mask all spins
	maskL    Spint      // mask for the bit left to the spins, image of spin 0
	dE_table [8]float64 // local 3 spins => dE
	Bound    Spint      // upper bound for looping over the spin states
}

// Create new system (spin state not initialized)
func New1DIsing(N uint, J, H float64) *ising1d {
	// TODO: more initialization methods?
	// TODO: E and M can both be cached

	// Pad the spins from both ends. Thus the indices are from 1..N
	if N+2 > spint_bits {
		return nil
	}

	var m, shift Spint
	m1 := make([]Spint, N)
	m2 := make([]Spint, N)
	m3 := make([]Spint, N)
	for i := range m1 {
		shift = Spint(i)
		m1[i] = 2 << shift
		m2[i] = 3 << shift
		m3[i] = 7 << shift
		m |= m1[i]
	}
	bound := Spint(1) << N
	// -J \sum Si Sj -H \sum Si
	var de_table [8]float64
	dJ := 4 * J
	dH := 2 * H
	de_table[0] = dJ - dH // 000 => 010: 2 => -2, -3 => -1
	de_table[1] = -dH
	de_table[2] = -dJ + dH
	de_table[3] = dH
	de_table[4] = -dH
	de_table[5] = -dJ - dH
	de_table[6] = dH
	de_table[7] = dJ + dH

	mag := &ising1d{n: N, j: J, h: H,
		masks1: m1, masks2: m2, masks3: m3, mask: m, maskL: Spint(2) << N,
		dE_table: de_table, Bound: bound}
	mag.updateZ()
	return mag
}

// Set spin state
func (m *ising1d) Set(spins Spint) {
	spins <<= 1
	spins &= m.mask
	// set index 0 and N for PBC
	spins |= (spins & m.masks1[0]) << m.n
	spins |= (spins & m.masks1[m.n-1]) >> m.n
	m.spins = spins
}

// Show spin state and system parameters
func (m *ising1d) String() string {
	// TODO: show all spins
	spin := (m.spins & m.mask) >> 1
	return fmt.Sprintf("%d spins: %b; J=%0.2f, H=%0.2f\n",
		m.n, spin, m.j, m.h)
}

// Set raw spin state
func (m *ising1d) SetRaw(spins Spint) {
	m.spins = spins
}

// Get raw spin state
func (m *ising1d) Raw() Spint {
	return m.spins
}
