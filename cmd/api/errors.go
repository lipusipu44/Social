package main

import (
	"log"
	"net/http"
)

//internalServerError
/*
This method will be used for all the errors in handler go files.
app *application is not required, but will be used when we will use
log in future and initiation of log will be happening in app struct,

r *http.Request is only used for log purpose, no role apart from that
*/
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request server error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) retryResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("serialization failure: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, "serialization failure, please retry")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("deadlock detected: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, "deadlock detected, please retry")
}
