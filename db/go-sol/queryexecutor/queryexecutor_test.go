package queryexecutor

import (
	// "sort"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

func TestLimit(t *testing.T) {
	rows := []Row{
		{1, "Harry", "Potter", "Sly"},
		{2, "Hermione", "Granger", "Rav"},
		{3, "Ron", "Weasley", "Huf"},
	}
	mapper := func(r Row) Row { return Row{r[1], r[3]} }

	got := run(q(
		&Limit{L: 1},
		&Projection{Mapper: mapper},
		&MemoryScan{Rows: rows},
	))

	want := []Row{
		{"Harry", "Sly"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("rows differ (-want +got):\n%s", diff)
	}
}

func TestSort(t *testing.T) {
	rows := []Row{
		{4, "Draco", "Malfoy", "Sly"},
		{1, "Harry", "Potter", "Sly"},
		{2, "Hermione", "Granger", "Rav"},
		{3, "Ron", "Weasley", "Huf"},
	}
	key := func(r Row) int {
		id, ok := r[0].(int)
		if !ok {
			t.Fatalf("can't convert %v to int", r[0])
		}
		return id
	}

	got := run(q(
		&Limit{L: 2},
		&Sort{Key: key, desc: true},
		&MemoryScan{Rows: rows},
	))

	want := []Row{
		{4, "Draco", "Malfoy", "Sly"},
		{3, "Ron", "Weasley", "Huf"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("rows differ (-want +got):\n%s", diff)
	}

}

func TestGroupBy(t *testing.T) {
	rows := []Row{
		{1, "Harry", "Gryffindor", 10},
		{2, "Hermione", "Gryffindor", 12},
		{3, "Draco", "Slytherin", 8},
		{4, "Pansy", "Slytherin", 6},
	}
	sumPoints := func(rs []Row) Row {
		sum := 0
		for _, r := range rs {
			val, ok := r[3].(int)
			if !ok {
				t.Fatalf("can't convert %v to int", r[3])
			}
			sum += val
		}
		return Row{sum}
	}

	got := run(q(
		&GroupBy{groupId: 2, agg: sumPoints},
		&MemoryScan{Rows: rows},
	))
	sort.Slice(got, func(i, j int) bool {
		return got[i][0].(string) < got[j][0].(string)
	})

	want := []Row{
		{"Gryffindor", 22}, // 10+12
		{"Slytherin", 14},  // 8+6
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("rows differ (-want +got):\n%s", diff)
	}
}

func substr(col int, needle string) FilterFunc {
	return func(r Row) bool {
		s, _ := r[col].(string)
		return strings.Contains(s, needle)
	}
}
