// Binary main runs an implementatioin of a DB query executor.
package main

import (
	"dbms/data"
	"dbms/queryexecutor"
	"fmt"
)

func main() {
	// schema: id, first_name, last_name, house, points
	rows := []queryexecutor.Row{
		{1, "Harry", "Potter", "Sly", 60},
		{2, "Hermione", "Granger", "Rav", 50.5},
		{3, "Ron", "Weasley", "Huf", 71.3},
		{4, "Draco", "Malfoy", "Sly", 11.3},
	}
	mapper := func(r queryexecutor.Row) queryexecutor.Row { return queryexecutor.Row{r[1], r[3]} }

	result := queryexecutor.Run(queryexecutor.Q(
		&queryexecutor.Limit{L: 2},
		&queryexecutor.Projection{Mapper: mapper},
		&queryexecutor.MemoryScan{Rows: rows},
	))

	if len(result) != 2 {
		fmt.Printf("len(result) = %v, want 2", len(result))
	} else {
		fmt.Println("OK")
	}

	const fileName = "data/movies_test.csv"
	data.Encode(fileName)
}
