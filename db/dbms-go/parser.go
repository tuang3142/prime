package main

import (
	"fmt"
	"regexp"
	"strings"
)

func jsonToNodes(input string) []Node {
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

	nodes := []Node{}
	parts := regexp.MustCompile(`\s*;\s*`).Split(input, -1)
	for _, part := range parts {
		switch {
		case part == "scan":
			// This is just a placeholder; actual MemoryScan must be inserted at call site with the db
			nodes = append(nodes, nil)

		case strings.HasPrefix(part, "filter="):
			patterns := regexp.MustCompile(`\|`).Split(part[len("filter="):], -1)
			var fs []FilterFunc
			for _, pat := range patterns {
				fs = append(fs, match(pat))
			}
			nodes = append(nodes, &Selection{filter: Or(fs...)})

		case strings.HasPrefix(part, "project="):
			fields := strings.Split(part[len("project="):], ",")
			nodes = append(nodes, &Projection{
				mapper: func(row Row) Row {
					var projected Row
					for _, f := range fields {
						var i int
						fmt.Sscanf(f, "%d", &i)
						projected = append(projected, row[i])
					}
					return projected
				},
			})

		case strings.HasPrefix(part, "sort="):
			spec := strings.TrimPrefix(part, "sort=")
			var col int
			desc := false
			if _, err := fmt.Sscanf(spec, "%d desc", &col); err == nil {
				desc = true
			} else {
				fmt.Sscanf(spec, "%d", &col)
			}
			nodes = append(nodes, &Sort{
				key: func(row Row) float64 {
					v, ok := row[col].(float64)
					if !ok {
						panic(fmt.Sprintf("row[%d] is not float64", col))
					}
					return v
				},
				desc: desc,
			})

		case strings.HasPrefix(part, "limit="):
			var lim int
			fmt.Sscanf(part[len("limit="):], "%d", &lim)
			nodes = append(nodes, &Limit{limit: lim})

		default:
			panic("unknown instruction: " + part)
		}
	}
	return nodes
}
