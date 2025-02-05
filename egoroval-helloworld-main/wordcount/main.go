//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	m := make(map[string]int)
	for _, file := range os.Args[1:] {
		data, err := os.ReadFile(file)
		check(err)
		strs := strings.Split(string(data), "\n")
		for _, s := range strs {
			m[s]++
		}
	}
	for s, n := range m {
		if n > 1 {
			fmt.Printf("%d	%s\n", n, s)
		}
	}
}
