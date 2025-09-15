package foo

import "testing"

func newFoo(t *testing.T, bar string) *Foo {
	t.Helper()
	return &Foo{Bar: bar}
}

func TestFoo_Read(t *testing.T) {
	foo := newFoo(t, "bar")

	got := foo.Read()

	if want := "bar"; got != want {
		t.Errorf("Foo.Read() = %v, but want %v", got, want)
	}
}
