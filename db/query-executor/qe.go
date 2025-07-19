// Binary main runs an implementatioin of a DB query executor.
package main

import (
	"fmt"
	"reflect"
	"regexp"
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

type MemoryScan struct {
	childAccessor
	table []Row
	i     int
}

func (m *MemoryScan) Next() Row {
	if m.i >= len(m.table) {
		return nil
	}
	row := m.table[m.i]
	m.i++
	return row
}

type Selection struct {
	childAccessor
	filter func(Row) bool
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

type Limit struct {
	childAccessor
	limit int
}

func (l *Limit) Next() Row {
	if l.limit <= 0 {
		return nil
	}
	l.limit -= 1
	return l.GetChild().Next()
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

func assert(got, want []Row) error {
	if len(got) != len(want) {
		return fmt.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			return fmt.Errorf("mismatch at index %d:\n-got  %v (type %T)\n+want %v (type %T)", i, got[i], got[i], want[i], want[i])
		}
	}
	return nil
}

// todo: make this a smaller version, an actual poc of the split db - the dal
// how do project irl use interface
// stretch goal: find a better design for the binary filtering
func main() {
	db := []Row{
		{"hp1", "Harry Potter and the Sorcerer's Stone", 7.6, true},
		{"hp2", "Harry Potter and the Chamber of Secrets", 7.4, false},
		{"hp3", "Harry Potter and the Prisoner of Azkaban", 7.9, true},
		{"inception", "Inception", 8.8, true},
		{"matrix", "The Matrix", 8.7, true},
		{"godfather", "The Godfather", 9.2, false},
		{"pulp", "Pulp Fiction", 8.9, true},
		{"lotr1", "The Lord of the Rings: The Fellowship of the Ring", 8.8, false},
		{"lotr2", "The Lord of the Rings: The Two Towers", 8.8, true},
		{"lotr3", "The Lord of the Rings: The Return of the King", 9.0, true},
	}
	filter := func(row Row) bool {
		r, _ := row[0].(string)
		m1, err := regexp.Match("^hp", []byte(r))
		if err != nil {
			panic(err)
		}
		m2, err := regexp.Match("^lotr", []byte(r))
		if err != nil {
			panic(err)
		}
		return m1 || m2
	}
	mapper := func(row Row) Row { return Row{row[0], row[1], row[2]} }
	key := func(row Row) float64 {
		rating, ok := row[2].(float64)
		if !ok {
			panic("can't convert rating to float64")
		}
		return rating
	}

	want := []Row{
		{"lotr1", "The Lord of the Rings: The Fellowship of the Ring", 8.8},
		{"hp3", "Harry Potter and the Prisoner of Azkaban", 7.9},
		{"hp1", "Harry Potter and the Sorcerer's Stone", 7.6},
		{"hp2", "Harry Potter and the Chamber of Secrets", 7.4},
	}
	got := run(q(
		&Projection{mapper: mapper},
		&Sort{key: key, desc: true},
		&Limit{limit: 4},
		&Selection{filter: filter},
		&MemoryScan{table: db},
	))
	if err := assert(got, want); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
