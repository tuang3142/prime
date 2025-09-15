package main

import (
	"pkgtest/foo"
)

func main() {
	f := foo.Foo{Bar: "bar"}
	print(f.Read())
}
