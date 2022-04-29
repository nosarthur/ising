package main

import (
	"bytes"
	"encoding/binary"
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
		Output string `arg:"positional" default:"props.txt"`
	}
	p := arg.MustParse(&args)
	if args.Input == "" {
		p.Fail("Input yaml file is needed.")
	}

	b, err := ioutil.ReadFile(args.Input)
	if err != nil {
		p.Fail("Fail to read input yaml data file")

	}
	params := magnets.Params{}
	err = yaml.Unmarshal(b, &params)
	if err != nil {
		p.Fail("Fail to parse input yaml file")
	}
	fmt.Println(params)

	b, err = ioutil.ReadFile("sweeps.dat")
	if err != nil {
		p.Fail("Fail to read input spin data file")
	}
	r := bytes.NewReader(b)
	spinConfigs := make([]magnets.Spint, params.NSteps)
	if err := binary.Read(r, binary.LittleEndian, &spinConfigs); err != nil {
		p.Fail("binary.Read failed:")
	}

	// save properties
	csvOut, err := os.Create(args.Output)
	if err != nil {
		p.Fail("Unable to open output file")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

	m := magnets.New1DIsing(params.NSpins, params.J, params.H)
	for i := range spinConfigs {
		m.SetRaw(spinConfigs[i])
		e := strconv.FormatFloat(m.E(), 'f', 2, 64)
		m := strconv.FormatFloat(m.S(), 'f', 2, 64)
		if err := w.Write([]string{strconv.Itoa(i), e, m}); err != nil {
			p.Fail("Cannot write to file")
		}

	}
	fmt.Println("Write to", args.Output)
}
