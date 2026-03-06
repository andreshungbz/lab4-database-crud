package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Hotel maps the hotel entity.
type Hotel struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
}

// ValidateHotel performs validation checks for a hotel record.
func ValidateHotel(v *validator.Validator, h *Hotel) {
	v.Check(h.Name != "", "name", "must be provided")
	v.Check(h.Street != "", "street", "must be provided")
	v.Check(h.City != "", "city", "must be provided")
	v.Check(h.State != "", "state", "must be provided")
	v.Check(h.Country != "", "country", "must be provided")
	v.Check(h.Phone != "", "phone", "must be provided")
}

// HotelModel holds the database handler.
type HotelModel struct {
	DB *sql.DB
}

// Insert creates a hotel record.
func (m HotelModel) Insert(h *Hotel) error {
	query := `
		INSERT INTO hotel (name, street, city, state, country, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []any{
		h.Name,
		h.Street,
		h.City,
		h.State,
		h.Country,
		h.Phone,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&h.ID)
}

// Get retrieves a single hotel record by id.
func (m HotelModel) Get(id int64) (*Hotel, error) {
	query := `
		SELECT id, name, street, city, state, country, phone
		FROM hotel
		WHERE id = $1`
	var h Hotel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&h.ID,
		&h.Name,
		&h.Street,
		&h.City,
		&h.State,
		&h.Country,
		&h.Phone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &h, nil
}

// GetAll retrieves multiple hotel records (filterable).
func (m HotelModel) GetAll(name string, filters Filters) ([]*Hotel, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),
			id,
			name,
			street,
			city,
			state,
			country,
			phone
		FROM hotel
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection(),
	)

	args := []any{name, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	hotels := []*Hotel{}
	for rows.Next() {
		var h Hotel

		err := rows.Scan(
			&totalRecords,
			&h.ID,
			&h.Name,
			&h.Street,
			&h.City,
			&h.State,
			&h.Country,
			&h.Phone,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		hotels = append(hotels, &h)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return hotels, metadata, nil
}

// Update modifies a hotel record by id.
func (m HotelModel) Update(h *Hotel) error {
	query := `
		UPDATE hotel
		SET name=$1,
			street=$2,
			city=$3,
			state=$4,
			country=$5,
			phone=$6
		WHERE id=$7`

	args := []any{
		h.Name,
		h.Street,
		h.City,
		h.State,
		h.Country,
		h.Phone,
		h.ID,
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

// Delete removes a hotel record by id (cascades to employee and room).
func (m HotelModel) Delete(id int64) error {
	query := `DELETE FROM hotel WHERE id = $1`

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
