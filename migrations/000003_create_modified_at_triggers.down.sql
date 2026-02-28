-- migrations/000003_create_modified_at_triggers.down.sql
-- Drops triggers then trigger functions for setting modified_at.

-- Triggers

DROP TRIGGER IF EXISTS guest_set_person_modified_at ON guest;
DROP TRIGGER IF EXISTS employee_set_person_modified_at ON employee;

DROP TRIGGER IF EXISTS person_set_modified_at ON person;
DROP TRIGGER IF EXISTS room_set_modified_at ON room;
DROP TRIGGER IF EXISTS reservation_set_modified_at ON reservation;
DROP TRIGGER IF EXISTS housekeeping_task_set_modified_at ON housekeeping_task;
DROP TRIGGER IF EXISTS maintenance_report_set_modified_at ON maintenance_report;

-- Trigger Functions

DROP FUNCTION IF EXISTS set_person_modified_at();
DROP FUNCTION IF EXISTS set_modified_at();
