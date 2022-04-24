package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nosarthur/ising/magnets"
	// "gopkg.in/yaml.v3"
)

func main() {
	J := 2.
	m := magnets.New1DIsing(10, -J, 1)
	fmt.Println("anti ferro:", m.F(), m.U(), m.U()-m.F())
	fmt.Println(m.M())
	m.Show()

	// return
	// time evolution
	csvOut, err := os.Create("./em.dat")
	if err != nil {
		log.Fatal("Unable to open output")
	}
	defer csvOut.Close()
	w := csv.NewWriter(csvOut)
	w.Comma = '\t'
	defer w.Flush()

	got := magnets.Monte(m, 1000)
	for i := range got {
		m.SetRaw(got[i])
		e := strconv.FormatFloat(m.E(), 'f', 2, 64)
		m := strconv.FormatFloat(m.S(), 'f', 2, 64)
		if err := w.Write([]string{strconv.Itoa(i), e, m}); err != nil {
			log.Fatal("Cannot write to file")
		}
	}
}
