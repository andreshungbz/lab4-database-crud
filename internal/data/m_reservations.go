package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Reservation maps the reservation entity and includes its registrations.
type Reservation struct {
	ID            int64           `json:"id"`
	GuestID       int64           `json:"guest_id"`
	CheckinDate   string          `json:"checkin_date"`
	CheckoutDate  string          `json:"checkout_date"`
	PaymentAmount float64         `json:"payment_amount"`
	PaymentMethod string          `json:"payment_method"`
	Source        string          `json:"source"`
	Completed     bool            `json:"completed"`
	Canceled      bool            `json:"canceled"`
	CreatedAt     string          `json:"created_at"`
	ModifiedAt    string          `json:"modified_at"`
	Registrations []*Registration `json:"registrations"`
}

// ValidateReservation performs validation checks for a reservation record.
func ValidateReservation(v *validator.Validator, r *Reservation) {
	v.Check(r.GuestID > 0, "guest_id", "must be provided")
	v.Check(r.CheckinDate != "", "checkin_date", "must be provided")
	v.Check(r.CheckoutDate != "", "checkout_date", "must be provided")
	v.Check(r.PaymentAmount >= 0, "payment_amount", "must be non-negative")
	v.Check(r.PaymentMethod != "", "payment_method", "must be provided")
	v.Check(r.Source != "", "source", "must be provided")
}

// ReservationModel holds the database handler.
type ReservationModel struct {
	DB *sql.DB
}

// Insert creates a reservation record and optionally registers a room.
func (m ReservationModel) Insert(r *Reservation, roomTypeID int, hotelID int) (int64, error) {
	query := `SELECT fn_create_reservation_workflow($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		r.GuestID,
		r.CheckinDate,
		r.CheckoutDate,
		r.PaymentMethod,
		r.Source,
		hotelID,
		roomTypeID,
		nil, // new reservation, no existing id
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservationID int64
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&reservationID)
	if err != nil {
		return 0, err
	}

	return reservationID, nil
}

// Get retrieves a single reservation along with its registrations.
func (m ReservationModel) Get(id int64) (*Reservation, error) {
	query := `
		SELECT id, guest_id, checkin_date, checkout_date,
		    payment_amount, payment_method, source,
		    completed, canceled, created_at, modified_at
		FROM reservation
		WHERE id=$1`
	var r Reservation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.GuestID,
		&r.CheckinDate,
		&r.CheckoutDate,
		&r.PaymentAmount,
		&r.PaymentMethod,
		&r.Source,
		&r.Completed,
		&r.Canceled,
		&r.CreatedAt,
		&r.ModifiedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Get registrations for this reservation
	regModel := RegistrationModel{DB: m.DB}
	regs, _, err := regModel.GetAll(
		r.ID,
		Filters{Page: 1, PageSize: 1000, Sort: "room_number", SortSafelist: []string{"room_number"}},
	)
	if err != nil {
		return nil, err
	}
	r.Registrations = regs

	return &r, nil
}

// GetAll retrieves all reservations with optional pagination.
func (m ReservationModel) GetAll(filters Filters) ([]*Reservation, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, guest_id, checkin_date, checkout_date,
		    payment_amount, payment_method, source,
		    completed, canceled, created_at, modified_at
		FROM reservation
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`,
		filters.sortColumn(), filters.sortDirection())

	args := []any{filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var reservations []*Reservation
	totalRecords := 0
	for rows.Next() {
		var r Reservation
		err := rows.Scan(
			&totalRecords,
			&r.ID,
			&r.GuestID,
			&r.CheckinDate,
			&r.CheckoutDate,
			&r.PaymentAmount,
			&r.PaymentMethod,
			&r.Source,
			&r.Completed,
			&r.Canceled,
			&r.CreatedAt,
			&r.ModifiedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// attach registrations
		regModel := RegistrationModel{DB: m.DB}
		regs, _, err := regModel.GetAll(
			r.ID,
			Filters{Page: 1, PageSize: 1000, Sort: "room_number", SortSafelist: []string{"room_number"}},
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		r.Registrations = regs

		reservations = append(reservations, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return reservations, metadata, nil
}

// Update modifies a reservation record.
func (m ReservationModel) Update(r *Reservation) error {
	query := `
		UPDATE reservation
		SET checkin_date=$1, checkout_date=$2, payment_amount=$3,
		    payment_method=$4, source=$5, completed=$6, canceled=$7,
		    modified_at=NOW()
		WHERE id=$8`

	args := []any{
		r.CheckinDate,
		r.CheckoutDate,
		r.PaymentAmount,
		r.PaymentMethod,
		r.Source,
		r.Completed,
		r.Canceled,
		r.ID,
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

// Delete removes a reservation record.
func (m ReservationModel) Delete(id int64) error {
	query := `DELETE FROM reservation WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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
