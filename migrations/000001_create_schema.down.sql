-- migrations/000001_create_schema.down.sql
-- Drops all tables and types in the reverse order they were created.

-- Housekeeping Activities
DROP TABLE IF EXISTS maintenance_report;
DROP TABLE IF EXISTS housekeeping_task;

-- Reservations
DROP TABLE IF EXISTS registration;
DROP TABLE IF EXISTS reservation;

-- Room & Room Type
DROP TABLE IF EXISTS room;
DROP TABLE IF EXISTS room_type;

-- Persons
DROP TABLE IF EXISTS housekeeper;
DROP TABLE IF EXISTS front_desk;
DROP TABLE IF EXISTS operations_manager;
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS guest;
DROP TABLE IF EXISTS person;

-- Hotel & Department
DROP TABLE IF EXISTS department;
DROP TABLE IF EXISTS hotel;

-- Types & Extensions
DROP TYPE IF EXISTS housekeeping_task_type;
DROP TYPE IF EXISTS reservation_source;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS shift_type;
DROP TYPE IF EXISTS room_status;
DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pgcrypto;
