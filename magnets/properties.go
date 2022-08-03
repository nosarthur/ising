package magnets

import (
	"math"
	"math/bits"
)

// Number of spins
func (m *ising1d) N() uint {
	return m.n
}

// Probability of the current spin state
func (m *ising1d) P() float64 {
	return math.Exp(-m.E()) / m.z
}

// Magnetization of the current spin state
func (m *ising1d) S() float64 {
	// 2*n1 - N
	return float64(2*bits.OnesCount16(m.spins&m.mask) - int(m.n))
}

// Energy of the current spin state
func (m *ising1d) E() float64 {
	// E = -J \sum Si Sj - H \sum Si
	var pair Spint
	var ess float64
	for i := range m.masks1 {
		// Note here the first pair is spin 1 and spin N
		pair = (m.spins & m.masks2[i]) >> i
		if pair == 1 || pair == 2 {
			ess -= 1
		} else {
			ess += 1
		}
	}
	return -(ess*m.j + m.h*m.S())
}

// Free energy
func (m *ising1d) F() float64 {
	return -math.Log(m.z)
}

// Magnetization
func (m *ising1d) M() (s float64) {
	for i := Spint(0); i < m.Bound; i++ {
		m.Set(i)
		s += m.P() * m.S()
	}
	return
}

// Internal energy
func (m *ising1d) U() (u float64) {
	for i := Spint(0); i < m.Bound; i++ {
		m.Set(i)
		u += m.P() * m.E()
	}
	return
}

// re-compute partition function
// It has side effect of setting all spins to 1
func (m *ising1d) updateZ() {
	m.z = 0
	for i := Spint(0); i < m.Bound; i++ {
		m.Set(i)
		m.z += math.Exp(-m.E())
	}
}

// Compute expectation value from a sample
func Estimate(f func(Magnet) float64, sample []Spint, m Magnet) (avg, std float64) {
	var o2 float64
	for i := range sample {
		m.SetRaw(sample[i])
		o := f(m)
		avg += o
		o2 += o * o
	}
	n := float64(len(sample))
	avg /= n
	std = math.Sqrt((o2/n - avg*avg) / (n - 1))
	return
}

// Compute free energy expectation value from a sample drawn on system 1
// using Zwanzig relation
func Zwanzig(m0 Magnet, m1 Magnet, sample1 []Spint) (avg, std float64) {
	f := func(m0 Magnet, m1 Magnet) float64 {
		return math.Exp(m1.E() - m0.E())
	}
	var o, o2 float64
	for i := range sample1 {
		m0.SetRaw(sample1[i])
		m1.SetRaw(sample1[i])
		o = f(m0, m1)
		avg += o
		o2 += o * o
	}
	n := float64(len(sample1))
	avg /= n
	// println("avg", avg)
	avg = -math.Log(avg)
	//TODO: error propagation of the log is not done
	std = math.Sqrt((o2/n - avg*avg) / (n - 1))
	return
}
