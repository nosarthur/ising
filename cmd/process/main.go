package main

import "github.com/alexflint/go-arg"

func main() {
	var args struct {
		Input  string `arg:"positional" help:"Input .dat file with raw spin configurations"`
		Output string `arg:"positional" default:"out.dat"`
	}
	p := arg.MustParse(&args)
	if args.Input == "" {
		p.Fail("Input .dat file is needed.")
	}

}
