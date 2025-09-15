package queryexecutor

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

func TestMemoryScan(t *testing.T) {
	rows := []Row{
		{1, "Harry J. Potter", "mage|orphan|main"},
		{2, "Tom M. Riddle", "mage|orphan|also main, but evil"},
	}

	got := run(q(&MemoryScan{Rows: rows}))

	if diff := cmp.Diff(rows, got); diff != "" {
		t.Errorf("rows differ (-want +got):\n%s", diff)
	}
}

func TestSelection(t *testing.T) {
	rows := []Row{
		{1, "Harry J. Potter", "mage|orphan|main"},
		{2, "Tom M. Riddle", "mage|orphan|also main, but evil"},
		{3, "Hermione Granger", "mage|muggle-born|main"},
		{4, "Severus Snape", "mage|professor"},
	}
	substr := func(col int, needle string) FilterFunc {
		return func(r Row) bool {
			s, _ := r[col].(string)
			return strings.Contains(s, needle)
		}
	}

	tests := []struct {
		name   string
		filter FilterFunc
		want   []Row
	}{
		{
			name:   "name_contains_Potter",
			filter: substr(1, "Potter"),
			want:   []Row{{1, "Harry J. Potter", "mage|orphan|main"}},
		},
		{
			name:   "no_match",
			filter: substr(1, "Dobby"),
			want:   nil,
		},
		{
			name:   "AND_orphan_and_mage",
			filter: And(substr(2, "orphan"), substr(2, "mage")),
			want: []Row{
				{1, "Harry J. Potter", "mage|orphan|main"},
				{2, "Tom M. Riddle", "mage|orphan|also main, but evil"},
			},
		},
		{
			name:   "OR_Hermione_or_Snape",
			filter: Or(substr(1, "Hermione"), substr(1, "Snape")),
			want: []Row{
				{3, "Hermione Granger", "mage|muggle-born|main"},
				{4, "Severus Snape", "mage|professor"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := run(q(
				&Selection{Filter: tc.filter},
				&MemoryScan{Rows: rows},
			))

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("rows differ (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProjection(t *testing.T) {
	rows := []Row{
		{1, "Harry", "Potter", "Sly"},
		{2, "Hermione", "Granger", "Rav"},
		{3, "Ron", "Weasley", "Huf"},
	}
	mapper := func(r Row) Row { return Row{r[1], r[3]} }

	got := run(q(
		&Projection{Mapper: mapper},
		&MemoryScan{Rows: rows},
	))

	want := []Row{
		{"Harry", "Sly"},
		{"Hermione", "Rav"},
		{"Ron", "Huf"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("rows differ (-want +got):\n%s", diff)
	}
}
