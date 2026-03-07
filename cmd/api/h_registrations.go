package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createRegistrationHandler calls Registration.Insert.
// Writes JSON of the created registration record.
func (app *application) createRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ReservationID int64 `json:"reservation_id"`
		HotelID       int64 `json:"hotel_id"`
		RoomNumber    int   `json:"room_number"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reg := &data.Registration{
		ReservationID: input.ReservationID,
		HotelID:       input.HotelID,
		RoomNumber:    input.RoomNumber,
	}

	v := validator.New()
	if data.ValidateRegistration(v, reg); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Registration.Insert(reg)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"registration": reg}, nil)
}

// showRegistrationHandler calls Registration.Get.
// Writes JSON of the retrieved registration record.
func (app *application) showRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	reservationID, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotelID, err := app.readInt64Param("hotelID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	roomNumber, err := app.readIntParam("roomNumber", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	reg, err := app.models.Registration.Get(reservationID, hotelID, roomNumber)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"registration": reg}, nil)
}

// listRegistrationsHandler calls Registration.GetAll.
// Writes JSON of the list of registrations belonging to a reservation.
func (app *application) listRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	reservationID, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "room_number")

	input.Filters.SortSafelist = []string{
		"hotel_id", "room_number",
		"-hotel_id", "-room_number",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	registrations, metadata, err := app.models.Registration.GetAll(
		reservationID,
		input.Filters,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"registrations": registrations,
		"metadata":      metadata,
	}, nil)
}

// updateRegistrationHandler calls Registration.Update.
// Writes JSON of the updated registration record.
func (app *application) updateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	reservationID, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotelID, err := app.readInt64Param("hotelID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	roomNumber, err := app.readIntParam("roomNumber", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		RoomNumber int `json:"room_number"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Registration.Update(reservationID, hotelID, roomNumber, input.RoomNumber)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	reg, err := app.models.Registration.Get(reservationID, hotelID, input.RoomNumber)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"registration": reg}, nil)
}

// deleteRegistrationHandler calls Registration.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	reservationID, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotelID, err := app.readInt64Param("hotelID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	roomNumber, err := app.readIntParam("roomNumber", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Registration.Delete(reservationID, hotelID, roomNumber)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "registration successfully deleted"}, nil)
}
