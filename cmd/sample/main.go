package main

import (
	"fmt"

	"github.com/nosarthur/ising/magnets"
	// "gopkg.in/yaml.v3"
)

func main() {
	J := 2.
	m := magnets.New1DIsing(10, J, 1)
	fmt.Println("ferro:     ", m.F(), m.U(), m.U()-m.F())
	fmt.Println(m.M())
	m = magnets.New1DIsing(10, -J, 1)
	fmt.Println("anti ferro:", m.F(), m.U(), m.U()-m.F())
	fmt.Println(m.M())
}
