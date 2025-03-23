package main

import (
	"net/http"
)

/*
This file is created to check if health-check API
is proper or not

Imp note - For swagger ensure we are not having a space before func and
swagger comments or it wont work
*/

// healthCheckHandler provides a simple health check for the API.
//
//	@Summary		Health Check
//	@Description	Returns the status of the API along with environment and version details.
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"API is healthy"
//	@Failure		500	{object}	error				"Internal server error"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create JSON response
	response := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version, //const from main.go
	}

	// call writeJSON wrapper from json.go for better json response handling
	if err := writeJSON(w, http.StatusOK, response); err != nil {
		writeJSONWrapper(w, http.StatusInternalServerError, err.Error())
	}

}
