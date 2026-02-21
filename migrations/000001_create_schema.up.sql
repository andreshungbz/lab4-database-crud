-- migrations/000001_create_schema.up.sql
-- Creates the entire database schema.

-- ====================================================================================
-- EXTENSIONS & TYPES
-- ====================================================================================

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE room_status AS ENUM ('V/C', 'O/C', 'O/D', 'V/D');
CREATE TYPE shift_type AS ENUM ('day', 'night');
CREATE TYPE payment_method AS ENUM ('cash', 'debit_card', 'credit_card');
CREATE TYPE reservation_source AS ENUM ('direct', 'Expedia', 'Booking.com');
CREATE TYPE housekeeping_task_type AS ENUM ('bed', 'bathroom', 'dusting');

-- ====================================================================================
-- HOTEL & DEPARTMENT
-- ====================================================================================

CREATE TABLE hotel (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    street TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL,
    phone TEXT NOT NULL
);

CREATE TABLE department (
    dept_name TEXT PRIMARY KEY,
    budget NUMERIC(18, 2)
);

-- ====================================================================================
-- PERSONS
-- ====================================================================================

CREATE TABLE person (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE guest (
    id BIGINT PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
    passport_number TEXT UNIQUE NOT NULL,
    contact_email CITEXT NOT NULL,
    contact_phone TEXT NOT NULL
);

CREATE TABLE employee (
    id BIGINT PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
    hotel_id BIGINT NOT NULL REFERENCES hotel(id) ON DELETE CASCADE,
    department TEXT NOT NULL REFERENCES department(dept_name),
    reports_to BIGINT REFERENCES employee(id),
    salary NUMERIC(18, 2) NOT NULL,
    ssn TEXT UNIQUE NOT NULL,
    work_email CITEXT UNIQUE NOT NULL,
    work_phone TEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    employed BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE operations_manager (
    id BIGINT PRIMARY KEY REFERENCES employee(id) ON DELETE CASCADE,
    hotel_owner BOOLEAN NOT NULL
);

CREATE TABLE front_desk (
    id BIGINT PRIMARY KEY REFERENCES employee(id) ON DELETE CASCADE,
    shift shift_type NOT NULL
);

CREATE TABLE housekeeper (
    id BIGINT PRIMARY KEY REFERENCES employee(id) ON DELETE CASCADE,
    shift shift_type NOT NULL
);

-- ====================================================================================
-- ROOM TYPE & ROOM
-- ====================================================================================

CREATE TABLE room_type (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    base_rate NUMERIC(12, 2) NOT NULL,
    max_occupancy INT NOT NULL,
    bed_count INT NOT NULL,
    has_balcony BOOLEAN NOT NULL
);

CREATE TABLE room (
    hotel_id INT REFERENCES hotel(id) ON DELETE CASCADE,
    number INT,
    room_type_id INT NOT NULL REFERENCES room_type(id) ON DELETE CASCADE,
    floor INT NOT NULL,
    status_code room_status NOT NULL DEFAULT 'V/C',
    PRIMARY KEY (hotel_id, number)
);

-- ====================================================================================
-- RESERVATIONS
-- ====================================================================================

CREATE TABLE reservation (
    id BIGSERIAL PRIMARY KEY,
    guest_id BIGINT NOT NULL REFERENCES guest(id) ON DELETE CASCADE,
    checkin_date DATE NOT NULL,
    checkout_date DATE NOT NULL CHECK (checkout_date > checkin_date),
    payment_amount NUMERIC(12, 2) NOT NULL,
    payment_method payment_method NOT NULL,
    source reservation_source NOT NULL,
    canceled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE registration (
    reservation_id BIGINT REFERENCES reservation(id) ON DELETE CASCADE,
    hotel_id BIGINT NOT NULL,
    room_number INT NOT NULL,
    PRIMARY KEY (reservation_id, hotel_id, room_number),
    FOREIGN KEY (hotel_id, room_number)
        REFERENCES room(hotel_id, number)
        ON DELETE CASCADE
);

-- ====================================================================================
-- HOUSEKEEPING ACTIVITIES
-- ====================================================================================

CREATE TABLE housekeeping_task (
    id BIGSERIAL PRIMARY KEY,
    hotel_id BIGINT NOT NULL,
    room_number INT NOT NULL,
    housekeeper_id BIGINT REFERENCES housekeeper(id),
    task_type housekeeping_task_type NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP(0) WITH TIME ZONE,
    FOREIGN KEY (hotel_id, room_number)
        REFERENCES room(hotel_id, number)
        ON DELETE CASCADE
);

CREATE TABLE maintenance_report (
    id BIGSERIAL PRIMARY KEY,
    hotel_id BIGINT NOT NULL,
    room_number INT NOT NULL,
    housekeeper_id BIGINT NOT NULL REFERENCES housekeeper(id),
    description TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP(0) WITH TIME ZONE,
    FOREIGN KEY (hotel_id, room_number)
        REFERENCES room(hotel_id, number)
        ON DELETE CASCADE
);
