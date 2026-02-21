package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes returns the HTTP router configured with all handlers, route-specific middleware,
// and global middleware.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Defined handlers for 404 and 205 status code
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck route
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Guest routes
	router.HandlerFunc(http.MethodGet, "/v1/guests/:passport", app.showGuestHandler)
	router.HandlerFunc(http.MethodGet, "/v1/guests", app.listGuestsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/guests", app.createGuestHandler)
	router.HandlerFunc(http.MethodPut, "/v1/guests/:passport", app.updateGuestHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/guests/:passport", app.updateGuestHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/guests/:passport", app.deleteGuestHandler)

	// Metrics debugging route
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.recoverPanic(router)
}
