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

type Stepping struct {
	current int
	steps   Steps
}

func NewStepping(str string) *Stepping {
	return &Stepping{
		0,
		GetSteps(str),
	}
}

// Do we still need to go deeper
func (s *Stepping) More() bool {
	return s.current < len(s.steps)
}

func (s *Stepping) StepIn() {
	s.current++
}

func (s *Stepping) Current() *Step {
	return &s.steps[s.current]
}

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

	addr := NewStepping(os.Args[1])
	dec := json.NewDecoder(input)
	navigate(dec, addr)
}

func isOk(err error) bool {
	if err == io.EOF {
		return false
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return false
	}
	return true
}

func navigate(d *json.Decoder, s *Stepping) {
	if !s.More() {
		printObject(d)
		return
	}

	t, err := d.Token()
	if !isOk(err) {
		// TODO mark error in stepping
		return
	}

	switch i := t.(type) {
	case json.Delim:
		if i == '[' {
			if s.Current().kind != INDEX {
				// ERROR
				return
			}
			for idx := 0; d.More(); idx++ {
				if idx == s.Current().index {
					s.StepIn()
					navigate(d, s)
					return
				}
			}
			// ERROR
			return
		}
		if i == '{' {
			if s.Current().kind != KEY {
				// ERROR
				return
			}
			for d.More() {
				if navigateKeyValue(d, s) {
					s.StepIn()
					navigate(d, s)
					return
				}
			}
			// ERROR
			return
		}
	default:
	}
}

func navigateKeyValue(d *json.Decoder, s *Stepping) bool {
	// key
	t, err := d.Token()
	if !isOk(err) {
		return false
	}
	if fmt.Sprint(t) == s.Current().name {
		return true
	}
	// value
	skipObject(d)
	return false
}

func skipObject(d *json.Decoder) {
	t, err := d.Token()
	if !isOk(err) {
		return
	}

	switch i := t.(type) {
	case json.Delim:
		if i == '[' {
			for d.More() {
				skipObject(d)
			}
			skipObject(d)
		} else if i == '{' {
			for d.More() {
				skipObject(d)
				skipObject(d)
			}
			skipObject(d)
		}
	default:
	}
}

func printObject(d *json.Decoder) {
	t, err := d.Token()
	if !isOk(err) {
		return
	}

	switch i := t.(type) {
	case json.Delim:
		fmt.Print(i)
		if i == '[' {
			if d.More() {
				printObject(d)
			}
			for d.More() {
				fmt.Print(",")
				printObject(d)
			}
			printObject(d)
		} else if i == '{' {
			if d.More() {
				printKeyValue(d)
			}
			for d.More() {
				fmt.Print(", ")
				printKeyValue(d)
			}
			printObject(d)
		}
	case string:
		fmt.Printf("\"%s\"", i)
	default:
		fmt.Print(t)
	}
}

func printKeyValue(d *json.Decoder) {
	printObject(d)
	fmt.Print(": ")
	printObject(d)
}
