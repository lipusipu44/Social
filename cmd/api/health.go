package main

import (
	"net/http"
)

/*
This file is created to check if health-check API
is proper or not
*/

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create JSON response
	response := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version, //const from main.go
	}

	// call writeJSON wrapper from json.go for better json response handling
	if err := writeJSON(w, http.StatusOK, response); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

}
