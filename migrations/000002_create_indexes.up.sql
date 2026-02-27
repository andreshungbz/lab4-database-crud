-- migrations/000002_create_indexes.up.sql
-- Creates indexes on the tables.
-- PostgreSQL creates indexes automatically for only primary keys and UNIQUE constraints, so
-- indexes are created on the foreign keys of the tables. Additional indexes are created on
-- columns that are likely to be accessed for searches or filtering.

CREATE INDEX IF NOT EXISTS idx_employee_hotel ON employee(hotel_id);
CREATE INDEX IF NOT EXISTS idx_employee_department ON employee(department);
CREATE INDEX IF NOT EXISTS idx_employee_reports_to ON employee(reports_to);

CREATE INDEX IF NOT EXISTS idx_guest_email ON guest(contact_email); -- likely to be searched

CREATE INDEX IF NOT EXISTS idx_reservation_guest ON reservation(guest_id);
CREATE INDEX IF NOT EXISTS idx_reservation_dates ON reservation(checkin_date, checkout_date); -- for availability searches

CREATE INDEX IF NOT EXISTS idx_registration_room ON registration(hotel_id, room_number);
CREATE INDEX IF NOT EXISTS idx_registration_reservation ON registration(reservation_id);

CREATE INDEX IF NOT EXISTS idx_room_type ON room(room_type_id);
CREATE INDEX IF NOT EXISTS idx_room_status ON room(status_code); -- for filters by room status

CREATE INDEX IF NOT EXISTS idx_housekeeping_housekeeper ON housekeeping_task(housekeeper_id);
CREATE INDEX IF NOT EXISTS idx_housekeeping_room ON housekeeping_task(hotel_id, room_number);
CREATE INDEX IF NOT EXISTS idx_housekeeping_created ON housekeeping_task(created_at); -- for possible task history ordering

CREATE INDEX IF NOT EXISTS idx_maintenance_housekeeper ON maintenance_report(housekeeper_id);
CREATE INDEX IF NOT EXISTS idx_maintenance_room ON maintenance_report(hotel_id, room_number);
CREATE INDEX IF NOT EXISTS idx_maintenance_created ON maintenance_report(created_at); -- for possible report history ordering

-- GIN Indexes for PostgreSQL full text-search
CREATE INDEX IF NOT EXISTS idx_person_name ON person USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS idx_person_country ON person USING GIN (to_tsvector('simple', country));
