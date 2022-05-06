package magnets

import (
	"fmt"
	"testing"

	"github.com/nosarthur/ising/magnets"
	"github.com/stretchr/testify/assert"
)

func TestFlip(t *testing.T) {
	// 10 spins
	m := magnets.New1DIsing(4, 2, 1)
	m.TryFlip(2, 0)
	fmt.Println(m.Raw())
	assert.Equal(t, m.Raw(), magnets.Spint(55))
	m.TryFlip(2, 0.9)
	fmt.Println(m.Raw())
	assert.Equal(t, m.Raw(), magnets.Spint(63))
}
