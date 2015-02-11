package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-hep/rio"
)

func main() {
	var fname string

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	fmt.Printf("::: inspecting file [%s]...\n", fname)
	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	r := rio.NewReader(f)
	if r == nil {
		fmt.Printf("*** error creating rio.Reader\n")
		os.Exit(2)
	}

	scan := rio.NewScanner(r)
	for scan.Scan() {
		// scans through the whole stream
		err = scan.Err()
		if err != nil {
			break
		}
		fmt.Printf(" -> %v\n", scan.Record().Name())
	}
	err = scan.Err()
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("::: inspecting file [%s]... [done]\n", fname)
}
