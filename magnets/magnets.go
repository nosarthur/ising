package magnets

import (
	"fmt"
)

type Magnet interface {
	// Step()
	DeltaE(i uint) float64
	E() float64
	M() float64
	Set(Spint)
	Show()
}

// 1D Ising magnet with PBC
type ising1d struct {
	// the size of Spins needs to be longer than N + 2 to deal with PBC
	spins    Spint   // spin configuration
	n        uint    // number of spins
	j        float64 // spin-spin interaction, in unit of kT
	h        float64 // external magnetic field, in unit of kT
	z        float64 // partition function
	masks1   []Spint // select single spin, i in 0..N-1
	masks2   []Spint // select 2 neighboring spins (i,i+1)
	masks3   []Spint // select 3 neighboring spins (i-1, i, i+1), i in 1..N
	mask     Spint   // mask all spins
	dE_table [8]float64
	Bound    Spint // upper bound for looping over the spin states
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
	fmt.Println("bound", bound)
	// -J \sum Si Sj
	var de_table [8]float64
	tmp := 4 * J
	de_table[0] = tmp // 000 => 010: 2 => -2
	de_table[7] = tmp
	de_table[2] = -tmp
	de_table[5] = -tmp

	mag := &ising1d{n: N, j: J, h: H, masks1: m1, masks2: m2, masks3: m3,
		mask: m, dE_table: de_table, Bound: bound}
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
func (m *ising1d) Show() {
	// TODO: show all spins
	fmt.Printf("%d spins: %b\n", m.n, (m.spins&m.mask)>>1)
	fmt.Printf("J: %0.2f, H: %0.2f\n", m.j, m.h)
}
