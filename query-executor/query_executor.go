package main

import (
	"fmt"
	"log"
	"reflect"
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

// can interface have pointer? inuitively, its doesnt make sense, but we could have one
func (c *ChildAccess) GetChild() Node  { return c.child }
func (c *ChildAccess) SetChild(n Node) { c.child = n }

type MemoryScan struct {
	ChildAccess

	table []Row // each node will store its owns state
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
	return []any{} // this works, for now, until we implement limit and such
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
	{"barswa", "Barn Swallow", 0.019, true},
}

func TestBasic() {
	native := func(r Row) bool {
		b, _ := r[3].(bool)
		return b
	}
	onlyName := func(r Row) Row {
		return []any{r[0], r[2]}
	}
	want := []Row{
		{"amerob", 0.077},
		{"baleag", 4.74},
		{"eursta", 0.082},
		{"barswa", 0.019},
	}

	got := Run(Q(
		NewProjection(onlyName),
		NewSelection(native),
		NewMemoryScan(db),
	))

	if !cmp(want, got) {
		fmt.Printf("mis-match!\n")
	} else {
		fmt.Printf("ok\n")
	}
}

func main() {
	TestBasic()
}
