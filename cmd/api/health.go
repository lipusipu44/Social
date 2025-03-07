package main

import (
	"encoding/json"
	"net/http"
)

/*
This file is created to check if health-check API
is proper or not
*/

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create JSON response
	response := map[string]string{"status": "OK"}

	// Encode JSON and write response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}
