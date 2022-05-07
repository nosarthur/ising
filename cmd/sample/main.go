package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	arg "github.com/alexflint/go-arg"
	"github.com/nosarthur/ising/magnets"
	"gopkg.in/yaml.v3"
)

func main() {
	var args struct {
		Input  string `arg:"positional" help:"Input Yaml file"`
		Output string `arg:"positional" default:"out.dat" help:"Output raw spin states"`
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
	fmt.Println(m)

	fmt.Println("exact F U TS:", m.F(), m.U(), m.U()-m.F())
	fmt.Println("exact M:", m.M())

	// Monte Carlo sweeps
	got := magnets.Monte(m, params.NSteps)

	// save spin configurations
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, got)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	if err = os.WriteFile(args.Output, buf.Bytes(), 0666); err != nil {
		p.Fail("Fail to write spin configurations to file")
	}
	fmt.Println("Write spin states to", args.Output)

	// save properties
	csvOut, err := os.Create("tmp.dat")
	if err != nil {
		p.Fail("Unable to open output file")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

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
	fmt.Println("Write properties to", "tmp.dat")
}
