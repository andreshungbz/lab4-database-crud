-- migrations/000003_seed_data.down.sql
-- Delete data and reset sequences.

DELETE FROM maintenance_report;
DELETE FROM housekeeping_task;
DELETE FROM registration;
DELETE FROM reservation;
DELETE FROM housekeeper;
DELETE FROM front_desk;
DELETE FROM operations_manager;
DELETE FROM employee;
DELETE FROM guest;
DELETE FROM person;
DELETE FROM room;
DELETE FROM room_type;
DELETE FROM department;
DELETE FROM hotel;

SELECT setval(pg_get_serial_sequence('hotel', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('person', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('employee', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('room_type', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('reservation', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('housekeeping_task', 'id'), 1, false);
SELECT setval(pg_get_serial_sequence('maintenance_report', 'id'), 1, false);