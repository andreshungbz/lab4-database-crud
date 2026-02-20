package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Room maps a hotel room entity.
type Room struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"-"`
	RoomNumber   int32     `json:"room_number"`
	RoomType     string    `json:"room_type"`
	MaxOccupancy int32     `json:"max_occupancy"`
	HasBalcony   bool      `json:"has_balcony"`
	Available    bool      `json:"available"`
}

// ValidateRoom does simple validation to ensure RoomNumber and MaxOccupancy
// are positive values.
func ValidateRoom(v *validator.Validator, room *Room) {
	v.Check(room.RoomNumber > 0, "room_number", "Room number must be a positive number")
	v.Check(room.MaxOccupancy > 0, "max_occupancy", "Max occupancy must be a positive number")
}

// RoomModel defines CRUD operations with the database for Room.
type RoomModel struct {
	DB *sql.DB
}

// Insert creates a new room database record.
func (m RoomModel) Insert(room *Room) error {
	query := `
		INSERT INTO room (room_number, room_type, max_occupancy, has_balcony, available)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	args := []any{room.RoomNumber, room.RoomType, room.MaxOccupancy, room.HasBalcony, room.Available}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&room.ID, &room.CreatedAt)
}

// Get retrieves a room database record by its ID.
func (m RoomModel) Get(id int64) (*Room, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, room_number, room_type, max_occupancy, has_balcony, available
		FROM room
		WHERE id = $1`

	var room Room

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.CreatedAt,
		&room.RoomNumber,
		&room.RoomType,
		&room.MaxOccupancy,
		&room.HasBalcony,
		&room.Available,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return &room, nil
}
