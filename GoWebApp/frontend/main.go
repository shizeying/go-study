package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/execute", executeHandler)
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	http.ListenAndServe(":8080", nil)
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var params Params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	result, err := ExecuteCmd(params)
	if err != nil {
		http.Error(w, "Failed to execute command", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
}
