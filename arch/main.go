package main

import (
	"varint"
)

func main() {
	f := &varint.Foo{Bar: "bar"}
	f.Hello()
}
