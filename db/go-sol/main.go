// Binary main runs an implementatioin of a DB query executor.
package main

import (
	"flag"
	"io"
	"sort"
)

type Row []any

type Node interface {
	SetChild(Node)
	GetChild() Node
	Next() Row
}

type childAccessor struct{ child Node }

func (c *childAccessor) GetChild() Node  { return c.child }
func (c *childAccessor) SetChild(n Node) { c.child = n }

type Reader interface {
	Read() ([]string, error)
}

// stretch exercise: check for memeory leak?
// file might not be closed at all
type CSVScan struct {
	childAccessor
	reader Reader
	eof    bool
}

func (f *CSVScan) Next() Row {
	if f.eof {
		return nil
	}

	record, err := f.reader.Read()
	if err == io.EOF {
		f.eof = true
		return nil
	}
	if err != nil {
		panic(err)
	}

	row := make(Row, len(record))
	for i, val := range record {
		row[i] = val
	}
	return row
}

// type MemoryScan struct {
// 	childAccessor
// 	table []Row
// 	i     int
// }

// func (m *MemoryScan) Next() Row {
// 	if m.i >= len(m.table) {
// 		return nil
// 	}
// 	row := m.table[m.i]
// 	m.i++
// 	return row
// }

type FilterFunc func(Row) bool

func And(fs ...FilterFunc) FilterFunc {
	return func(r Row) bool {
		ret := true
		for _, f := range fs {
			ret = ret && f(r)
		}
		return ret
	}
}

func Or(fs ...FilterFunc) FilterFunc {
	return func(r Row) bool {
		ret := false
		for _, f := range fs {
			ret = ret || f(r)
		}
		return ret
	}
}

type Selection struct {
	childAccessor
	filter FilterFunc
}

func (s *Selection) Next() Row {
	for {
		row := s.GetChild().Next()
		if row == nil {
			return nil
		}
		if s.filter(row) {
			return row
		}
	}
}

type Projection struct {
	childAccessor
	mapper func(Row) Row
}

func (p *Projection) Next() Row {
	row := p.GetChild().Next()
	if row != nil {
		return p.mapper(row)
	}
	return nil
}

type Sort struct {
	childAccessor
	key    func(Row) float64 // need better design
	sorted bool
	desc   bool
	table  []Row
	i      int
}

func (s *Sort) Next() Row {
	if !s.sorted {
		for row := s.GetChild().Next(); row != nil; row = s.GetChild().Next() {
			s.table = append(s.table, row)
		}
		sort.Slice(s.table, func(i, j int) bool {
			x, y := s.key(s.table[i]), s.key(s.table[j])
			if s.desc {
				return x > y
			}
			return x < y
		})
		s.sorted = true
	}
	if s.i >= len(s.table) {
		return nil
	}
	row := s.table[s.i]
	s.i++
	return row
}

type AggFunc func([]Row) Row

type GroupBy struct {
	childAccessor
	agg     AggFunc
	groupId int
	grouped bool
	table   []Row
	i       int
}

func (g *GroupBy) Next() Row {
	if !g.grouped {
		groupMap := make(map[any][]Row)
		for r := g.GetChild().Next(); r != nil; r = g.GetChild().Next() {
			key := r[g.groupId]
			groupMap[key] = append(groupMap[key], r)
		}
		for key, r := range groupMap {
			aggRow := g.agg(r)
			g.table = append(g.table, append(Row{key}, aggRow...))
		}
		g.grouped = true
	}
	if g.i >= len(g.table) {
		return nil
	}
	r := g.table[g.i]
	g.i++
	return r
}
func run(head Node) []Row {
	var result []Row
	for row := head.Next(); row != nil; row = head.Next() {
		result = append(result, row)
	}
	return result
}

func q(nodes ...Node) Node {
	for i, node := range nodes[:len(nodes)-1] {
		next := nodes[i+1]
		node.SetChild(next)
	}
	return nodes[0]
}

var db []Row = []Row{
	{"hp1", "Harry Potter and the Sorcerer's Stone", 7.6, true, "hp"},
	{"hp2", "Harry Potter and the Chamber of Secrets", 7.4, false, "hp"},
	{"hp3", "Harry Potter and the Prisoner of Azkaban", 7.9, true, "hp"},
	{"inception", "Inception", 8.8, true, "sci"},
	{"matrix", "The Matrix", 8.7, true, "sci"},
	{"godfather", "The Godfather", 9.2, false, "dra"},
	{"pulp", "Pulp Fiction", 8.9, true, "dra"},
	{"lotr1", "The Lord of the Rings: The Fellowship of the Ring", 8.8, false, "lotr"},
	{"lotr2", "The Lord of the Rings: The Two Towers", 8.8, true, "lotr"},
	{"lotr3", "The Lord of the Rings: The Return of the King", 9.0, true, "lotr"},
}

func main() {
	t := flag.Bool("t", false, "Run test")
	flag.Parse()
	if *t {
		runTest()
		runTestReadFromCSV()
		return
	}

	// Running DB as a server
	// http.HandleFunc("/query", handleQuery)
	// fmt.Println("Read-only DB server running at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
