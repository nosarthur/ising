package magnets_test

import (
	"testing"

	"github.com/nosarthur/ising/magnets"
	"github.com/stretchr/testify/assert"
)

func TestEandM(t *testing.T) {
	// 10 spins, H=1
	cases := []struct {
		J      float64
		spins  magnets.Spint
		want_s float64
		want_e float64
	}{
		{1, 0, -10, 0},
		{1, 3, -6, 0},
		{1, 5, -6, 4},
		{2, 1, -8, -4},
	}
	for _, c := range cases {
		m := magnets.New1DIsing(10, c.J, 1)
		m.Set(c.spins)
		got_s := m.S()
		assert.InDelta(t, c.want_s, got_s, 0.1)
		got_e := m.E()
		assert.InDelta(t, c.want_e, got_e, 0.1)
	}
}

func TestP(t *testing.T) {
	cases := []struct {
		n_spins uint
		want_ps []float64
	}{
		{2, []float64{0.018, 0.002, 0.002, 0.977}},
		{3, []float64{0.002, 0, 0, 0.002, 0, 0.002, 0.002, 0.989}},
	}
	for _, c := range cases {
		m := magnets.New1DIsing(c.n_spins, 1, 1)
		for i := magnets.Spint(0); i < m.Bound; i++ {
			m.Set(i)
			assert.InDelta(t, c.want_ps[i], m.P(), 0.01)
		}
	}
}
