package main

import (
	"fmt"
	"reflect"
	"regexp"
)

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

func runTestReadFromCSV() {
	isAnimation := func(r Row) bool {
		genre, ok := r[2].(string)
		if !ok {
			panic("can't get genre")
		}
		matched, err := regexp.Match(`Animation`, []byte(genre))
		return err != nil && matched
	}
	mapper := func(r Row) Row { return Row{r[1], r[2]} }
	// TODO: isGood join with rating
	got := run(q(
		&Projection{mapper: mapper},
		&Limit{limit: 5},
		&Selection{filter: isAnimation},
		&FileScan{},
	))
	for _, r := range got {
		fmt.Println(r)
	}
}

func runTest() {
	match := func(regex string) FilterFunc {
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
		{"lotr", 3},
		{"hp", 3},
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

	agg := func(rows []Row) Row {
		var s float64
		for _, row := range rows {
			rating := row[2].(float64)
			s += rating
		}
		return Row{s}
	}
	want = []Row{
		[]any{"hp1", 22.9},
		[]any{"lotr", 26.6},
	}
	got = run(q(
		&GroupBy{groupId: 4, agg: agg},
		&Selection{filter: Or(match("^hp"), match("^lotr"))},
		&MemoryScan{table: db},
	))
	if err := assert(got, want); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
