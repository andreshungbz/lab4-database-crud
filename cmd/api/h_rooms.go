package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createRoomHandler calls Room.Insert.
// Writes JSON of the created room record.
func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	hotelID, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Number     int    `json:"number"`
		RoomTypeID int    `json:"room_type_id"`
		Floor      int    `json:"floor"`
		StatusCode string `json:"status_code"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	room := &data.Room{
		HotelID:    hotelID,
		Number:     input.Number,
		RoomTypeID: input.RoomTypeID,
		Floor:      input.Floor,
		StatusCode: input.StatusCode,
	}

	v := validator.New()
	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Room.Insert(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	room, err = app.models.Room.Get(hotelID, room.Number)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"room": room}, nil)
}

// showRoomHandler calls Room.Get.
// Writes JSON of the retrieved room record.
func (app *application) showRoomHandler(w http.ResponseWriter, r *http.Request) {
	hotelID, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	roomNumber, err := app.readIntParam("number", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	room, err := app.models.Room.Get(hotelID, roomNumber)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
}

// listRoomsHandler calls Room.GetAll.
// Writes JSON of the list of rooms belonging to a hotel and a metadata object.
func (app *application) listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	hotelID, err := app.readInt64Param("id", r)
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
	input.Filters.Sort = app.readURLString(qs, "sort", "number")
	input.Filters.SortSafelist = []string{
		"number", "floor", "status_code",
		"-number", "-floor", "-status_code",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	rooms, metadata, err := app.models.Room.GetAll(hotelID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"rooms": rooms, "metadata": metadata}, nil)
}

// updateRoomHandler calls Room.Update.
// Writes JSON of the updated room record.
func (app *application) updateRoomHandler(w http.ResponseWriter, r *http.Request) {
	hotelID, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	number, err := app.readIntParam("number", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	room, err := app.models.Room.Get(hotelID, number)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		RoomTypeID *int    `json:"room_type_id"`
		Floor      *int    `json:"floor"`
		StatusCode *string `json:"status_code"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.RoomTypeID != nil {
		room.RoomTypeID = *input.RoomTypeID
	}

	if input.Floor != nil {
		room.Floor = *input.Floor
	}

	if input.StatusCode != nil {
		room.StatusCode = *input.StatusCode
	}

	v := validator.New()
	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Room.Update(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	room, err = app.models.Room.Get(hotelID, number)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
}

// deleteRoomHandler calls Room.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	hotelID, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	number, err := app.readIntParam("number", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Room.Delete(hotelID, number)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "room successfully deleted"}, nil)
}
