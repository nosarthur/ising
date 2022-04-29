package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/alexflint/go-arg"
	"github.com/nosarthur/ising/magnets"
	"gopkg.in/yaml.v3"
)

func main() {
	var args struct {
		Input  string `arg:"positional" help:"Input Yaml file"`
		Output string `arg:"positional" default:"out.txt"`
	}
	p := arg.MustParse(&args)
	if args.Input == "" {
		p.Fail("Input yaml file is needed.")
	}

	b, err := ioutil.ReadFile(args.Input)
	if err != nil {
		p.Fail("Fail to read input yaml file")

	}

	params := magnets.Params{}
	err = yaml.Unmarshal(b, &params)
	if err != nil {
		p.Fail("Fail to parse input yaml file")
	}

	m := magnets.New1DIsing(params.NSpins, params.J, params.H)
	fmt.Println("exact F U TS:", m.F(), m.U(), m.U()-m.F())
	fmt.Println("M:", m.M())
	m.Show()

	// save properties
	csvOut, err := os.Create(args.Output)
	if err != nil {
		p.Fail("Unable to open output file")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

	for i := magnets.Spint(0); i < m.Bound; i++ {
		m.Set(i)
		s := strconv.FormatFloat(m.S(), 'f', 2, 64)
		e := strconv.FormatFloat(m.E(), 'f', 2, 64)
		if err := w.Write([]string{s, e}); err != nil {
			p.Fail("Cannot write to file")
		}
	}
	fmt.Println("Write to", args.Output)

}
