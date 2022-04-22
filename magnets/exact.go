package magnets

import (
	"math"
	"math/bits"
)

// Probability of the current spin state
func (m *ising1d) P() float64 {
	return math.Exp(-m.E()) / m.z
}

// Magnetization
func (m *ising1d) M() float64 {
	// 2*n1 - N
	return float64(2*bits.OnesCount16(m.spins&m.mask) - int(m.n))
}

// Return the energy difference to flip the i'th spin
func (m *ising1d) DeltaE(i uint) float64 {
	// idx + 1 % N == i
	return 0
}

// Energy
func (m *ising1d) E() float64 {
	// E = -J \sum Si Sj - H \sum Si

	var pair Spint
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
	return -(e*m.j + m.h*m.M())
}

// Free energy
func (m *ising1d) F() float64 {
	return math.Log(m.z)
}

// re-compute partition function
func (m *ising1d) updateZ() {
	m.z = 0
	for i := Spint(0); i < m.Bound; i++ {
		m.Set(i)
		m.z += math.Exp(-m.E())
		// fmt.Println("--", i, ':', m.E())
	}
}
