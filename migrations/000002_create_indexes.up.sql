-- migrations/000002_create_indexes.up.sql
-- Creates indexes on the tables.
-- PostgreSQL creates indexes automatically for only primary keys and UNIQUE constraints, so
-- indexes are created on the foreign keys of the tables. Additional indexes are created on
-- columns that are likely to be accessed for searches or filtering.

CREATE INDEX idx_employee_hotel ON employee(hotel_id);
CREATE INDEX idx_employee_department ON employee(department);
CREATE INDEX idx_employee_reports_to ON employee(reports_to);

CREATE INDEX idx_guest_email ON guest(contact_email); -- likely to be searched

CREATE INDEX idx_reservation_guest ON reservation(guest_id);
CREATE INDEX idx_reservation_dates ON reservation(checkin_date, checkout_date); -- for availability searches

CREATE INDEX idx_registration_room ON registration(hotel_id, room_number);
CREATE INDEX idx_registration_reservation ON registration(reservation_id);

CREATE INDEX idx_room_type ON room(room_type_id);
CREATE INDEX idx_room_status ON room(status_code); -- for filters by room status

CREATE INDEX idx_housekeeping_housekeeper ON housekeeping_task(housekeeper_id);
CREATE INDEX idx_housekeeping_room ON housekeeping_task(hotel_id, room_number);
CREATE INDEX idx_housekeeping_created ON housekeeping_task(created_at); -- for possible task history ordering

CREATE INDEX idx_maintenance_housekeeper ON maintenance_report(housekeeper_id);
CREATE INDEX idx_maintenance_room ON maintenance_report(hotel_id, room_number);
CREATE INDEX idx_maintenance_created ON maintenance_report(created_at); -- for possible report history ordering
