package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Room maps the room entity.
type Room struct {
	HotelID    int64     `json:"hotel_id"`
	Number     int       `json:"number"`
	RoomTypeID int       `json:"-"`
	Floor      int       `json:"floor"`
	StatusCode string    `json:"status_code"`
	ModifiedAt time.Time `json:"modified_at"`

	RoomType RoomType `json:"room_type,omitzero"`
}

// ValidateRoom performs validation checks for a room record.
func ValidateRoom(v *validator.Validator, r *Room) {
	v.Check(r.HotelID > 0, "hotel_id", "must be provided")
	v.Check(r.Number > 0, "number", "must be provided")
	v.Check(r.RoomTypeID > 0, "room_type_id", "must be provided")
	v.Check(r.Floor > 0, "floor", "must be provided")
	v.Check(r.StatusCode != "", "status_code", "must be provided")
}

// RoomModel holds the database handler.
type RoomModel struct {
	DB *sql.DB
}

// Insert creates a room record.
func (m RoomModel) Insert(r *Room) error {
	query := `
		INSERT INTO room (hotel_id, number, room_type_id, floor, status_code)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING modified_at`

	args := []any{
		r.HotelID,
		r.Number,
		r.RoomTypeID,
		r.Floor,
		r.StatusCode,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&r.ModifiedAt)
}

// Get retrieves a single room record by hotel id and room number.
func (m RoomModel) Get(hotelID int64, number int) (*Room, error) {
	query := `
		SELECT
			r.hotel_id,
			r.number,
			r.room_type_id,
			r.floor,
			r.status_code,
			r.modified_at,
			rt.id,
			rt.title,
			rt.base_rate,
			rt.max_occupancy,
			rt.bed_count,
			rt.has_balcony
		FROM room r
		JOIN room_type rt ON r.room_type_id = rt.id
		WHERE r.hotel_id = $1 AND r.number = $2`
	var r Room

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, hotelID, number).Scan(
		&r.HotelID,
		&r.Number,
		&r.RoomTypeID,
		&r.Floor,
		&r.StatusCode,
		&r.ModifiedAt,
		&r.RoomType.ID,
		&r.RoomType.Title,
		&r.RoomType.BaseRate,
		&r.RoomType.MaxOccupancy,
		&r.RoomType.BedCount,
		&r.RoomType.HasBalcony,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &r, nil
}

// GetAll retrieves all rooms belonging to a specific hotel.
func (m RoomModel) GetAll(hotelID int64, filters Filters) ([]*Room, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			count(*) OVER(),
			r.hotel_id,
			r.number,
			r.room_type_id,
			r.floor,
			r.status_code,
			r.modified_at,
			rt.id,
			rt.title,
			rt.base_rate,
			rt.max_occupancy,
			rt.bed_count,
			rt.has_balcony
		FROM room r
		JOIN room_type rt ON r.room_type_id = rt.id
		WHERE r.hotel_id = $1
		ORDER BY %s %s, r.number ASC
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection(),
	)

	args := []any{hotelID, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	rooms := []*Room{}
	for rows.Next() {
		var r Room

		err := rows.Scan(
			&totalRecords,
			&r.HotelID,
			&r.Number,
			&r.RoomTypeID,
			&r.Floor,
			&r.StatusCode,
			&r.ModifiedAt,
			&r.RoomType.ID,
			&r.RoomType.Title,
			&r.RoomType.BaseRate,
			&r.RoomType.MaxOccupancy,
			&r.RoomType.BedCount,
			&r.RoomType.HasBalcony,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		rooms = append(rooms, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return rooms, metadata, nil
}

// Update modifies a room record by hotel id and room number.
func (m RoomModel) Update(r *Room) error {
	query := `
		UPDATE room
		SET room_type_id=$1,
			floor=$2,
			status_code=$3,
			modified_at=NOW()
		WHERE hotel_id=$4 AND number=$5`

	args := []any{
		r.RoomTypeID,
		r.Floor,
		r.StatusCode,
		r.HotelID,
		r.Number,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes a room record by hotel id and room number (cascades).
func (m RoomModel) Delete(hotelID int64, number int) error {
	query := `DELETE FROM room WHERE hotel_id = $1 AND number = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, hotelID, number)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
