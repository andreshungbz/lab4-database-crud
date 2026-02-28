package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createRoomTypeHandler reads JSON input to create a room type.
func (app *application) createRoomTypeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string  `json:"title"`
		BaseRate     float64 `json:"base_rate"`
		MaxOccupancy int     `json:"max_occupancy"`
		BedCount     int     `json:"bed_count"`
		HasBalcony   bool    `json:"has_balcony"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	rt := &data.RoomType{
		Title:        input.Title,
		BaseRate:     input.BaseRate,
		MaxOccupancy: input.MaxOccupancy,
		BedCount:     input.BedCount,
		HasBalcony:   input.HasBalcony,
	}

	v := validator.New()
	if data.ValidateRoomType(v, rt); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RoomType.Insert(rt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/room-types/%d", rt.ID))

	app.writeJSON(w, http.StatusCreated, envelope{"room_type": rt}, headers)
}

// showRoomTypeHandler returns a JSON response of a room type by its ID.
func (app *application) showRoomTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rt, err := app.models.RoomType.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"room_type": rt}, nil)
}

// listRoomTypesHandler returns all room types (filterable).
func (app *application) listRoomTypesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "title", "base_rate",
		"-id", "-title", "-base_rate",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	roomTypes, metadata, err := app.models.RoomType.GetAll(input.Title, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"room_types": roomTypes, "metadata": metadata}, nil)
}

// updateRoomTypeHandler reads JSON input and updates the corresponding room type by ID.
func (app *application) updateRoomTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rt, err := app.models.RoomType.Get(id)
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
		Title        *string  `json:"title"`
		BaseRate     *float64 `json:"base_rate"`
		MaxOccupancy *int     `json:"max_occupancy"`
		BedCount     *int     `json:"bed_count"`
		HasBalcony   *bool    `json:"has_balcony"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		rt.Title = *input.Title
	}
	if input.BaseRate != nil {
		rt.BaseRate = *input.BaseRate
	}
	if input.MaxOccupancy != nil {
		rt.MaxOccupancy = *input.MaxOccupancy
	}
	if input.BedCount != nil {
		rt.BedCount = *input.BedCount
	}
	if input.HasBalcony != nil {
		rt.HasBalcony = *input.HasBalcony
	}

	v := validator.New()
	if data.ValidateRoomType(v, rt); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RoomType.Update(rt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"room_type": rt}, nil)
}

// deleteRoomTypeHandler removes a room type by its ID (cascades).
func (app *application) deleteRoomTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.RoomType.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "room type successfully deleted"}, nil)
}
