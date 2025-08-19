package main

import (
	"encoding/json"
	"net/http"
)

type QueryRequest struct {
	Query string `json:"query"` // e.g. "filter=^hp|^lotr;project=0,1,2;sort=2;limit=4"
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST supported", http.StatusMethodNotAllowed)
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	nodes := jsonToNodes(req.Query)
	nodes = append(nodes, &MemoryScan{table: db})
	head := q(nodes...)

	result := run(head)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
