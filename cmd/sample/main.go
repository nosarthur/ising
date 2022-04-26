package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	arg "github.com/alexflint/go-arg"
	"github.com/nosarthur/ising/magnets"
	"gopkg.in/yaml.v3"
)

type Params struct {
	NSpins uint
	J, H   float64
	NSteps uint
}

func main() {
	var args struct {
		Input  string `arg:"positional" help:"Input Yaml file"`
		Output string `arg:"positional" default:"out.dat"`
	}
	p := arg.MustParse(&args)
	if args.Input == "" {
		p.Fail("Input yaml file is needed.")
	}

	yfile, err := ioutil.ReadFile(args.Input)

	if err != nil {
		p.Fail("Fail to read input file")

	}
	params := Params{}
	err = yaml.Unmarshal(yfile, &params)
	if err != nil {
		p.Fail("Fail to parse input yaml file")
	}

	J := 2.
	m := magnets.New1DIsing(10, J, 1)
	fmt.Println("F U TS:", m.F(), m.U(), m.U()-m.F())
	fmt.Println("M:", m.M())
	m.Show()

	// return
	// time evolution
	csvOut, err := os.Create(args.Output)
	if err != nil {
		p.Fail("Unable to open output file")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

	got := magnets.Monte(m, 1000)
	if err := w.Write([]string{"Index", "Energy", "Magnetization"}); err != nil {
		p.Fail("Cannot write to file")
	}
	for i := range got {
		m.SetRaw(got[i])
		e := strconv.FormatFloat(m.E(), 'f', 2, 64)
		m := strconv.FormatFloat(m.S(), 'f', 2, 64)
		if err := w.Write([]string{strconv.Itoa(i), e, m}); err != nil {
			p.Fail("Cannot write to file")
		}
	}
	fmt.Println("Write to", args.Output)
}
