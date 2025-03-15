package main

import (
	"encoding/json"
	"net/http"
)

/*
when the data is written to response Writer,
that's treated as response payload from curl, the json
in w will be shown in response
*/
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	/*
		This encodes the data struct into JSON and writes it to the response body.
		If encoding fails, it returns an error.
	*/
	return json.NewEncoder(w).Encode(data)
}

//readJSON
/*
This function reads JSON data from an HTTP request payload
and decodes it into the provided data structure.
It also enforces a size limit and prevents unknown fields.
*/
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	/*
		This prevents large payloads from
		overwhelming the server.
		Any request body larger than ~1MB will be rejected.
	*/

	maxBytes := 1_048_578 // 1 mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	//This initializes a decoder to read JSON from the request body.
	decoder := json.NewDecoder(r.Body)
	/*
		If the JSON request contains fields
		that are not defined in the data struct,
		it rejects the request.
	*/
	decoder.DisallowUnknownFields()
	/*
		This attempts to decode the JSON into the data struct.
		If there’s an error
		(e.g., invalid JSON, unknown fields, or exceeding size limits),
		it returns an error.
	*/
	return decoder.Decode(data)
}

//writeJSONError
/*
This function is a helper for sending JSON-formatted error messages in HTTP responses.
*/
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//Defines an anonymous struct with a single field, error, to hold the error message.
	type envelope struct {
		Error string `json:"error"` //for json error formatting
	}
	return writeJSON(w, status, envelope{Error: message})
}

//writeJSONWrapper
/*
This is a wrapper to wrap around the response coming from writeJSON
in envelop struct.
*/
func writeJSONWrapper(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, envelope{Data: data})
}
