package main

import "net/http"

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleWrite(w, r)
	default:
		http.Error(w, "Not Supported", http.StatusMethodNotAllowed)
	}
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		initCategories(w, r)
	default:
		http.Error(w, "Not Supported", http.StatusMethodNotAllowed)
	}
}
