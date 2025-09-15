package foo

// TODO what is go.work file?
// i wonder how we can set up a go project
type Foo struct{ Bar string }

// Read returns Bar
func (f *Foo) Read() string {
	return f.Bar
}
