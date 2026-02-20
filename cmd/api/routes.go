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

	// Room routes
	router.HandlerFunc(http.MethodGet, "/v1/rooms/:id", app.showRoomHandler)
	router.HandlerFunc(http.MethodPost, "/v1/rooms", app.createRoomHandler)

	// Metrics debugging route
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.recoverPanic(router)
}
