package magnets_test

import (
	"testing"

	"github.com/nosarthur/ising/magnets"
	"github.com/stretchr/testify/assert"
)

func TestEandM(t *testing.T) {
	// 10 spins
	m := magnets.New1DIsing(0, 10, 1, 1)
	cases := []struct {
		spins  magnets.Spint
		want_m float64
		want_e float64
	}{
		{0, -10, 0},
		{3, -6, 0},
		{5, -6, 4},
	}
	for _, c := range cases {
		m.Set(c.spins)
		got_m := m.M()
		assert.InDelta(t, c.want_m, got_m, 0.1)
		got_e := m.E()
		assert.InDelta(t, c.want_e, got_e, 0.1)
	}
}