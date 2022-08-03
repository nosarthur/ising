package magnets

import (
	"math"
	"math/rand"
)

// Make Monte Carlo sweeps
func Monte(m Magnet, nMC uint) []Spint {
	// var i, j uint
	samples := make([]Spint, nMC)
	for i := uint(0); i < nMC; i++ {
		for j := uint(0); j < m.N(); j++ {
			m.TryFlip(j, rand.Float64())
		}
		// after one sweep
		samples[i] = m.Raw()
	}
	return samples
}

// Accept or reject the flip of the i'th spin
func (m *ising1d) TryFlip(i uint, rnd float64) {
	local3 := (m.masks3[i] & m.spins) >> i
	if rnd < math.Exp(-m.dE_table[local3]) {
		mask := m.masks1[i]
		m.spins ^= mask
		// PBC: 1; N ... 1; N
		// thus i=0 and i=N-1 requires special treatments
		if i == 0 {
			m.spins ^= m.maskL
		} else if i == m.n {
			m.spins ^= 1
		}
	}
}

//TODO: uniformly sample raw spin state by sampling uniformly random integers
func RandomSample(n uint, m ising1d) []Spint {
	sample := make([]Spint, n)
	bound := int32(m.Bound)
	for i := uint(0); i < n; i++ {
		r := Spint(rand.Int31n(bound))
		m.Set(r)
		sample[i] = m.Raw()
	}
	return sample
}
