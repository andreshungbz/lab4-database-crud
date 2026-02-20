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
	Rooms RoomModel
}

// NewModels returns all Models configured with the database handler.
func NewModels(db *sql.DB) Models {
	return Models{
		Rooms: RoomModel{DB: db},
	}
}
