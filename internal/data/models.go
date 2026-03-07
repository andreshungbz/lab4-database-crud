package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models groups all database models used in the application.
type Models struct {
	Guest             GuestModel
	Hotel             HotelModel
	Department        DepartmentModel
	Employee          EmployeeModel
	Room              RoomModel
	RoomType          RoomTypeModel
	HousekeepingTask  HousekeepingTaskModel
	MaintenanceReport MaintenanceReportModel
	Registration      RegistrationModel
	Reservation       ReservationModel
}

// NewModels returns all Models configured with the database handler.
func NewModels(db *sql.DB) Models {
	return Models{
		Guest:             GuestModel{DB: db},
		Hotel:             HotelModel{DB: db},
		Department:        DepartmentModel{DB: db},
		Employee:          EmployeeModel{DB: db},
		Room:              RoomModel{DB: db},
		RoomType:          RoomTypeModel{DB: db},
		HousekeepingTask:  HousekeepingTaskModel{DB: db},
		MaintenanceReport: MaintenanceReportModel{DB: db},
		Registration:      RegistrationModel{DB: db},
		Reservation:       ReservationModel{DB: db},
	}
}
