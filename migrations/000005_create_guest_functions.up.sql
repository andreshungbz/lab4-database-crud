-- migrations/000005_create_guest_functions.up.sql
-- Creates functions for guests.

-- ====================================================================================
-- CREATE FUNCTION fn_create_guest returns the person id and created_at for a newly created
-- guest from on the passed guest and person details.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_guest(
    -- guest attributes
    p_passport TEXT,
    p_contact_email CITEXT,
    p_contact_phone TEXT,
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT
)
RETURNS TABLE (
    id BIGINT,
    created_at TIMESTAMP(0) WITH TIME ZONE
)
AS $$
DECLARE
    v_guest_id BIGINT;
BEGIN
    -- insert person entry
    INSERT INTO person (name, gender, street, city, country)
    VALUES (p_name, p_gender, p_street, p_city, p_country)
    RETURNING person.id INTO v_guest_id;

    -- insert guest entry
    INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    VALUES (v_guest_id, p_passport, p_contact_email, p_contact_phone);

    RETURN QUERY
    SELECT
        p.id,
        p.created_at
    FROM person p
    WHERE p.id = v_guest_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- READ FUNCTION fn_get_guest_by_passport returns the guest and person data for an existing
-- guest by their passport.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_guest(
    p_passport TEXT
)
RETURNS TABLE (
    id BIGINT,
    passport_number TEXT,
    contact_email CITEXT,
    contact_phone TEXT,
    name TEXT,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE
)
AS $$
BEGIN
    RETURN QUERY
    SELECT
        g.id,
        g.passport_number,
        g.contact_email,
        g.contact_phone,
        p.name,
        p.gender,
        p.street,
        p.city,
        p.country,
        p.created_at,
        p.modified_at
    FROM guest g
    JOIN person p ON p.id = g.id
    WHERE g.passport_number = p_passport;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- UPDATE FUNCTION fn_update_guest updates person and guest details
-- based on passport number.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_update_guest(
    p_passport TEXT,
    -- guest attributes
    p_contact_email CITEXT,
    p_contact_phone TEXT,
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT
)
RETURNS VOID
AS $$
DECLARE
    v_guest_id BIGINT;
BEGIN
    -- find guest id from passport
    SELECT g.id
    INTO v_guest_id
    FROM guest g
    WHERE g.passport_number = p_passport;

    IF NOT FOUND THEN
        RAISE EXCEPTION
            '[guest-not-found] Guest with passport % does not exist',
            p_passport;
    END IF;

    -- update guest
    UPDATE guest
    SET
        contact_email = p_contact_email,
        contact_phone = p_contact_phone
    WHERE id = v_guest_id;

    -- update person
    UPDATE person
    SET
        name = p_name,
        gender = p_gender,
        street = p_street,
        city = p_city,
        country = p_country
    WHERE id = v_guest_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- DELETE FUNCTION fn_delete_guest deletes a guest (via person table)
-- based on passport number. Associated reservation and registration records are also
-- deleted.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_delete_guest(
    p_passport TEXT
)
RETURNS VOID
AS $$
DECLARE
    v_guest_id BIGINT;
BEGIN
    -- find guest id
    SELECT id
    INTO v_guest_id
    FROM guest
    WHERE passport_number = p_passport;

    IF NOT FOUND THEN
        RAISE EXCEPTION
            '[guest-not-found] Guest with passport % does not exist',
            p_passport;
    END IF;

    -- delete from person (cascades to guest, reservation, and registration)
    DELETE FROM person
    WHERE id = v_guest_id;

END;
$$ LANGUAGE plpgsql;
