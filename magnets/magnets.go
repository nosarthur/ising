package magnets

import (
	"fmt"
	"math/bits"
)

type Magnet interface {
	// Step()
	Energy() float32
	Magnetization() int
	Show()
}

type magnet struct {
	spins uint32
	N     int // number of spins
	J, H  float64
}

func New1DIsing() Magnet {
	return &magnet{111, 10, 1, 1}

}

func (m *magnet) Magnetization() int {
	// 2*n1 - N
	return 2*bits.OnesCount32(m.spins) - m.N
}

func (m *magnet) Energy() float32 {
	return 0
}

func (m *magnet) Show() {
	// TODO: show all spins
	fmt.Printf("%d: %b\n", m.N, m.spins)
	fmt.Printf("J: %0.2f, H: %0.2f\n", m.J, m.H)
}
