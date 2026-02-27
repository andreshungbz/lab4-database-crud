-- migrations/000002_create_indexes.down.sql
-- Drops created indexes.

DROP INDEX IF EXISTS idx_employee_hotel;
DROP INDEX IF EXISTS idx_employee_department;
DROP INDEX IF EXISTS idx_employee_reports_to;

DROP INDEX IF EXISTS idx_guest_email;

DROP INDEX IF EXISTS idx_reservation_guest;
DROP INDEX IF EXISTS idx_reservation_dates;

DROP INDEX IF EXISTS idx_registration_room;
DROP INDEX IF EXISTS idx_registration_reservation;

DROP INDEX IF EXISTS idx_room_type;
DROP INDEX IF EXISTS idx_room_status;

DROP INDEX IF EXISTS idx_housekeeping_housekeeper;
DROP INDEX IF EXISTS idx_housekeeping_room;
DROP INDEX IF EXISTS idx_housekeeping_created;

DROP INDEX IF EXISTS idx_maintenance_housekeeper;
DROP INDEX IF EXISTS idx_maintenance_room;
DROP INDEX IF EXISTS idx_maintenance_created;

DROP INDEX IF EXISTS idx_person_name;
DROP INDEX IF EXISTS idx_person_country;
