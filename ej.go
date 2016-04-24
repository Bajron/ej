package main

import (
	"encoding/json"
	"fmt"
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
	return Step{INDEX, fmt.Sprintf("[%d]", i), i}
}

type stepError struct {
	info string
}

func (se stepError) Error() string {
	return se.info
}

type Steps []Step

type Stepping struct {
	current int
	steps   Steps
	err     error
}

func NewStepping(str string) *Stepping {
	return &Stepping{
		0,
		GetSteps(str),
		nil,
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

func (s *Stepping) SoFar() Steps {
	return s.steps[:s.current]
}

func (s Steps) String() string {
	r := ""
	if len(s) > 0 {
		r = s[0].name
	}
	for i := 1; i < len(s); i++ {
		if s[i].kind == KEY {
			r += "."
		}
		r += s[i].name
	}
	return r
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

	address := NewStepping(os.Args[1])
	decoder := json.NewDecoder(input)
	navigate(decoder, address)

	if address.err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", address.err)
		fmt.Fprintf(os.Stderr, "found so far: %s\n", address.SoFar())
		os.Exit(1)
	}
}

func nextToken(d *json.Decoder, s *Stepping) (json.Token, bool) {
	if s.err != nil {
		return json.Token(nil), false
	}

	t, err := d.Token()

	if err != nil {
		s.err = err
		return t, false
	}
	return t, true
}

func navigate(d *json.Decoder, s *Stepping) {
	if !s.More() {
		printObject(d, s)
		return
	}

	t, ok := nextToken(d, s)
	if !ok {
		return
	}

	switch i := t.(type) {
	case json.Delim:
		if i == '[' {
			if s.Current().kind != INDEX {
				s.err = stepError{"Met array, when not expected"}
				return
			}
			for idx := 0; d.More(); idx++ {
				if idx == s.Current().index {
					s.StepIn()
					navigate(d, s)
					return
				}
				skipObject(d, s)
			}
			s.err = stepError{"Requested index not found"}
			return
		}
		if i == '{' {
			if s.Current().kind != KEY {
				s.err = stepError{"Met object, when not expected"}
				return
			}
			for d.More() {
				if navigateKeyValue(d, s) {
					s.StepIn()
					navigate(d, s)
					return
				}
			}
			s.err = stepError{"Requested key not found"}
			return
		}
	default:
		s.err = stepError{"Cannot look deeper"}
	}
}

func navigateKeyValue(d *json.Decoder, s *Stepping) bool {
	// key
	t, ok := nextToken(d, s)
	if !ok {
		return false
	}
	if fmt.Sprint(t) == s.Current().name {
		return true
	}
	// value
	skipObject(d, s)
	return false
}

func skipObject(d *json.Decoder, s *Stepping) {
	t, ok := nextToken(d, s)
	if !ok {
		return
	}

	switch i := t.(type) {
	case json.Delim:
		if i == '[' {
			for d.More() {
				skipObject(d, s)
			}
			skipObject(d, s)
		} else if i == '{' {
			for d.More() {
				skipObject(d, s)
				skipObject(d, s)
			}
			skipObject(d, s)
		}
	default:
	}
}

func printObject(d *json.Decoder, s *Stepping) {
	t, ok := nextToken(d, s)
	if !ok {
		return
	}

	switch i := t.(type) {
	case json.Delim:
		fmt.Print(i)
		if i == '[' {
			if d.More() {
				printObject(d, s)
			}
			for d.More() {
				fmt.Print(",")
				printObject(d, s)
			}
			printObject(d, s)
		} else if i == '{' {
			if d.More() {
				printKeyValue(d, s)
			}
			for d.More() {
				fmt.Print(",")
				printKeyValue(d, s)
			}
			printObject(d, s)
		}
	case string:
		fmt.Printf("\"%s\"", i)
	default:
		fmt.Print(t)
	}
}

func printKeyValue(d *json.Decoder, s *Stepping) {
	printObject(d, s)
	fmt.Print(":")
	printObject(d, s)
}
