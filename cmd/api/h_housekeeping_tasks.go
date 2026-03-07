package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// createHousekeepingTaskHandler calls HousekeepingTask.Insert.
// Writes JSON of the created housekeeping task record.
func (app *application) createHousekeepingTaskHandler(w http.ResponseWriter, r *http.Request) {
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
		HousekeeperID *int64 `json:"housekeeper_id"`
		TaskType      string `json:"task_type"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	task := &data.HousekeepingTask{
		HotelID:       hotelID,
		RoomNumber:    roomNumber,
		HousekeeperID: input.HousekeeperID,
		TaskType:      input.TaskType,
	}

	v := validator.New()
	if data.ValidateHousekeepingTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.HousekeepingTask.Insert(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	task, err = app.models.HousekeepingTask.Get(task.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"housekeeping_task": task}, nil)
}

// showHousekeepingTaskHandler calls HousekeepingTask.Get.
// Writes JSON of the retrieved housekeeping task record.
func (app *application) showHousekeepingTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("taskID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.models.HousekeepingTask.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"housekeeping_task": task}, nil)
}

// listHousekeepingTasksHandler calls HousekeepingTask.GetAll.
// Writes JSON of the list of housekeeping tasks belonging to a hotel room and a metadata object.
func (app *application) listHousekeepingTasksHandler(w http.ResponseWriter, r *http.Request) {
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
		HousekeeperID *int64
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.HousekeeperID, err = app.readOptionalInt64(qs, "housekeeper_id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "created_at")
	input.Filters.SortSafelist = []string{
		"id", "task_type", "completed", "created_at",
		"-id", "-task_type", "-completed", "-created_at",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	tasks, metadata, err := app.models.HousekeepingTask.GetAll(
		hotelID,
		roomNumber,
		input.HousekeeperID,
		input.Filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"housekeeping_tasks": tasks,
		"metadata":           metadata,
	}, nil)
}

// updateHousekeepingTaskHandler calls HousekeepingTask.Update.
// Writes JSON of the updated housekeeping task record.
func (app *application) updateHousekeepingTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("taskID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.models.HousekeepingTask.Get(id)
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
		HousekeeperID *int64  `json:"housekeeper_id"`
		TaskType      *string `json:"task_type"`
		Completed     *bool   `json:"completed"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.HousekeeperID != nil {
		task.HousekeeperID = input.HousekeeperID
	}

	if input.TaskType != nil {
		task.TaskType = *input.TaskType
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	v := validator.New()
	if data.ValidateHousekeepingTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.HousekeepingTask.Update(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	task, err = app.models.HousekeepingTask.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"housekeeping_task": task}, nil)
}

// deleteHousekeepingTaskHandler calls HousekeepingTask.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteHousekeepingTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("taskID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.HousekeepingTask.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "housekeeping task successfully deleted"}, nil)
}
