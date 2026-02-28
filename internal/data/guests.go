package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Guest maps the guest entity, which is a subtype of the person entity.
type Guest struct {
	// guest attributes
	ID             int64  `json:"-"`
	PassportNumber string `json:"passport_number"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
	// person attributes
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"-"`
}

// ValidateGuest checks for the passport number.
func ValidateGuest(v *validator.Validator, guest *Guest) {
	v.Check(guest.PassportNumber != "", "passport_number", "must be provided")
}

// GuestModel holds a handler to the database
type GuestModel struct {
	DB *sql.DB
}

// Insert creates a record in tables person and guest.
func (g GuestModel) Insert(guest *Guest) error {
	query := `SELECT * FROM fn_create_guest($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		guest.PassportNumber,
		guest.ContactEmail,
		guest.ContactPhone,
		guest.Name,
		guest.Gender,
		guest.Street,
		guest.City,
		guest.Country,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return g.DB.QueryRowContext(ctx, query, args...).Scan(
		// scan remaining attributes
		&guest.ID,
		&guest.CreatedAt,
	)
}

// Get reads a guest's passport and returns a Guest.
func (g GuestModel) Get(passport string) (*Guest, error) {
	query := `SELECT * FROM fn_get_guest($1)`
	var guest Guest

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := g.DB.QueryRowContext(ctx, query, passport).Scan(
		// scan all attributes
		&guest.ID,
		&guest.PassportNumber,
		&guest.ContactEmail,
		&guest.ContactPhone,
		&guest.Name,
		&guest.Gender,
		&guest.Street,
		&guest.City,
		&guest.Country,
		&guest.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &guest, nil
}

// GetAll reads all guests in the database, filterable.
func (g GuestModel) GetAll(name string, country string, filters Filters) ([]*Guest, Metadata, error) {
	// PostgreSQL full-text search notes
	// - to_tsvector in simple configuration breaks string to lower case lexemes.
	// - plainto_tsquery in simple configuration normalizes the query term.
	//		e.g. "John Smith" -> 'john' + 'smith'
	// - @@ is the matches operator.
	query := fmt.Sprintf(`
		SELECT
			count(*) OVER(), 
			g.id,
			g.passport_number,
			g.contact_email,
			g.contact_phone,
			p.name,
			p.gender,
			p.street,
			p.city,
			p.country,
			p.created_at
		FROM guest g
		JOIN person p ON p.id = g.id
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', country) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	// limit is used for page_size, and offset is used for page
	args := []any{name, country, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// retrieves rows from the database
	rows, err := g.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	// construct the array of guests
	totalRecords := 0
	guests := []*Guest{}
	for rows.Next() {
		var guest Guest

		err := rows.Scan(
			&totalRecords,
			// scan all attributes
			&guest.ID,
			&guest.PassportNumber,
			&guest.ContactEmail,
			&guest.ContactPhone,
			&guest.Name,
			&guest.Gender,
			&guest.Street,
			&guest.City,
			&guest.Country,
			&guest.CreatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		guests = append(guests, &guest)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// construct Metadata object
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return guests, metadata, nil
}

// Update modifies the appropriate person and guest records for a guest.
func (g GuestModel) Update(guest *Guest) error {
	query := `SELECT fn_update_guest($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		guest.PassportNumber,
		guest.ContactEmail,
		guest.ContactPhone,
		guest.Name,
		guest.Gender,
		guest.Street,
		guest.City,
		guest.Country,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// No rows being affected means the guest's passport is not in the database

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes a guest from the database and their associated reservations
// and registrations.
func (g GuestModel) Delete(passport string) error {
	query := `SELECT fn_delete_guest($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, passport)
	if err != nil {
		return err
	}

	// No rows being affected means the guest's passport is not in the database

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
