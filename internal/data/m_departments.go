package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Department maps the department entity.
type Department struct {
	DeptName string  `json:"dept_name"`
	Budget   float64 `json:"budget"`
}

// ValidateDepartment performs validation checks for a department record.
func ValidateDepartment(v *validator.Validator, d *Department) {
	v.Check(d.DeptName != "", "dept_name", "must be provided")
	v.Check(d.Budget >= 0, "budget", "must be greater than or equal to 0")
}

// DepartmentModel holds the database handler.
type DepartmentModel struct {
	DB *sql.DB
}

// Insert creates a department record.
func (m DepartmentModel) Insert(d *Department) error {
	query := `
		INSERT INTO department (dept_name, budget)
		VALUES ($1, $2)`

	args := []any{
		d.DeptName,
		d.Budget,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// Get retrieves a single department record by dept_name.
func (m DepartmentModel) Get(name string) (*Department, error) {
	query := `
		SELECT dept_name, budget
		FROM department
		WHERE dept_name = $1`
	var d Department

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, name).Scan(
		&d.DeptName,
		&d.Budget,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &d, nil
}

// GetAll retrieves multiple department records (filterable).
func (m DepartmentModel) GetAll(name string, filters Filters) ([]*Department, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),
			dept_name,
			budget
		FROM department
		WHERE (to_tsvector('simple', dept_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, dept_name ASC
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
	departments := []*Department{}

	for rows.Next() {
		var d Department

		err := rows.Scan(
			&totalRecords,
			&d.DeptName,
			&d.Budget,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		departments = append(departments, &d)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return departments, metadata, nil
}

// Update modifies a department record by dept_name.
func (m DepartmentModel) Update(d *Department) error {
	query := `
		UPDATE department
		SET budget = $1
		WHERE dept_name = $2`

	args := []any{
		d.Budget,
		d.DeptName,
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

// Delete removes a department record by ID (cascades to employee).
func (m DepartmentModel) Delete(name string) error {
	query := `DELETE FROM department WHERE dept_name = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, name)
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
