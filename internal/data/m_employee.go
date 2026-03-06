package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/validator"
)

// Employee maps the employee entity (subtype of person).
type Employee struct {
	// employee attributes
	ID           int64   `json:"id"`
	HotelID      int64   `json:"hotel_id"`
	Department   string  `json:"department"`
	ManagerID    *int64  `json:"manager_id,omitempty"`
	Salary       float64 `json:"salary"`
	SSN          string  `json:"ssn"`
	WorkEmail    string  `json:"work_email"`
	WorkPhone    string  `json:"work_phone"`
	PasswordHash []byte  `json:"-"`
	// role attributes
	Role       string  `json:"role"`
	HotelOwner *bool   `json:"hotel_owner,omitempty"`
	Shift      *string `json:"shift,omitempty"`
	// person attributes
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

// ValidateEmployee performs validation checks for an employee record.
func ValidateEmployee(v *validator.Validator, e *Employee) {
	v.Check(e.HotelID > 0, "hotel_id", "must be provided")
	v.Check(e.Department != "", "department", "must be provided")
	v.Check(e.Salary >= 0, "salary", "must be provided")
	v.Check(e.SSN != "", "ssn", "must be provided")
	v.Check(e.WorkEmail != "", "work_email", "must be provided")
	v.Check(e.WorkPhone != "", "work_phone", "must be provided")
	v.Check(e.Role != "", "role", "must be provided")

	v.Check(e.Name != "", "name", "must be provided")
	v.Check(e.Gender != "", "gender", "must be provided")
	v.Check(e.Street != "", "street", "must be provided")
	v.Check(e.City != "", "city", "must be provided")
	v.Check(e.Country != "", "country", "must be provided")

	switch e.Role {
	case "operations_manager":
		v.Check(e.HotelOwner != nil, "hotel_owner", "must be provided for operations_manager")
	case "front_desk", "housekeeper":
		v.Check(e.Shift != nil && *e.Shift != "", "shift", "must be provided for front_desk or housekeeper")
	}
}

// EmployeeModel holds the database handler.
type EmployeeModel struct {
	DB *sql.DB
}

// Insert creates a new employee (person + employee + role-specific).
func (m EmployeeModel) Insert(e *Employee) error {
	query := `SELECT * FROM fn_create_employee($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	args := []any{
		// person attributes
		e.Name,
		e.Gender,
		e.Street,
		e.City,
		e.Country,
		// employee attributes
		e.HotelID,
		e.Department,
		e.ManagerID,
		e.Salary,
		e.SSN,
		e.WorkEmail,
		e.WorkPhone,
		e.PasswordHash,
		e.Role,
		// role-specific attributes
		e.HotelOwner,
		e.Shift,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&e.ID, &e.CreatedAt, &e.ModifiedAt)
}

// Get retrieves a single employee record by id.
func (m EmployeeModel) Get(id int64) (*Employee, error) {
	query := `SELECT * FROM fn_get_employee($1)`

	var e Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.HotelID,
		&e.Department,
		&e.ManagerID,
		&e.Salary,
		&e.SSN,
		&e.WorkEmail,
		&e.WorkPhone,
		&e.PasswordHash,
		&e.Role,
		&e.HotelOwner,
		&e.Shift,
		&e.Name,
		&e.Gender,
		&e.Street,
		&e.City,
		&e.Country,
		&e.CreatedAt,
		&e.ModifiedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &e, nil
}

// GetAll retrieves multiple employee records with optional role filter and pagination.
func (m EmployeeModel) GetAll(role string, filters Filters) ([]*Employee, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			count(*) OVER(),
			e.id,
			e.hotel_id,
			e.department,
			e.manager_id,
			e.salary,
			e.ssn,
			e.work_email,
			e.work_phone,
			CASE
				WHEN om.id IS NOT NULL THEN 'operations_manager'
				WHEN fd.id IS NOT NULL THEN 'front_desk'
				WHEN hk.id IS NOT NULL THEN 'housekeeper'
			END AS role,
			om.hotel_owner,
			COALESCE(fd.shift, hk.shift) AS shift,
			p.name,
			p.gender,
			p.street,
			p.city,
			p.country,
			p.created_at,
			p.modified_at
		FROM employee e
		JOIN person p ON p.id = e.id
		LEFT JOIN operations_manager om ON om.id = e.id
		LEFT JOIN front_desk fd ON fd.id = e.id
		LEFT JOIN housekeeper hk ON hk.id = e.id
		WHERE ($1 = '' OR
			(om.id IS NOT NULL AND $1 = 'operations_manager') OR
			(fd.id IS NOT NULL AND $1 = 'front_desk') OR
			(hk.id IS NOT NULL AND $1 = 'housekeeper')
		)
		ORDER BY %s %s, e.id ASC
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection(),
	)

	args := []any{role, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	employees := []*Employee{}
	for rows.Next() {
		var e Employee
		err := rows.Scan(
			&totalRecords,
			&e.ID,
			&e.HotelID,
			&e.Department,
			&e.ManagerID,
			&e.Salary,
			&e.SSN,
			&e.WorkEmail,
			&e.WorkPhone,
			&e.Role,
			&e.HotelOwner,
			&e.Shift,
			&e.Name,
			&e.Gender,
			&e.Street,
			&e.City,
			&e.Country,
			&e.CreatedAt,
			&e.ModifiedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		employees = append(employees, &e)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return employees, metadata, nil
}

// Update modifies an existing employee record (person + employee + role-specific).
func (m EmployeeModel) Update(e *Employee) error {
	query := `
	SELECT fn_update_employee(
    	$1,$2,$3,$4,$5,
    	$6,
    	$7,$8,$9,$10,
    	$11,$12,
    	$13,$14,$15
	)`

	args := []any{
		// person attributes
		e.Name,
		e.Gender,
		e.Street,
		e.City,
		e.Country,
		// employee attributes
		e.ID,
		e.HotelID,
		e.Department,
		e.ManagerID,
		e.Salary,
		e.WorkEmail,
		e.WorkPhone,
		// role-specific attributes
		e.Role,
		e.HotelOwner,
		e.Shift,
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

// Delete removes an employee record by id.
func (m EmployeeModel) Delete(id int64) error {
	query := `SELECT fn_delete_employee($1)`

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
