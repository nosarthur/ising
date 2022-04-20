package magnets

import (
	"fmt"
	"math/bits"
)

type Magnet interface {
	// Step()
	DeltaE(i uint) float64
	Energy() float64
	Magnetization() float64
	Set(spint)
	Show()
}

// 1D Ising magnet with PBC
type magnet struct {
	// the size of Spins needs to be longer than N + 2 to deal with PBC
	spins    spint // spin configuration
	N        uint  // number of spins
	J, H     float64
	masks1   []spint // select single spin, i in 0..N-1
	masks2   []spint // select 2 neighboring spins (i,i+1)
	masks3   []spint // select 3 neighboring spins (i-1, i, i+1), i in 1..N
	mask     spint   // mask all spins
	dE_table [8]float64
}

// Create new system (spin state not initialized)
func New1DIsing(spins spint, N uint, J, H float64) *magnet {
	// TODO: more initialization methods?
	// Pad the spins from both ends. Thus the indices are from 1..N
	if N+2 > spint_bits {
		return nil
	}

	var m, shift spint
	m1 := make([]spint, N)
	m2 := make([]spint, N)
	m3 := make([]spint, N)
	for i := range m1 {
		shift = spint(i)
		m1[i] = 2 << shift
		m2[i] = 3 << shift
		m3[i] = 7 << shift
		m |= m1[i]
	}
	// -J \sum Si Sj
	var de_table [8]float64
	tmp := 4 * J
	de_table[0] = tmp // 000 => 010: 2 => -2
	de_table[7] = tmp
	de_table[2] = -tmp
	de_table[5] = -tmp
	mag := &magnet{N: N, J: J, H: H, masks1: m1, masks2: m2, masks3: m3,
		mask: m, dE_table: de_table}
	mag.Set(spins)
	return mag
}

func (m *magnet) Magnetization() float64 {
	// 2*n1 - N
	return float64(2*bits.OnesCount16(m.spins&m.mask) - int(m.N))
}

// Return the energy difference to flip the i'th spin
func (m *magnet) DeltaE(i uint) float64 {
	// idx + 1 % N == i
	return 0
}

func (m *magnet) Energy() float64 {
	// E = -J \sum Si Sj - H \sum Si

	var pair spint
	var e float64
	for i := range m.masks1 {
		// Note here the first pair is spin 1 and spin N
		pair = (m.spins & m.masks2[i]) >> i
		if pair == 1 || pair == 2 {
			e -= 1
		} else {
			e += 1
		}
	}
	// fmt.Println("e:", e, m.Magnetization())
	return -(e*m.J + m.H*m.Magnetization())
}

// Set spin state
func (m *magnet) Set(spins spint) {
	spins <<= 1
	spins &= m.mask
	// set index 0 and N for PBC
	spins |= (spins & m.masks1[0]) << m.N
	spins |= (spins & m.masks1[m.N-1]) >> m.N
	m.spins = spins
}

// Show spin state and system parameters
func (m *magnet) Show() {
	// TODO: show all spins
	fmt.Printf("%d: %b\n", m.N, m.spins)
	fmt.Printf("J: %0.2f, H: %0.2f\n", m.J, m.H)
}
