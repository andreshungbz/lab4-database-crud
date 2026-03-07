package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// ====================================================================================
// Handlers for Reservation Entity
// ====================================================================================

// createReservationHandler creates a reservation and optionally registers a room.
func (app *application) createReservationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GuestID       int64  `json:"guest_id"`
		CheckinDate   string `json:"checkin_date"`
		CheckoutDate  string `json:"checkout_date"`
		PaymentMethod string `json:"payment_method"`
		Source        string `json:"source"`
		RoomTypeID    int    `json:"room_type_id"`
		HotelID       int    `json:"hotel_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reservation := &data.Reservation{
		GuestID:       input.GuestID,
		CheckinDate:   input.CheckinDate,
		CheckoutDate:  input.CheckoutDate,
		PaymentMethod: input.PaymentMethod,
		Source:        input.Source,
	}

	v := validator.New()
	if data.ValidateReservation(v, reservation); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	id, err := app.models.Reservation.Insert(reservation, input.RoomTypeID, input.HotelID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newRes, err := app.models.Reservation.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"reservation": newRes}, nil)
}

// showReservationHandler returns a reservation with its registrations.
func (app *application) showReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	reservation, err := app.models.Reservation.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"reservation": reservation}, nil)
}

// listReservationsHandler returns all reservations.
func (app *application) listReservationsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "-id", "guest_id", "-guest_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reservations, metadata, err := app.models.Reservation.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"reservations": reservations, "metadata": metadata}, nil)
}

// updateReservationHandler modifies a reservation.
func (app *application) updateReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		CheckinDate   string  `json:"checkin_date"`
		CheckoutDate  string  `json:"checkout_date"`
		PaymentAmount float64 `json:"payment_amount"`
		PaymentMethod string  `json:"payment_method"`
		Source        string  `json:"source"`
		Completed     bool    `json:"completed"`
		Canceled      bool    `json:"canceled"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reservation := &data.Reservation{
		ID:            id,
		CheckinDate:   input.CheckinDate,
		CheckoutDate:  input.CheckoutDate,
		PaymentAmount: input.PaymentAmount,
		PaymentMethod: input.PaymentMethod,
		Source:        input.Source,
		Completed:     input.Completed,
		Canceled:      input.Canceled,
	}

	err = app.models.Reservation.Update(reservation)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	newRes, err := app.models.Reservation.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"reservation": newRes}, nil)
}

// deleteReservationHandler removes a reservation.
func (app *application) deleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reservationID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Reservation.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "reservation successfully deleted"}, nil)
}
