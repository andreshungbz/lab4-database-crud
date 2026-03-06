package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createEmployeeHandler calls Employee.Insert.
// Writes JSON of the created employee record.
func (app *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		HotelID    int64   `json:"hotel_id"`
		Department string  `json:"department"`
		ManagerID  *int64  `json:"manager_id,omitempty"`
		Salary     float64 `json:"salary"`
		SSN        string  `json:"ssn"`
		WorkEmail  string  `json:"work_email"`
		WorkPhone  string  `json:"work_phone"`
		Password   string  `json:"password"`
		Role       string  `json:"role"`
		HotelOwner *bool   `json:"hotel_owner,omitempty"`
		Shift      *string `json:"shift,omitempty"`
		Name       string  `json:"name"`
		Gender     string  `json:"gender"`
		Street     string  `json:"street"`
		City       string  `json:"city"`
		Country    string  `json:"country"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	employee := &data.Employee{
		HotelID:      input.HotelID,
		Department:   input.Department,
		ManagerID:    input.ManagerID,
		Salary:       input.Salary,
		SSN:          input.SSN,
		WorkEmail:    input.WorkEmail,
		WorkPhone:    input.WorkPhone,
		PasswordHash: []byte(input.Password), // TODO: hashing
		Role:         input.Role,
		HotelOwner:   input.HotelOwner,
		Shift:        input.Shift,
		Name:         input.Name,
		Gender:       input.Gender,
		Street:       input.Street,
		City:         input.City,
		Country:      input.Country,
	}

	v := validator.New()
	if data.ValidateEmployee(v, employee); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Employee.Insert(employee)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"employee": employee}, nil)
}

// showEmployeeHandler calls Employee.Get.
// Writes JSON of the retrieved employee record.
func (app *application) showEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	employee, err := app.models.Employee.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employee": employee}, nil)
}

// listEmployeesHandler calls Employee.GetAll (optional filters can be added later).
// Writes JSON of the list of employees.
func (app *application) listEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
		Role string
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Role = app.readURLString(qs, "role", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "name", "department",
		"-id", "-name", "-department",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	employees, metadata, err := app.models.Employee.GetAll(input.Role, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employees": employees, "metadata": metadata}, nil)
}

// updateEmployeeHandler calls Employee.Update.
// Writes JSON of the updated employee record.
func (app *application) updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	employee, err := app.models.Employee.Get(id)
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
		HotelID    *int64   `json:"hotel_id,omitempty"`
		Department *string  `json:"department"`
		ManagerID  *int64   `json:"manager_id,omitempty"`
		Salary     *float64 `json:"salary"`
		WorkEmail  *string  `json:"work_email"`
		WorkPhone  *string  `json:"work_phone"`
		Role       *string  `json:"role"`
		HotelOwner *bool    `json:"hotel_owner,omitempty"`
		Shift      *string  `json:"shift,omitempty"`
		Name       *string  `json:"name"`
		Gender     *string  `json:"gender"`
		Street     *string  `json:"street"`
		City       *string  `json:"city"`
		Country    *string  `json:"country"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.HotelID != nil {
		employee.HotelID = *input.HotelID
	}

	if input.Department != nil {
		employee.Department = *input.Department
	}
	if input.ManagerID != nil {
		employee.ManagerID = input.ManagerID
	}
	if input.Salary != nil {
		employee.Salary = *input.Salary
	}
	if input.WorkEmail != nil {
		employee.WorkEmail = *input.WorkEmail
	}
	if input.WorkPhone != nil {
		employee.WorkPhone = *input.WorkPhone
	}
	if input.Role != nil {
		employee.Role = *input.Role
	}
	if input.HotelOwner != nil {
		employee.HotelOwner = input.HotelOwner
	}
	if input.Shift != nil {
		employee.Shift = input.Shift
	}
	if input.Name != nil {
		employee.Name = *input.Name
	}
	if input.Gender != nil {
		employee.Gender = *input.Gender
	}
	if input.Street != nil {
		employee.Street = *input.Street
	}
	if input.City != nil {
		employee.City = *input.City
	}
	if input.Country != nil {
		employee.Country = *input.Country
	}

	v := validator.New()
	if data.ValidateEmployee(v, employee); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Employee.Update(employee)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employee": employee}, nil)
}

// deleteEmployeeHandler calls Employee.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Employee.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "employee successfully deleted"}, nil)
}
