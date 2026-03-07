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

	// employee routes
	router.HandlerFunc(http.MethodGet, "/v1/employees/:id", app.showEmployeeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/employees", app.listEmployeesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/employees", app.createEmployeeHandler)
	router.HandlerFunc(http.MethodPut, "/v1/employees/:id", app.updateEmployeeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/employees/:id", app.updateEmployeeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/employees/:id", app.deleteEmployeeHandler)

	// room routes
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms", app.listRoomsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number", app.showRoomHandler)
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms", app.createRoomHandler)
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

	// housekeeping_task routes
	router.HandlerFunc(http.MethodGet, "/v1/housekeeping_tasks/:taskID", app.showHousekeepingTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number/housekeeping_tasks", app.listHousekeepingTasksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms/:number/housekeeping_tasks", app.createHousekeepingTaskHandler)
	router.HandlerFunc(http.MethodPut, "/v1/housekeeping_tasks/:taskID", app.updateHousekeepingTaskHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/housekeeping_tasks/:taskID", app.updateHousekeepingTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/housekeeping_tasks/:taskID", app.deleteHousekeepingTaskHandler)

	// maintenance_report routes
	router.HandlerFunc(http.MethodGet, "/v1/maintenance_reports/:reportID", app.showMaintenanceReportHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number/maintenance_reports", app.listMaintenanceReportsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms/:number/maintenance_reports", app.createMaintenanceReportHandler)
	router.HandlerFunc(http.MethodPut, "/v1/maintenance_reports/:reportID", app.updateMaintenanceReportHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/maintenance_reports/:reportID", app.updateMaintenanceReportHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/maintenance_reports/:reportID", app.deleteMaintenanceReportHandler)

	// registration routes
	router.HandlerFunc(http.MethodGet, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.showRegistrationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/registrations/:reservationID", app.listRegistrationsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/registrations", app.createRegistrationHandler)
	router.HandlerFunc(http.MethodPut, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.updateRegistrationHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.updateRegistrationHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.deleteRegistrationHandler)

	// reservation routes
	router.HandlerFunc(http.MethodGet, "/v1/reservations/:reservationID", app.showReservationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/reservations", app.listReservationsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/reservations", app.createReservationHandler)
	router.HandlerFunc(http.MethodPut, "/v1/reservations/:reservationID", app.updateReservationHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/reservations/:reservationID", app.updateReservationHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/reservations/:reservationID", app.deleteReservationHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
