-- migrations/000003_create_modified_at_triggers.up.sql
-- Creates triggers for all tables with modified_at fields.

-- ====================================================================================
-- TRIGGER FUNCTION set_modified_at updates the modified_at field to the current time.
-- ====================================================================================

CREATE OR REPLACE FUNCTION set_modified_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- TRIGGER FUNCTION set_person_modified_at is used for updates to guest or employee.
-- ====================================================================================

CREATE OR REPLACE FUNCTION set_person_modified_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE person
    SET modified_at = NOW()
    WHERE id = NEW.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ====================================================================================
-- TRIGGERS
-- ====================================================================================

CREATE TRIGGER person_set_modified_at
BEFORE UPDATE ON person
FOR EACH ROW
EXECUTE FUNCTION set_modified_at();

CREATE TRIGGER room_set_modified_at
BEFORE UPDATE ON room
FOR EACH ROW
EXECUTE FUNCTION set_modified_at();

CREATE TRIGGER reservation_set_modified_at
BEFORE UPDATE ON reservation
FOR EACH ROW
EXECUTE FUNCTION set_modified_at();

CREATE TRIGGER housekeeping_task_set_modified_at
BEFORE UPDATE ON housekeeping_task
FOR EACH ROW
EXECUTE FUNCTION set_modified_at();

CREATE TRIGGER maintenance_report_set_modified_at
BEFORE UPDATE ON maintenance_report
FOR EACH ROW
EXECUTE FUNCTION set_modified_at();

-- guest updates on person
CREATE TRIGGER guest_set_person_modified_at
AFTER UPDATE ON guest
FOR EACH ROW
EXECUTE FUNCTION set_person_modified_at();

-- employee updates on person
CREATE TRIGGER employee_set_person_modified_at
AFTER UPDATE ON employee
FOR EACH ROW
EXECUTE FUNCTION set_person_modified_at();
