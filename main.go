package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
)

const AppName = "numstat"

type Options struct {
	ShowAverage bool `short:"a" long:"avg"     description:"Show average"`
	ShowSum     bool `short:"s" long:"sum"     description:"Show sum"`
	ShowMaximum bool `short:"x" long:"max"     description:"Show maximum"`
	ShowMinimum bool `short:"n" long:"min"     description:"Show minimum"`
	ShowAsJson  bool `short:"j" long:"json"    description:"Show as JSON"`
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

func (o *Options) isAllDisabled() bool {
	return !o.ShowAverage && !o.ShowSum && !o.ShowMaximum && !o.ShowMinimum
}

func (o *Options) enableAll() {
	o.ShowAverage = true
	o.ShowSum = true
	o.ShowMaximum = true
	o.ShowMinimum = true

}

var opts Options

var maximum float64
var maximumFirst = true

var minimum float64
var minimumFirst = true

var sum float64
var count float64

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = AppName
	parser.Usage = "[OPTIONS] FILES..."

	args, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	r, err := argf.From(args)
	if err != nil {
		panic(err)
	}

	if opts.isAllDisabled() {
		opts.enableAll()
	}

	err = numstat(r, os.Stdout, os.Stderr, opts)
	if err != nil {
		panic(err)
	}
}

func numstat(r io.Reader, stdout io.Writer, stderr io.Writer, opts Options) error {
	if opts.ShowVersion {
		io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return nil
	}

	reader := bufio.NewReader(r)
	var line []byte
	var err error
	for {
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		f, convErr := strconv.ParseFloat(string(line), 64)
		if convErr != nil {
			fmt.Fprintf(stderr, "number conversion error: %s\n", convErr)
			continue
		}

		if opts.ShowMaximum {
			feedMaximum(f)
		}
		if opts.ShowMinimum {
			feedMinimum(f)
		}
		if opts.ShowAverage || opts.ShowSum {
			sum += f
		}
		if opts.ShowAverage {
			count += 1
		}
	}

	if opts.ShowAsJson {
		resultJson := make(map[string]float64)

		if opts.ShowMaximum {
			resultJson["max"] = maximum
		}
		if opts.ShowMinimum {
			resultJson["min"] = minimum
		}
		if opts.ShowSum {
			resultJson["sum"] = sum
		}
		if opts.ShowAverage {
			resultJson["avg"] = sum / count
		}

		encoder := json.NewEncoder(stdout)
		encodeErr := encoder.Encode(resultJson)
		if encodeErr != nil {
			panic(encodeErr)
		}
	} else {
		if opts.ShowMaximum {
			fmt.Fprint(stdout, "Max: ", maximum, "\n")
		}
		if opts.ShowMinimum {
			fmt.Fprint(stdout, "Min: ", minimum, "\n")
		}
		if opts.ShowSum {
			fmt.Fprint(stdout, "Sum: ", sum, "\n")
		}
		if opts.ShowAverage {
			fmt.Fprint(stdout, "Avg: ", sum/count, "\n")
		}
	}

	return nil
}

func feedMaximum(f float64) {
	if maximumFirst {
		maximum = f
		maximumFirst = false
	} else {
		if f > maximum {
			maximum = f
		}
	}
}

func feedMinimum(f float64) {
	if minimumFirst {
		minimum = f
		minimumFirst = false
	} else {
		if f < minimum {
			minimum = f
		}
	}
}
