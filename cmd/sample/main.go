package main

import (
	"fmt"

	"github.com/nosarthur/ising/magnets"
	// "gopkg.in/yaml.v3"
)

func main() {
	m := magnets.New1DIsing(10, 1, 1)
	m.Show()
	fmt.Println(m.M())
}
