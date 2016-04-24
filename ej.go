package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	if len(str) == 0 {
		return s
	}
	// it's what's left, although it's on the rhs ;)
	left := str
	var dot, bracket int
	for {
		dot = strings.Index(left, ".")
		bracket = strings.Index(left, "[")
		if dot == -1 && bracket == -1 {
			break
		}
		if dot != -1 && dot < bracket || bracket == -1 {
			if dot != 0 {
				s = append(s, Key(left[0:dot]))
			}
			left = left[dot+1:]
		} else {
			if bracket != 0 {
				s = append(s, Key(left[0:bracket]))
			}
			closed := strings.Index(left, "]")
			num, _ := strconv.Atoi(left[bracket+1 : closed])
			s = append(s, Index(num))
			left = left[closed+1:]
		}
	}

	if len(left) > 0 {
		s = append(s, Key(left))
	}

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
