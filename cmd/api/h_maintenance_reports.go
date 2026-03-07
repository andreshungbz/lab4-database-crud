package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createMaintenanceReportHandler calls MaintenanceReport.Insert.
// Writes JSON of the created maintenance report record.
func (app *application) createMaintenanceReportHandler(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		HousekeeperID int64  `json:"housekeeper_id"`
		Description   string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	report := &data.MaintenanceReport{
		HotelID:       hotelID,
		RoomNumber:    roomNumber,
		HousekeeperID: input.HousekeeperID,
		Description:   input.Description,
	}

	v := validator.New()
	if data.ValidateMaintenanceReport(v, report); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.MaintenanceReport.Insert(report)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	report, err = app.models.MaintenanceReport.Get(report.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"maintenance_report": report}, nil)
}

// showMaintenanceReportHandler calls MaintenanceReport.Get.
// Writes JSON of the retrieved maintenance report record.
func (app *application) showMaintenanceReportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reportID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	report, err := app.models.MaintenanceReport.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"maintenance_report": report}, nil)
}

// listMaintenanceReportsHandler calls MaintenanceReport.GetAll.
// Writes JSON of the list of maintenance reports belonging to a hotel room and a metadata object.
func (app *application) listMaintenanceReportsHandler(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "created_at")
	input.Filters.SortSafelist = []string{
		"id", "description", "completed", "created_at",
		"-id", "-description", "-completed", "-created_at",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reports, metadata, err := app.models.MaintenanceReport.GetAll(
		hotelID,
		roomNumber,
		input.Filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"maintenance_reports": reports,
		"metadata":            metadata,
	}, nil)
}

// updateMaintenanceReportHandler calls MaintenanceReport.Update.
// Writes JSON of the updated maintenance report record.
func (app *application) updateMaintenanceReportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reportID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	report, err := app.models.MaintenanceReport.Get(id)
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
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Description != nil {
		report.Description = *input.Description
	}

	if input.Completed != nil {
		report.Completed = *input.Completed
	}

	v := validator.New()
	if data.ValidateMaintenanceReport(v, report); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.MaintenanceReport.Update(report)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	report, err = app.models.MaintenanceReport.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"maintenance_report": report}, nil)
}

// deleteMaintenanceReportHandler calls MaintenanceReport.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteMaintenanceReportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("reportID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.MaintenanceReport.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "maintenance report successfully deleted"}, nil)
}
