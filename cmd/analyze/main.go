package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/alexflint/go-arg"
	"github.com/nosarthur/ising/magnets"
	"gopkg.in/yaml.v3"
)

type args struct {
	Yaml    string `arg:"positional" help:"Input Yaml file"`
	MCInput string `arg:"positional" help:"Input spin states file"`
	Output  string `arg:"positional" default:"props.txt"`
}

func (args) Description() string {
	return "This program computes the observables of an ensemble"
}

func main() {
	var args args

	p := arg.MustParse(&args)

	b, err := ioutil.ReadFile(args.Yaml)
	if err != nil {
		p.Fail("Fail to read input yaml data file")

	}
	params := magnets.Params{}
	err = yaml.Unmarshal(b, &params)
	if err != nil {
		p.Fail("Fail to parse input yaml file")
	}
	fmt.Println(params)

	b, err = ioutil.ReadFile(args.MCInput)
	if err != nil {
		p.Fail("Fail to read input spin data file")
	}
	r := bytes.NewReader(b)
	rawSpinConfigs := make([]magnets.Spint, params.NSteps)
	if err := binary.Read(r, binary.LittleEndian, &rawSpinConfigs); err != nil {
		p.Fail("binary.Read failed:")
	}

	// save properties for each spin config
	csvOut, err := os.Create(args.Output)
	if err != nil {
		p.Fail("Unable to open output file")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

	mag := magnets.New1DIsing(params.NSpins, params.J, params.H)
	for i := range rawSpinConfigs {
		mag.SetRaw(rawSpinConfigs[i])
		e := strconv.FormatFloat(mag.E(), 'f', 2, 64)
		m := strconv.FormatFloat(mag.S(), 'f', 2, 64)
		if err := w.Write([]string{strconv.Itoa(i), e, m}); err != nil {
			p.Fail("Cannot write to file")
		}

	}
	fmt.Println("Write to", args.Output)

	// statistics
	toCompute := []func(magnets.Magnet) float64{
		func(m magnets.Magnet) float64 { return m.S() },
		func(m magnets.Magnet) float64 { return m.E() },
		func(m magnets.Magnet) float64 { return math.Exp(-m.E()) },
	}
	for i := range toCompute {
		f := toCompute[i]
		m, dm := magnets.Estimate(f, rawSpinConfigs, mag)
		fmt.Printf("%v +- %v\n", m, dm)
	}

}
