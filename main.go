package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/go-logfmt/logfmt"
)

func main() {
	var verbose bool
	verboseFlag(&verbose)
	flag.Parse()

	dec := json.NewDecoder(os.Stdin)
	enc := logfmt.NewEncoder(os.Stdout)

	for {
		m := make(map[string]interface{})
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("failure to parse json: %s", err)

			dec = json.NewDecoder(io.MultiReader(advanceLine(dec.Buffered()), os.Stdin))
			continue
		}

		for k, v := range m {
			if err := enc.EncodeKeyval(k, v); err != nil {
				if verbose {
					log.Printf("failure to encode to logfmt: %s", err)
				}
			}
		}

		if err := enc.EndRecord(); err != nil {
			log.Fatalf("failure to write a record terminator: %s", err)
		}
	}
}

func verboseFlag(v *bool) {
	usage := "Complain to stderr on errors"
	flag.BoolVar(v, "v", false, usage)
	flag.BoolVar(v, "verbose", false, usage)
}

func advanceLine(rdr io.Reader) io.Reader {
	r := bufio.NewReader(rdr)
	r.ReadBytes('\n')
	return r
}
