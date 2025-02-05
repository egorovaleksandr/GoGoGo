//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	for _, url := range os.Args[1:] {
		res, err := http.Get(url)
		check(err)
		body, err := io.ReadAll(res.Body)
		check(err)
		res.Body.Close()
		fmt.Printf("%s", body)
	}
}
