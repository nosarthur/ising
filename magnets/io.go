package magnets

import "fmt"

// simulation setting
type Params struct {
	NSpins, NSteps uint
	J, H           float64
	NToss, NRun    uint
}

func (p Params) String() string {
	return fmt.Sprintf("NSpins=%v, NSteps=%v, J=%v, H=%v\n",
		p.NSpins, p.NSteps, p.J, p.H)
}
