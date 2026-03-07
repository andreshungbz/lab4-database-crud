package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// MaintenanceReport maps the maintenance_report entity.
type MaintenanceReport struct {
	ID            int64     `json:"id"`
	HotelID       int64     `json:"hotel_id"`
	RoomNumber    int       `json:"room_number"`
	HousekeeperID int64     `json:"housekeeper_id"`
	Description   string    `json:"description"`
	Completed     bool      `json:"completed"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

// ValidateMaintenanceReport performs validation checks for a maintenance_report record.
func ValidateMaintenanceReport(v *validator.Validator, r *MaintenanceReport) {
	v.Check(r.HotelID > 0, "hotel_id", "must be provided")
	v.Check(r.RoomNumber > 0, "room_number", "must be provided")
	v.Check(r.HousekeeperID > 0, "housekeeper_id", "must be provided")
	v.Check(r.Description != "", "description", "must be provided")
}

// MaintenanceReportModel holds the database handler.
type MaintenanceReportModel struct {
	DB *sql.DB
}

// Insert creates a maintenance_report record.
func (m MaintenanceReportModel) Insert(r *MaintenanceReport) error {
	query := `
		INSERT INTO maintenance_report
		(hotel_id, room_number, housekeeper_id, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, modified_at, completed`

	args := []any{r.HotelID, r.RoomNumber, r.HousekeeperID, r.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&r.ID,
		&r.CreatedAt,
		&r.ModifiedAt,
		&r.Completed,
	)
}

// Get retrieves a single maintenance_report record by id.
func (m MaintenanceReportModel) Get(id int64) (*MaintenanceReport, error) {
	query := `
		SELECT id, hotel_id, room_number, housekeeper_id, description,
		    completed, created_at, modified_at
		FROM maintenance_report
		WHERE id=$1`
	var r MaintenanceReport

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.HotelID,
		&r.RoomNumber,
		&r.HousekeeperID,
		&r.Description,
		&r.Completed,
		&r.CreatedAt,
		&r.ModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &r, nil
}

// GetAll retrieves multiple maintenance_report records (filterable).
func (m MaintenanceReportModel) GetAll(hotelID int64, roomNumber int, filters Filters) ([]*MaintenanceReport, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, hotel_id, room_number, housekeeper_id, description,
		    completed, created_at, modified_at
		FROM maintenance_report
		WHERE hotel_id=$1 AND room_number=$2
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []any{hotelID, roomNumber, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	reports := []*MaintenanceReport{}
	totalRecords := 0
	for rows.Next() {
		var r MaintenanceReport
		err := rows.Scan(&totalRecords, &r.ID, &r.HotelID, &r.RoomNumber, &r.HousekeeperID,
			&r.Description, &r.Completed, &r.CreatedAt, &r.ModifiedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		reports = append(reports, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return reports, metadata, nil
}

// Update modifies a maintenance_report record by id.
func (m MaintenanceReportModel) Update(r *MaintenanceReport) error {
	query := `
		UPDATE maintenance_report
		SET description=$1,
		    completed=$2,
		    modified_at=NOW()
		WHERE id=$3`

	args := []any{r.Description, r.Completed, r.ID}

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

// Delete removes a maintenance_report record by id.
func (m MaintenanceReportModel) Delete(id int64) error {
	query := `DELETE FROM maintenance_report WHERE id=$1`
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
