// Binary main runs an implementatioin of a DB query executor.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

var db []Row = []Row{
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

func main() {
	t := flag.Bool("t", false, "Run test")
	flag.Parse()
	if *t {
		// runTest()
		runTestReadFromCSV()
		return
	}

	http.HandleFunc("/query", handleQuery)
	fmt.Println("Read-only DB server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
