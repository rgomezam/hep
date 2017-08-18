// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2yoda converts ROOT files containing hbook-like values (H1D, H2D, ...)
// into YODA files.
//
// Example:
//
//  $> root2yoda file1.root file2.root > out.yoda
package main // import "go-hep.org/x/hep/cmd/root2yoda"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/rootio"
)

func main() {
	var out io.WriteCloser = os.Stdout
	defer out.Close()

	log.SetFlags(0)
	log.SetPrefix("root2yoda: ")

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root2yoda [options] file1.root [file2.root [...]] > out.yoda

ex:
 $> root2yoda file1.root file2.root > out.yoda

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	for _, fname := range flag.Args() {
		log.Printf("processing %s\n", fname)
		f, err := rootio.Open(fname)
		if err != nil {
			log.Fatalf("failed to open [%s]: %v\n", fname, err)
			os.Exit(1)
		}
		defer f.Close()
		for _, k := range uniq(f.Keys()) {
			walk(out, k)
		}
	}

	err := out.Close()
	if err != nil {
		log.Fatalf("error closing output file: %v\n", err)
	}
}

func walk(w io.Writer, k rootio.Key) {
	obj := k.Value()
	switch obj := obj.(type) {
	case rootio.Directory:
		for _, k := range uniq(obj.Keys()) {
			walk(w, k)
		}
	case yodacnv.Marshaler:
		raw, err := obj.MarshalYODA()
		if err != nil {
			log.Fatalf("root2yoda: error converting %v: %v", k.Name(), err)
		}
		_, err = w.Write(raw)
		if err != nil {
			log.Fatalf("root2yoda: error writing %v: %v", k.Name(), err)
		}
	}
}

func uniq(keys []rootio.Key) []rootio.Key {
	set := make(map[string]rootio.Key, len(keys))
	for _, k := range keys {
		kk, dup := set[k.Name()]
		if dup && kk.Cycle() > k.Cycle() {
			continue
		}
		set[k.Name()] = k
	}
	out := make([]rootio.Key, 0, len(set))
	for _, k := range set {
		out = append(out, k)
	}
	return out
}