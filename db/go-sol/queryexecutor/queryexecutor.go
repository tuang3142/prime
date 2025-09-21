package queryexecutor

import (
	"os"
	"sort"
)

type Row []any

type Node interface {
	SetChild(Node)
	GetChild() Node
	Next() Row
	Close()
}

type childAccessor struct{ child Node }

func (c *childAccessor) GetChild() Node  { return c.child }
func (c *childAccessor) SetChild(n Node) { c.child = n }
func (c *childAccessor) Close() {
	child := c.GetChild()
	if child != nil {
		child.Close()
	}
}

type MemoryScan struct {
	childAccessor
	Rows []Row
	i    int
}

func (m *MemoryScan) Next() Row {
	if m.i >= len(m.Rows) {
		return nil
	}
	row := m.Rows[m.i]
	m.i++
	return row
}

type FileScan struct {
	childAccessor
	FileName string
	Closed   bool
	file     *os.File
}

func (f *FileScan) Next() {
	if f.file == nil {
		var err error
		f.file, err = os.Open(f.FileName)
		if err != nil {
			panic(err)
		}
	}
}

func (f *FileScan) Close() {
	f.Closed = true
}

type Selection struct {
	childAccessor
	Filter FilterFunc
}

func (s *Selection) Next() Row {
	for {
		row := s.GetChild().Next()
		if row == nil {
			return nil
		}
		if s.Filter(row) {
			return row
		}
	}
}

type Projection struct {
	childAccessor
	Mapper MapperFunc
}

func (p *Projection) Next() Row {
	row := p.GetChild().Next()
	if row != nil {
		return p.Mapper(row)
	}
	return nil
}

type Limit struct {
	childAccessor
	L int
}

func (l *Limit) Next() Row {
	if l.L <= 0 {
		return nil
	}
	l.L -= 1
	return l.GetChild().Next()
}

type Sort struct {
	childAccessor
	Key    SortFunc
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
			x, y := s.Key(s.table[i]), s.Key(s.table[j])
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

// GroupBy categorizes rows by a column which will be used as key in map (ex: map[string][]Row),
// then performs AggFunc on those rows (ex: sum their value, get max value, count).
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

func Run(head Node) []Row { return run(head) }

func run(head Node) []Row {
	defer head.Close() // is this a good pattern?

	var result []Row
	for row := head.Next(); row != nil; row = head.Next() {
		result = append(result, row)
	}
	return result
}

func Q(nodes ...Node) Node { return q(nodes...) }

func q(nodes ...Node) Node {
	for i, node := range nodes[:len(nodes)-1] {
		next := nodes[i+1]
		node.SetChild(next)
	}
	return nodes[0]
}
