//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func QueryURL(url string, chanl chan int, num int) {
	start := time.Now()
	_, err := http.Get(url)
	if err == nil {
		fmt.Println("ok")
	}
	fmt.Println(time.Since(start), url)
	chanl <- num
}

func main() {
	chanl := make(chan int)
	start := time.Now()
	for i, url := range os.Args[1:] {
		go QueryURL(url, chanl, i)
	}
	for range os.Args[1:] {
		msg := <-chanl
		fmt.Println(msg)
	}
	fmt.Println(time.Since(start), "end")
}
