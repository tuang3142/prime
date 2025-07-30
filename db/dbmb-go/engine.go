package main

import (
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

type FileScan struct {
	childAccessor
	reader Reader
	eof    bool
}

func (f *FileScan) Next() Row {
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
