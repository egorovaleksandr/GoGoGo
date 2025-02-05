//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Evaluator struct {
	wordList map[string][]string
	stack    []int
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{make(map[string][]string), make([]int, 0)}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	isDef := false
	commands := strings.Split(strings.ToLower(row), " ")
	forDef := 0
	var tmp []string
	var err error
	for i, c := range commands {
		if c == ":" {
			isDef = true
			forDef = i + 1
			if _, ok := e.wordList[commands[forDef]]; ok {
				tmp = e.wordList[commands[forDef]]
				e.wordList[commands[forDef]] = e.wordList[commands[forDef]][:0]
			}
			continue
		}
		if c == ";" {
			isDef = false
			continue
		}
		if isDef && i == forDef {
			if _, er := strconv.Atoi(commands[forDef]); er == nil {
				return e.stack, errors.New("redefine number")
			}
			continue
		}
		if isDef {
			if _, ok := e.wordList[c]; ok {
				if c == commands[forDef] {
					e.wordList[commands[forDef]] = append(e.wordList[commands[forDef]], tmp...)
				} else {
					e.wordList[commands[forDef]] = append(e.wordList[commands[forDef]], e.wordList[c]...)
				}
				continue
			}
			e.wordList[commands[forDef]] = append(e.wordList[commands[forDef]], c)
			continue
		}
		if wl, ok := e.wordList[c]; ok {
			for _, dc := range wl {
				err = e.DoCommand(dc)
			}
			continue
		}
		err = e.DoCommand(c)
	}
	return e.stack, err
}

func (e *Evaluator) Pop() (int, error) {
	n := len(e.stack)
	if n < 1 {
		return 0, errors.New("few arguments")
	}
	res := e.stack[n-1]
	e.stack = e.stack[:n-1]
	return res, nil
}

func (e *Evaluator) DoCommand(command string) error {
	var last, prev int
	var err error
	switch command {
	case "drop":
		_, err = e.Pop()
		if err != nil {
			return err
		}
	case "dup":
		last, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, last, last)
	case "+":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, prev+last)
	case "-":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, prev-last)
	case "*":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, prev*last)
	case "/":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		if last == 0 {
			return errors.New("division by zero")
		}
		e.stack = append(e.stack, prev/last)
	case "over":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, prev, last, prev)
	case "swap":
		last, _ = e.Pop()
		prev, err = e.Pop()
		if err != nil {
			return err
		}
		e.stack = append(e.stack, last, prev)
	default:
		num, err := strconv.Atoi(command)
		if err != nil {
			return errors.New("undefined command")
		}
		e.stack = append(e.stack, num)
	}
	return nil
}
