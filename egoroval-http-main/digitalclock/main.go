//go:build !solution

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var f string = "Content-Type"
var s string = "image/png"
var t string = "15:04:05"
var m map[int]string = map[int]string{0: Zero, 1: One, 2: Two, 3: Three, 4: Four, 5: Five, 6: Six, 7: Seven, 8: Eight, 9: Nine}

func handler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse("http://" + r.Host + r.URL.String())
	kq := 1
	tq := ""
	var x, y int
	if err != nil {
		panic(err)
	}
	q := u.Query()
	tq = time.Now().Format(t)
	if len(q["time"]) != 0 && q["time"][0] != "" {
		tq = q["time"][0]
	}
	if len(tq) != 8 {
		http.Error(w, "", 400)
	}
	_, err = time.Parse(t, tq)
	if err != nil {
		http.Error(w, "", 400)
	}
	if len(q["k"]) != 0 && q["k"][0] != "" {
		kq, err = strconv.Atoi(q["k"][0])
	}
	if err != nil || kq < 1 || kq > 30 {
		http.Error(w, "", 400)
	}
	img := image.NewRGBA(image.Rect(0, 0, 56*kq, 12*kq))
	for _, value := range tq {
		if value == ':' {
			img, x, y = mC(x, y, img, kq)
			continue
		}
		img, x, y = mD(x, y, img, value, kq)
	}
	err = png.Encode(w, img)
	if err != nil {
		return
	}
	w.Header().Set(f, s)
	w.WriteHeader(200)
}

func mD(xs int, ys int, img *image.RGBA, digit int32, k int) (*image.RGBA, int, int) {
	x := xs
	y := ys
	digitDraw := m[(int(digit - '0'))]
	for i, sampleItem := range digitDraw {
		if digitDraw[i] == 10 {
			continue
		}
		for y_ := y * k; y_ < (y+1)*k; y_++ {
			for x_ := x * k; x_ < (x+1)*k; x_++ {
				if string(sampleItem) == "1" {
					img.Set(x_, y_, Cyan)
					continue
				}
				img.Set(x_, y_, color.White)
			}
		}
		if x == xs+7 {
			x = xs
			y++
			continue
		}
		x++
	}
	return img, xs + 8, 0
}

func mC(xs int, ys int, img *image.RGBA, k int) (*image.RGBA, int, int) {
	x := xs
	y := ys
	for i, sampleItem := range Colon {
		if Colon[i] == 10 {
			continue
		}
		for y_ := y * k; y_ < (y+1)*k; y_++ {
			for x_ := x * k; x_ < (x+1)*k; x_++ {
				if string(sampleItem) == "1" {
					img.Set(x_, y_, Cyan)
					continue
				}
				img.Set(x_, y_, color.White)
			}
		}
		if x == 3+xs {
			x = xs
			y++
			continue
		}
		x++
	}
	return img, xs + 4, 0
}

func getR(num int) int {
	portPtr := flag.Int("port", num, "port string")
	flag.Parse()
	return *portPtr
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:"+strconv.Itoa(getR(rand.Intn(10000))), nil)
	if err != nil {
		panic(err)
	}
}
