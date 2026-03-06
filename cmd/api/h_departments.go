package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createDepartmentHandler calls Department.Insert.
// Writes JSON of the created department record and its resource location.
func (app *application) createDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DeptName string  `json:"dept_name"`
		Budget   float64 `json:"budget"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	d := &data.Department{
		DeptName: input.DeptName,
		Budget:   input.Budget,
	}

	v := validator.New()
	if data.ValidateDepartment(v, d); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Department.Insert(d)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"department": d}, nil)
}

// showDepartmentHandler calls Department.Get.
// Writes JSON of the retrieved department record.
func (app *application) showDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	name := app.readStringParam("dept_name", r)
	if name == "" {
		app.notFoundResponse(w, r)
		return
	}

	d, err := app.models.Department.Get(name)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"department": d}, nil)
}

// listDepartmentsHandler calls Department.GetAll.
// Writes JSON of the list of filtered department records and a metadata object.
func (app *application) listDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readURLString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "dept_name")
	input.Filters.SortSafelist = []string{
		"dept_name", "budget",
		"-dept_name", "-budget",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	departments, metadata, err := app.models.Department.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"departments": departments, "metadata": metadata}, nil)
}

// updateDepartmentHandler calls Department.Update.
// Writes JSON of the updated department record.
func (app *application) updateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	name := app.readStringParam("dept_name", r)
	if name == "" {
		app.notFoundResponse(w, r)
		return
	}

	d, err := app.models.Department.Get(name)

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
		Budget *float64 `json:"budget"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Budget != nil {
		d.Budget = *input.Budget
	}

	v := validator.New()
	if data.ValidateDepartment(v, d); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Department.Update(d)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"department": d}, nil)
}

// deleteDepartmentHandler calls Department.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	name := app.readStringParam("dept_name", r)
	if name == "" {
		app.notFoundResponse(w, r)
		return
	}

	err := app.models.Department.Delete(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "department successfully deleted"}, nil)
}
