// Package main contains a very crude implementation of the query executor
package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
)

type Row []any

type Node interface {
	Next() Row
	SetChild(Node)
	GetChild() Node
}

type ChildAccess struct {
	child Node
}

func (c *ChildAccess) GetChild() Node  { return c.child }

func (c *ChildAccess) SetChild(n Node) { c.child = n }

type MemoryScan struct {
	ChildAccess

	table []Row
	idx   int
}

func NewMemoryScan(table []Row) *MemoryScan {
	return &MemoryScan{table: table, idx: 0}
}

func (m *MemoryScan) Next() Row {
	if m.idx >= len(m.table) {
		return nil
	}
	r := m.table[m.idx]
	m.idx++
	return r
}

type Condition func(Row) bool

type Selection struct {
	ChildAccess

	cond Condition
}

func NewSelection(f Condition) *Selection {
	return &Selection{cond: f}
}

func (s *Selection) Next() Row {
	row := s.GetChild().Next()
	if row == nil {
		return nil
	}
	if s.cond(row) {
		return row
	}
	return []any{} // empty row which will be skipped by the parent (see Sort() and Run())
}

type Mapper func(Row) Row

type Projection struct {
	ChildAccess

	mp Mapper
}

func NewProjection(mp Mapper) *Projection {
	return &Projection{mp: mp}
}

func (p *Projection) Next() Row {
	row := p.GetChild().Next()
	if row == nil || len(row) == 0 {
		return row
	}
	return p.mp(row)
}

type Limit struct {
	ChildAccess

	cnt int
}

func (l *Limit) Next() Row {
	row := l.GetChild().Next()
	if row == nil || len(row) == 0 {
		return row
	}

	if l.cnt == 0 {
		return []any{}
	}
	l.cnt--
	return row
}

func NewLimit(l int) *Limit {
	return &Limit{cnt: l}
}

type Sort struct {
	ChildAccess

	table  []Row
	idx    int
	sortBy int // index of the "column" to be used as the sort key
	sorted bool
}

func NewSort(sortBy int) *Sort {
	return &Sort{sortBy: sortBy, sorted: false, idx: 0}
}

func (s *Sort) Next() Row {
	// if not sorted: get all the rows from child
	// sort them, then return row by row
	if !s.sorted {
		s.sort()
	}
	// similar to table scan
	if s.idx >= len(s.table) {
		return nil
	}
	row := s.table[s.idx]
	s.idx++
	return row
}

func (s *Sort) sort() {
	for {
		rec := s.GetChild().Next()
		if rec == nil {
			break
		}
		if len(rec) == 0 {
			continue
		}
		s.table = append(s.table, rec)
	}
	// perform sorting
	sort.Slice(s.table, func(i, j int) bool {
		r1, r2 := s.table[i], s.table[j]
		v1, _ := r1[s.sortBy].(string)
		v2, _ := r2[s.sortBy].(string)
		return v1 < v2
	})
	s.sorted = true
}

// Q procudes a linked-list from a series of node.
// Assume that nodes are not empty.
func Q(nodes ...Node) Node {
	var head Node
	head = nodes[0]
	for _, node := range nodes[1:] {
		head.SetChild(node)
		head = head.GetChild()
	}
	return nodes[0]
}

// Run executes the query by calling `next` on the (presumed) root.
func Run(q Node) []Row {
	var result []Row
	for rec := q.Next(); rec != nil; rec = q.Next() {
		if len(rec) != 0 { // skip unqualified, empty row
			result = append(result, rec)
		}
	}
	return result
}

func cmp(want, got []Row) bool {
	if len(want) != len(got) {
		log.Printf("Diff: length\n")
		return false
	}
	for i := range want {
		if !reflect.DeepEqual(want[i], got[i]) {
			log.Printf("Diff:\n-%v\n+%v\n", want[i], got[i])
			return false
		}
	}
	return true
}

// schema
// id string, name string, weight float64, native bool
var db []Row = []Row{
	{"ostric", "Ostrich", 104.0, false},
	{"amerob", "American Robin", 0.077, true},
	{"baleag", "Bald Eagle", 4.74, true},
	{"eursta", "European Starling", 0.082, true},
	{"barswa1", "Barn Swallow", 0.019, true},
	{"barswa2", "Barn Swallow", 0.019, true},
	{"barswa3", "Barn Swallow", 0.019, true},
	{"barswa5", "Barn Swallow", 11.019, true},
	{"barswa4", "Barn Swallow", 10.019, true},
}

var dbAdvance []Row = []Row{
	{"ostric", "Ostrich", 104.0, false},
	{"amerob", "American Robin", 0.077, true},
	{"baleag", "Bald Eagle", 4.74, true},
	{"eursta", "European Starling", 0.082, true},
	{"barswa5", "Barn Swallow", 5.019, true},
	{"barswa4", "Barn Swallow", 4.019, true},
	{"barswa1", "Barn Swallow", 1.019, true},
	{"barswa2", "Barn Swallow", 2.019, true},
	{"barswa3", "Barn Swallow", 3.019, true},
}

func TestBasic() {
	native := func(r Row) bool {
		b, _ := r[3].(bool)
		return b
	}
	onlyNameWeight := func(r Row) Row {
		return []any{r[0], r[2]}
	}
	want := []Row{
		{"amerob", 0.077},
		{"baleag", 4.74},
		{"eursta", 0.082},
		{"barswa1", 0.019},
		{"barswa2", 0.019},
		{"barswa3", 0.019},
		{"barswa5", 11.019},
		{"barswa4", 10.019},
	}

	got := Run(Q(
		NewProjection(onlyNameWeight),
		NewSelection(native),
		NewMemoryScan(db),
	))

	if !cmp(want, got) {
		fmt.Printf("failed!\n")
	} else {
		fmt.Printf("ok\n")
	}
}

var and = func(c1, c2 Condition) Condition {
	return func(row Row) bool {
		return c1(row) && c2(row)
	}
}
var or = func(c1, c2 Condition) Condition {
	return func(row Row) bool {
		return c1(row) || c2(row)
	}
}

func TestAdvance() {
	native := func(r Row) bool {
		b, _ := r[3].(bool)
		return b
	}
	heavy := func(r Row) bool {
		weight, _ := r[2].(float64)
		return weight >= 1
	}
	onlyNameWeight := func(r Row) Row {
		return []any{r[0], r[2]}
	}
	want := []Row{
		{"baleag", 4.74},
		{"barswa1", 1.019},
		{"barswa2", 2.019},
		{"barswa3", 3.019},
		{"barswa4", 4.019},
		{"barswa5", 5.019},
	}

	got := Run(Q(
		NewSort(0),
		NewLimit(10),
		NewProjection(onlyNameWeight),
		NewSelection(and(native, heavy)), // filter with combined condition
		NewMemoryScan(dbAdvance),
	))

	if !cmp(want, got) {
		fmt.Printf("failed!\n")
	} else {
		fmt.Printf("ok\n")
	}
}

func main() {
	TestBasic()
	TestAdvance()
}
