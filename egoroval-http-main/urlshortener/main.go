//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"math/rand"
	"net/http"
	"strconv"
)

var m map[string]string = make(map[string]string)
var f string = "/shorten"
var s string = "/go/"

type PURL struct {
	Key string
	URL string
}

type URL struct {
	URL string
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "", 400)
	} else {
		err = r.Body.Close()
		if err != nil {
			panic(err)
		}
		if m[url.URL] != "" {
			marshaled, _ := json.Marshal(PURL{m[url.URL], url.URL})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(marshaled)
			if err != nil {
				panic(err)
			}
		} else {
			keyString := strconv.Itoa(rand.Intn(1_000))
			m[keyString] = url.URL
			m[url.URL] = keyString
			marshaled, err := json.Marshal(PURL{keyString, url.URL})
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(marshaled)
			if err != nil {
				panic(err)
			}
		}
	}
}

func goHandler(w http.ResponseWriter, r *http.Request) {
	if m[r.URL.String()[4:]] == "" {
		http.Error(w, "", http.StatusNotFound)
	} else {
		http.Redirect(w, r, m[r.URL.String()[4:]], http.StatusFound)
	}
}

func getR(num int) int {
	portPtr := flag.Int("port", num, "port string")
	flag.Parse()
	return *portPtr
}

func main() {
	http.HandleFunc(f, shortenHandler)
	http.HandleFunc(s, goHandler)
	err := http.ListenAndServe("localhost:"+strconv.Itoa(getR(rand.Intn(10_000))), nil)
	if err != nil {
		panic(err)
	}
}
