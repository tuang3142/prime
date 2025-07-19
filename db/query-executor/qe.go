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

type Filter_F func(Row) bool

func And(fs ...Filter_F) Filter_F {
	return func(r Row) bool {
		ret := true
		for _, f := range fs {
			ret = ret && f(r)
		}
		return ret
	}
}

func Or(fs ...Filter_F) Filter_F {
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
	filter Filter_F
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

// TODO: tbd
func jsonToNode(json string) []Node { return []Node{} }

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

func main() {
	db := []Row{
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
	match := func(regex string) Filter_F {
		return func(r Row) bool {
			id := r[0].(string)
			matched, err := regexp.Match(regex, []byte(id))
			if err != nil {
				panic(err)
			}
			return matched
		}
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
		&Selection{filter: Or(match("^hp"), match("^lotr"))},
		&MemoryScan{table: db},
	))
	if err := assert(got, want); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}

	count := func(rows []Row) Row { return Row{len(rows)} }
	want = []Row{
		{"hp", 3},
		{"lotr", 3},
	}
	got = run(q(
		&GroupBy{groupId: 4, agg: count},
		&Selection{filter: Or(match("^hp"), match("^lotr"))},
		&MemoryScan{table: db},
	))
	if err := assert(got, want); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}

	agg := func(rows []Row) Row { return Row{rows[0]} } // group rows by their id, then produce a copy of the row itself.
	// not very pretty but it proves the point
	want = []Row{
		[]any{"hp1", []any{"hp1", "Harry Potter and the Sorcerer's Stone", 7.6, true, "hp"}},
	}
	got = run(q(
		&GroupBy{groupId: 0, agg: agg},
		&Limit{limit: 1},
		&Selection{filter: Or(match("^hp"), match("^lotr"))},
		&MemoryScan{table: db},
	))
	if err := assert(got, want); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
