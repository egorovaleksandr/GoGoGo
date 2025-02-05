//go:build !solution

package otp

import (
	"io"
)

type CipherReader struct {
	r_    io.Reader
	prng_ io.Reader
}

func (cr CipherReader) Read(p []byte) (n int, err error) {
	n, err = cr.r_.Read(p)
	rn := make([]byte, n)
	m, _ := cr.prng_.Read(rn)
	if m < n {
		n = m
	}
	for i := 0; i < n; i++ {
		p[i] = p[i] ^ rn[i]
	}
	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return CipherReader{r, prng}
}

type CipherWriter struct {
	w_    io.Writer
	prng_ io.Reader
}

func (wr CipherWriter) Write(p []byte) (n int, err error) {
	rn := make([]byte, len(p))
	m, _ := wr.prng_.Read(rn)
	rp := make([]byte, m)
	for i := 0; i < m; i++ {
		rp[i] = p[i] ^ rn[i]
	}
	n, err = wr.w_.Write(rp)
	return
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return CipherWriter{w, prng}
}
