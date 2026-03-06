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

	// Metrics debugging route
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// DATABASE SCHEMA ROUTES

	// guest routes
	router.HandlerFunc(http.MethodGet, "/v1/guests/:passport_number", app.showGuestHandler)
	router.HandlerFunc(http.MethodGet, "/v1/guests", app.listGuestsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/guests", app.createGuestHandler)
	router.HandlerFunc(http.MethodPut, "/v1/guests/:passport_number", app.updateGuestHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/guests/:passport_number", app.updateGuestHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/guests/:passport_number", app.deleteGuestHandler)

	// hotel routes
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id", app.showHotelHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/hotels", app.createHotelHandler)
	router.HandlerFunc(http.MethodPut, "/v1/hotels/:id", app.updateHotelHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/hotels/:id", app.updateHotelHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/hotels/:id", app.deleteHotelHandler)

	// department routes
	router.HandlerFunc(http.MethodGet, "/v1/departments/:dept_name", app.showDepartmentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/departments", app.listDepartmentsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/departments", app.createDepartmentHandler)
	router.HandlerFunc(http.MethodPut, "/v1/departments/:dept_name", app.updateDepartmentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/departments/:dept_name", app.updateDepartmentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/departments/:dept_name", app.deleteDepartmentHandler)

	// room routes
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms", app.listRoomsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms", app.createRoomHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number", app.showRoomHandler)
	router.HandlerFunc(http.MethodPut, "/v1/hotels/:id/rooms/:number", app.updateRoomHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/hotels/:id/rooms/:number", app.updateRoomHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/hotels/:id/rooms/:number", app.deleteRoomHandler)

	// room_type routes
	router.HandlerFunc(http.MethodGet, "/v1/room_types/:id", app.showRoomTypeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/room_types", app.listRoomTypesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/room_types", app.createRoomTypeHandler)
	router.HandlerFunc(http.MethodPut, "/v1/room_types/:id", app.updateRoomTypeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/room_types/:id", app.updateRoomTypeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/room_types/:id", app.deleteRoomTypeHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
