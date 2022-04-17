package main

import (
	"github.com/nosarthur/ising/magnets"
	// "gopkg.in/yaml.v3"
)

func main() {
	m := magnets.New1DIsing()
	m.Show()
}
