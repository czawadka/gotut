package main

import (
"io"
"os"
"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == nil {
		for i := 0; i < n; i++ {
			var start byte
			if 'a' <= p[i] && p[i] <= 'z' {
				start = 'a'
			} else if 'A' <= p[i] && p[i] <= 'Z' {
				start = 'A'
			}
			if start > 0 {
				p[i] = start+((p[i]-start)+13)%26
			}
		}
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
