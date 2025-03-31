package main

import (
	"go.uber.org/zap"
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
	app.zapLogger.Error("internal server error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Warn("bad request server error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("not found error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) retryResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("serialization failure error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusConflict, "serialization failure, please retry")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("deadlock detected error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusConflict, "deadlock detected, please retry")
}

func (app *application) uniqueContraintConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("unique Constraint detected error", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusConflict, "unique key violation error detected")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("unauthorized error response", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	//at-least without this I was not getting pop up for admin uid and pwd in browser
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized error detected")
}

func (app *application) unauthorizedTokenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.zapLogger.Error("unauthorized error response", zap.String("METHOD", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusUnauthorized, "unauthorized error detected")
}
