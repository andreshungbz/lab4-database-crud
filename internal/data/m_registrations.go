package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Registration maps the registration entity.
type Registration struct {
	ReservationID int64 `json:"reservation_id"`
	HotelID       int64 `json:"hotel_id"`
	RoomNumber    int   `json:"room_number"`
}

// ValidateRegistration performs validation checks for a registration record.
func ValidateRegistration(v *validator.Validator, r *Registration) {
	v.Check(r.ReservationID > 0, "reservation_id", "must be provided")
	v.Check(r.HotelID > 0, "hotel_id", "must be provided")
	v.Check(r.RoomNumber > 0, "room_number", "must be provided")
}

// RegistrationModel holds the database handler.
type RegistrationModel struct {
	DB *sql.DB
}

// Insert creates a registration record.
func (m RegistrationModel) Insert(r *Registration) error {
	query := `SELECT fn_create_registration($1,$2,$3)`

	args := []any{r.ReservationID, r.HotelID, r.RoomNumber}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// Get retrieves a single registration record.
func (m RegistrationModel) Get(reservationID int64, hotelID int64, roomNumber int) (*Registration, error) {
	query := `
		SELECT reservation_id, hotel_id, room_number
		FROM registration
		WHERE reservation_id=$1 AND hotel_id=$2 AND room_number=$3`
	var r Registration

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, reservationID, hotelID, roomNumber).Scan(
		&r.ReservationID,
		&r.HotelID,
		&r.RoomNumber,
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

// GetAll retrieves all registration records belonging to a reservation.
func (m RegistrationModel) GetAll(reservationID int64, filters Filters) ([]*Registration, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), reservation_id, hotel_id, room_number
		FROM registration
		WHERE reservation_id=$1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection())

	args := []any{reservationID, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	registrations := []*Registration{}
	totalRecords := 0
	for rows.Next() {
		var r Registration

		err := rows.Scan(
			&totalRecords,
			&r.ReservationID,
			&r.HotelID,
			&r.RoomNumber,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		registrations = append(registrations, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return registrations, metadata, nil
}

// Update modifies a registration record.
func (m RegistrationModel) Update(reservationID int64, hotelID int64, roomNumber int, newRoomNumber int) error {
	query := `
		UPDATE registration
		SET room_number=$1
		WHERE reservation_id=$2 AND hotel_id=$3 AND room_number=$4`

	args := []any{newRoomNumber, reservationID, hotelID, roomNumber}

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

// Delete removes a registration record.
func (m RegistrationModel) Delete(reservationID int64, hotelID int64, roomNumber int) error {
	query := `
		DELETE FROM registration
		WHERE reservation_id=$1 AND hotel_id=$2 AND room_number=$3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, reservationID, hotelID, roomNumber)
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
