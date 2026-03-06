package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createGuestHandler calls Guest.Insert.
// Writes JSON of the created guest record and its resource location.
func (app *application) createGuestHandler(w http.ResponseWriter, r *http.Request) {
	// Read JSON input into a Guest

	var input struct {
		PassportNumber string `json:"passport_number"`
		ContactEmail   string `json:"contact_email"`
		ContactPhone   string `json:"contact_phone"`
		Name           string `json:"name"`
		Gender         string `json:"gender"`
		Street         string `json:"street"`
		City           string `json:"city"`
		Country        string `json:"country"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	guest := &data.Guest{
		PassportNumber: input.PassportNumber,
		ContactEmail:   input.ContactEmail,
		ContactPhone:   input.ContactPhone,
		Name:           input.Name,
		Gender:         input.Gender,
		Street:         input.Street,
		City:           input.City,
		Country:        input.Country,
	}

	// validate
	v := validator.New()
	if data.ValidateGuest(v, guest); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// insert into database
	err = app.models.Guest.Insert(guest)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// add a header to indicate where the new resource is
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/guests/%s", guest.PassportNumber))

	// return JSON response of newly created guest
	err = app.writeJSON(w, http.StatusCreated, envelope{"guest": guest}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showGuestHandler calls Guest.Get.
// Writes JSON of the retrieved guest record.
func (app *application) showGuestHandler(w http.ResponseWriter, r *http.Request) {
	// read passport parameter
	passport := app.readStringParam("passport_number", r)
	if passport == "" {
		app.notFoundResponse(w, r)
		return
	}

	// retrieve guest from database
	guest, err := app.models.Guest.Get(passport)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// return JSON response of retrieved guest
	err = app.writeJSON(w, http.StatusOK, envelope{"guest": guest}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listGuestsHandler calls Guest.GetAll.
// Writes JSON of the list of filtered guest records and a metadata object.
func (app *application) listGuestsHandler(w http.ResponseWriter, r *http.Request) {
	// create input for filters (pagination + sort)
	var input struct {
		Name    string
		Country string
		data.Filters
	}

	// create validator and url.Values map
	v := validator.New()
	qs := r.URL.Query()

	// read parameters for filtering (search, pagination and sorting)
	input.Name = app.readURLString(qs, "name", "")                        // guest name
	input.Country = app.readURLString(qs, "country", "")                  // guest country
	input.Filters.Page = app.readInt(qs, "page", 1, v)                    // default: 1st page
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)          // default: 20 items per page
	input.Filters.Sort = app.readURLString(qs, "sort", "passport_number") // default: sort by passport number ascending
	input.Filters.SortSafelist = []string{                                // allowed sorting options
		"passport_number", "name", "created_at",
		"-passport_number", "-name", "-created_at",
	}

	// validate filters
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// retrieve records and pagination metadata from the database
	guests, metadata, err := app.models.Guest.GetAll(input.Name, input.Country, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// return JSON response of the list of guests
	err = app.writeJSON(w, http.StatusOK, envelope{"guests": guests, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateGuestHandler calls Guest.Update.
// Writes JSON of the updated guest record.
func (app *application) updateGuestHandler(w http.ResponseWriter, r *http.Request) {
	// read passport parameter
	passport := app.readStringParam("passport_number", r)
	if passport == "" {
		app.notFoundResponse(w, r)
		return
	}

	// retrieve guest from database
	guest, err := app.models.Guest.Get(passport)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Read JSON input

	var input struct {
		ContactEmail *string `json:"contact_email"`
		ContactPhone *string `json:"contact_phone"`
		Name         *string `json:"name"`
		Gender       *string `json:"gender"`
		Street       *string `json:"street"`
		City         *string `json:"city"`
		Country      *string `json:"country"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ContactEmail != nil {
		guest.ContactEmail = *input.ContactEmail
	}
	if input.ContactPhone != nil {
		guest.ContactPhone = *input.ContactPhone
	}
	if input.Name != nil {
		guest.Name = *input.Name
	}
	if input.Gender != nil {
		guest.Gender = *input.Gender
	}
	if input.Street != nil {
		guest.Street = *input.Street
	}
	if input.City != nil {
		guest.City = *input.City
	}
	if input.Country != nil {
		guest.Country = *input.Country
	}

	// validate
	v := validator.New()
	if data.ValidateGuest(v, guest); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// update record in the database
	err = app.models.Guest.Update(guest)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// return JSON response of updated guest
	err = app.writeJSON(w, http.StatusOK, envelope{"guest": guest}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteGuestHandler calls Guest.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteGuestHandler(w http.ResponseWriter, r *http.Request) {
	// read passport parameter
	passport := app.readStringParam("passport_number", r)
	if passport == "" {
		app.notFoundResponse(w, r)
		return
	}

	// delete guest and associated records from the database
	err := app.models.Guest.Delete(passport)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// return JSON response indicating success
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "guest successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
