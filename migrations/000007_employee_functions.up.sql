-- migrations/employee_functions.up.sql
-- Creates functions for employees (subtype of person).

-- ====================================================================================
-- CREATE FUNCTION fn_create_employee returns the person id and created_at for a newly
-- created employee. Handles conditional insertion into role-specific tables.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_employee(
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT,
    -- employee attributes
    p_hotel_id BIGINT,
    p_department TEXT,
    p_manager_id BIGINT,
    p_salary NUMERIC,
    p_ssn TEXT,
    p_work_email CITEXT,
    p_work_phone TEXT,
    p_password_hash BYTEA,
    -- role-specific attributes
    p_role TEXT,  -- "operations_manager", "front_desk", "housekeeper"
    p_hotel_owner BOOLEAN DEFAULT NULL,
    p_shift TEXT DEFAULT NULL
)
RETURNS TABLE (
    id BIGINT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE
)
AS $$
DECLARE
    v_employee_id BIGINT;
BEGIN
    -- insert person entry
    INSERT INTO person (name, gender, street, city, country)
    VALUES (p_name, p_gender, p_street, p_city, p_country)
    RETURNING person.id INTO v_employee_id;

    -- insert employee entry
    INSERT INTO employee (id, hotel_id, department, manager_id, salary, ssn, work_email, work_phone, password_hash)
    VALUES (v_employee_id, p_hotel_id, p_department, p_manager_id, p_salary, p_ssn, p_work_email, p_work_phone, p_password_hash);

    -- conditional insertion into role-specific tables
    IF p_role = 'operations_manager' THEN
        INSERT INTO operations_manager (id, hotel_owner)
        VALUES (v_employee_id, p_hotel_owner);
    ELSIF p_role = 'front_desk' THEN
        INSERT INTO front_desk (id, shift)
        VALUES (v_employee_id, p_shift::shift_type);
    ELSIF p_role = 'housekeeper' THEN
        INSERT INTO housekeeper (id, shift)
        VALUES (v_employee_id, p_shift::shift_type);
    END IF;

    RETURN QUERY
    SELECT
        p.id,
        p.created_at,
        p.modified_at
    FROM person p
    WHERE p.id = v_employee_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- READ FUNCTION fn_get_employee returns employee and person data for a given employee id.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_employee(
    p_id BIGINT
)
RETURNS TABLE (
    id BIGINT,
    hotel_id BIGINT,
    department TEXT,
    manager_id BIGINT,
    salary NUMERIC,
    ssn TEXT,
    work_email CITEXT,
    work_phone TEXT,
    password_hash BYTEA,
    role TEXT,
    hotel_owner BOOLEAN,
    shift shift_type,
    name TEXT,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE
)
AS $$
DECLARE
    v_role TEXT;
BEGIN
    -- determine role
    IF EXISTS (SELECT 1 FROM operations_manager om WHERE om.id = p_id) THEN
        v_role := 'operations_manager';
    ELSIF EXISTS (SELECT 1 FROM front_desk fd WHERE fd.id = p_id) THEN
        v_role := 'front_desk';
    ELSIF EXISTS (SELECT 1 FROM housekeeper hk WHERE hk.id = p_id) THEN
        v_role := 'housekeeper';
    ELSE
        RAISE EXCEPTION '[employee-not-found] Employee with id % does not exist', p_id;
    END IF;

    RETURN QUERY
    SELECT
        e.id,
        e.hotel_id,
        e.department,
        e.manager_id,
        e.salary,
        e.ssn,
        e.work_email,
        e.work_phone,
        e.password_hash,
        v_role,
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
    WHERE e.id = p_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- UPDATE FUNCTION fn_update_employee updates person, employee, and optional
-- role-specific fields.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_update_employee(
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT,
    -- employee attributes
    p_id BIGINT,
    p_hotel_id BIGINT,
    p_department TEXT,
    p_manager_id BIGINT,
    p_salary NUMERIC,
    p_work_email CITEXT,
    p_work_phone TEXT,
    -- role-specific attributes
    p_role TEXT,
    p_hotel_owner BOOLEAN DEFAULT NULL,
    p_shift TEXT DEFAULT NULL
)
RETURNS VOID
AS $$
BEGIN
    -- update employee
    UPDATE employee
    SET hotel_id = p_hotel_id,
        department = p_department,
        manager_id = p_manager_id,
        salary = p_salary,
        work_email = p_work_email,
        work_phone = p_work_phone
    WHERE id = p_id;

    -- update role-specific tables
    IF p_role = 'operations_manager' THEN
        UPDATE operations_manager
        SET hotel_owner = p_hotel_owner
        WHERE id = p_id;
    ELSIF p_role = 'front_desk' THEN
        UPDATE front_desk
        SET shift = p_shift::shift_type
        WHERE id = p_id;
    ELSIF p_role = 'housekeeper' THEN
        UPDATE housekeeper
        SET shift = p_shift::shift_type
        WHERE id = p_id;
    END IF;

    -- update person
    UPDATE person
    SET name = p_name,
        gender = p_gender,
        street = p_street,
        city = p_city,
        country = p_country
    WHERE id = p_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- DELETE FUNCTION fn_delete_employee deletes an employee (via person table) and cascades.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_delete_employee(
    p_id BIGINT
)
RETURNS VOID
AS $$
BEGIN
    -- delete from person (cascades to employee and subtype tables)
    DELETE FROM person
    WHERE id = p_id;
END;
$$ LANGUAGE plpgsql;
