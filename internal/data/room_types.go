package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// RoomType maps the room_type entity.
type RoomType struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	BaseRate     float64 `json:"base_rate"`
	MaxOccupancy int     `json:"max_occupancy"`
	BedCount     int     `json:"bed_count"`
	HasBalcony   bool    `json:"has_balcony"`
}

// ValidateRoomType checks for a required title and positive numbers.
func ValidateRoomType(v *validator.Validator, rt *RoomType) {
	v.Check(rt.Title != "", "title", "must be provided")
	v.Check(rt.BaseRate > 0, "base_rate", "must be greater than 0")
	v.Check(rt.MaxOccupancy > 0, "max_occupancy", "must be greater than 0")
	v.Check(rt.BedCount > 0, "bed_count", "must be greater than 0")
}

// RoomTypeModel holds a handler to the database.
type RoomTypeModel struct {
	DB *sql.DB
}

// Insert creates a record in table room_type.
func (m RoomTypeModel) Insert(rt *RoomType) error {
	query := `
		INSERT INTO room_type (title, base_rate, max_occupancy, bed_count, has_balcony)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	args := []any{
		rt.Title,
		rt.BaseRate,
		rt.MaxOccupancy,
		rt.BedCount,
		rt.HasBalcony,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&rt.ID)
}

// Get retrieves room type information by its ID.
func (m RoomTypeModel) Get(id int64) (*RoomType, error) {
	query := `
		SELECT id, title, base_rate, max_occupancy, bed_count, has_balcony
		FROM room_type
		WHERE id = $1`
	var rt RoomType

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&rt.ID,
		&rt.Title,
		&rt.BaseRate,
		&rt.MaxOccupancy,
		&rt.BedCount,
		&rt.HasBalcony,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rt, nil
}

// GetAll reads all room types in the database (filterable).
func (m RoomTypeModel) GetAll(title string, filters Filters) ([]*RoomType, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),
			id,
			title,
			base_rate,
			max_occupancy,
			bed_count,
			has_balcony
		FROM room_type
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection(),
	)

	args := []any{title, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	roomTypes := []*RoomType{}
	for rows.Next() {
		var rt RoomType

		err := rows.Scan(
			&totalRecords,
			&rt.ID,
			&rt.Title,
			&rt.BaseRate,
			&rt.MaxOccupancy,
			&rt.BedCount,
			&rt.HasBalcony,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		roomTypes = append(roomTypes, &rt)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return roomTypes, metadata, nil
}

// Update modifies a room type.
func (m RoomTypeModel) Update(rt *RoomType) error {
	query := `
		UPDATE room_type
		SET title = $1,
			base_rate = $2,
			max_occupancy = $3,
			bed_count = $4,
			has_balcony = $5
		WHERE id = $6`

	args := []any{
		rt.Title,
		rt.BaseRate,
		rt.MaxOccupancy,
		rt.BedCount,
		rt.HasBalcony,
		rt.ID,
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

// Delete removes a room type (cascades)
func (m RoomTypeModel) Delete(id int64) error {
	query := `DELETE FROM room_type WHERE id = $1`

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
