package main

import (
	"net/http"
	"fmt"
)

type String string

func (handler String) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprint(writer, "Hello!")
}

type Struct struct {
	Greeting string
	Punct string
	Who string
}
func (handler Struct) ServeHTTP(out http.ResponseWriter, req *http.Request) {
	a := req.FormValue("a")
	fmt.Fprintf(out, "%s%s%s - a=%s", handler.Greeting, handler.Punct, handler.Who, a)
}

func main() {
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}
