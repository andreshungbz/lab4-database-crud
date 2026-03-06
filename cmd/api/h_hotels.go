package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createHotelHandler calls Hotel.Insert.
// Writes JSON of the created hotel record and its resource location.
func (app *application) createHotelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
		Phone   string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	h := &data.Hotel{
		Name:    input.Name,
		Street:  input.Street,
		City:    input.City,
		State:   input.State,
		Country: input.Country,
		Phone:   input.Phone,
	}

	v := validator.New()
	if data.ValidateHotel(v, h); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Hotel.Insert(h)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/hotels/%d", h.ID))

	app.writeJSON(w, http.StatusCreated, envelope{"hotel": h}, headers)
}

// showHotelHandler calls Hotel.Get.
// Writes JSON of the retrieved hotel record.
func (app *application) showHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	h, err := app.models.Hotel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotel": h}, nil)
}

// listHotelsHandler calls Hotel.GetAll.
// Writes JSON of the list of filtered hotel records and a metadata object.
func (app *application) listHotelsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readURLString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "name", "city",
		"-id", "-name", "-city",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hotels, metadata, err := app.models.Hotel.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotels": hotels, "metadata": metadata}, nil)
}

// updateHotelHandler calls Hotel.Update.
// Writes JSON of the updated hotel record.
func (app *application) updateHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	h, err := app.models.Hotel.Get(id)
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
		Name    *string `json:"name"`
		Street  *string `json:"street"`
		City    *string `json:"city"`
		State   *string `json:"state"`
		Country *string `json:"country"`
		Phone   *string `json:"phone"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		h.Name = *input.Name
	}
	if input.Street != nil {
		h.Street = *input.Street
	}
	if input.City != nil {
		h.City = *input.City
	}
	if input.State != nil {
		h.State = *input.State
	}
	if input.Country != nil {
		h.Country = *input.Country
	}
	if input.Phone != nil {
		h.Phone = *input.Phone
	}

	v := validator.New()

	if data.ValidateHotel(v, h); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Hotel.Update(h)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotel": h}, nil)
}

// deleteHotelHandler calls Hotel.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Hotel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "hotel successfully deleted"}, nil)
}
