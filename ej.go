package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	KEY = iota
	INDEX
)

type Step struct {
	kind  int
	name  string
	index int
}

func Key(name string) Step {
	return Step{KEY, name, -1}
}
func Index(i int) Step {
	return Step{INDEX, "", i}
}

type Steps []Step

func GetSteps(str string) Steps {
	s := make(Steps, 0, 16)
	return s
}

func main() {
	var input *os.File
	var err error

	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, `Usage:
%s <pattern> [file]`, os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		input, err = os.Open(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open provided file (%s) -- %s", os.Args[2], err)
			os.Exit(2)
		}
	} else {
		input = os.Stdin
	}

	dec := json.NewDecoder(input)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
		}
		fmt.Printf("%T: %v", t, t)
		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}
}
