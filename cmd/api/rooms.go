package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createRoomHandler reads JSON for a hotel room and writes it back to the client.
func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// use an intermediate struct to check validity of JSON
	var input struct {
		RoomNumber   int32  `json:"room_number"`
		RoomType     string `json:"room_type"`
		MaxOccupancy int32  `json:"max_occupancy"`
		HasBalcony   bool   `json:"has_balcony"`
		Available    bool   `json:"available"`
	}

	// attempt readJSON
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// create a Room from the input values
	room := &data.Room{
		RoomNumber:   input.RoomNumber,
		RoomType:     input.RoomType,
		MaxOccupancy: input.MaxOccupancy,
		HasBalcony:   input.HasBalcony,
		Available:    input.Available,
	}

	// validate JSON values
	v := validator.New()
	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// insert room into database
	err = app.models.Rooms.Insert(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// apply an HTTP header showing where the new resource is located
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/rooms/%d", room.ID))

	// use writeJSON for the HTTP response
	err = app.writeJSON(w, http.StatusCreated, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showRoomHandler writes JSON of a single hotel room to the HTTP response.
func (app *application) showRoomHandler(w http.ResponseWriter, r *http.Request) {
	// read and validate URL id parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// retrieve database record
	room, err := app.models.Rooms.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// use writeJSON for the HTTP response
	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
