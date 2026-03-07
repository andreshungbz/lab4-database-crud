package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// HousekeepingTask maps the housekeeping_task entity.
type HousekeepingTask struct {
	ID            int64     `json:"id"`
	HotelID       int64     `json:"hotel_id"`
	RoomNumber    int       `json:"room_number"`
	HousekeeperID *int64    `json:"housekeeper_id,omitempty"`
	TaskType      string    `json:"task_type"`
	Completed     bool      `json:"completed"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

// ValidateHousekeepingTask performs validation checks for a housekeeping_task record.
func ValidateHousekeepingTask(v *validator.Validator, t *HousekeepingTask) {
	v.Check(t.HotelID > 0, "hotel_id", "must be provided")
	v.Check(t.RoomNumber > 0, "room_number", "must be provided")
	v.Check(t.TaskType != "", "task_type", "must be provided")
}

// HousekeepingTaskModel holds the database handler.
type HousekeepingTaskModel struct {
	DB *sql.DB
}

// Insert creates a housekeeping_task record.
func (m HousekeepingTaskModel) Insert(t *HousekeepingTask) error {
	query := `
		INSERT INTO housekeeping_task
		(hotel_id, room_number, housekeeper_id, task_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, modified_at, completed`

	args := []any{
		t.HotelID,
		t.RoomNumber,
		t.HousekeeperID,
		t.TaskType,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&t.ID,
		&t.CreatedAt,
		&t.ModifiedAt,
		&t.Completed,
	)
}

// Get retrieves a single housekeeping_task record by id.
func (m HousekeepingTaskModel) Get(id int64) (*HousekeepingTask, error) {
	query := `
		SELECT id, hotel_id, room_number, housekeeper_id,
		    task_type, completed, created_at, modified_at
		FROM housekeeping_task
		WHERE id = $1`
	var t HousekeepingTask

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.HotelID,
		&t.RoomNumber,
		&t.HousekeeperID,
		&t.TaskType,
		&t.Completed,
		&t.CreatedAt,
		&t.ModifiedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &t, nil
}

// GetAll retrieves multiple housekeeping_task records (filterable).
func (m HousekeepingTaskModel) GetAll(hotelID int64, roomNumber int, housekeeperID *int64, filters Filters) ([]*HousekeepingTask, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			count(*) OVER(),
			id,
			hotel_id,
			room_number,
			housekeeper_id,
			task_type,
			completed,
			created_at,
			modified_at
		FROM housekeeping_task
		WHERE hotel_id = $1
		AND room_number = $2
		AND ($3::BIGINT IS NULL OR housekeeper_id = $3)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`,
		filters.sortColumn(), filters.sortDirection(),
	)

	args := []any{
		hotelID,
		roomNumber,
		housekeeperID,
		filters.limit(),
		filters.offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tasks := []*HousekeepingTask{}
	for rows.Next() {
		var t HousekeepingTask

		err := rows.Scan(
			&totalRecords,
			&t.ID,
			&t.HotelID,
			&t.RoomNumber,
			&t.HousekeeperID,
			&t.TaskType,
			&t.Completed,
			&t.CreatedAt,
			&t.ModifiedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

// Update modifies a housekeeping_task record by id.
func (m HousekeepingTaskModel) Update(t *HousekeepingTask) error {
	query := `
		UPDATE housekeeping_task
		SET housekeeper_id=$1,
		    task_type=$2,
		    completed=$3,
		    modified_at=NOW()
		WHERE id=$4`

	args := []any{
		t.HousekeeperID,
		t.TaskType,
		t.Completed,
		t.ID,
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

// Delete removes a housekeeping_task record by id.
func (m HousekeepingTaskModel) Delete(id int64) error {
	query := `DELETE FROM housekeeping_task WHERE id = $1`

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
