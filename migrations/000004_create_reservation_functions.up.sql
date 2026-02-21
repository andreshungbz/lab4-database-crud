-- migrations/000004_create_reservation_functions.up.sql
-- Creates functions for creating reservations.

-- ====================================================================================
-- FUNCTION fn_get_or_create_guest returns the id of an existing guest
-- or a newly created one.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_or_create_guest(
    p_passport TEXT,
    p_contact_email TEXT,
    p_contact_phone TEXT,
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT
)
RETURNS BIGINT
AS $$
DECLARE
    v_guest_id BIGINT;
BEGIN
    -- find guest id based on passport number
    SELECT g.id
    INTO v_guest_id
    FROM guest g
    WHERE g.passport_number = p_passport;

    -- return guest id if it already exists
    IF v_guest_id IS NOT NULL THEN
        RETURN v_guest_id;
    END IF;

    -- if guest does not exist, then create one
    BEGIN
        INSERT INTO person (name, gender, street, city, country)
        VALUES (p_name, p_gender, p_street, p_city, p_country)
        RETURNING id INTO v_guest_id;

        INSERT INTO guest (id, passport_number, contact_email, contact_phone)
        VALUES (v_guest_id, p_passport, p_contact_email, p_contact_phone);

        RETURN v_guest_id;

    -- concurrent insert check
    EXCEPTION WHEN unique_violation THEN
        SELECT id INTO v_guest_id
        FROM guest
        WHERE passport_number = p_passport;

        RETURN v_guest_id;
    END;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- FUNCTION fn_calculate_payment takes a room types base rate and multiplies it
-- by the number of nights defined from the checkin and checkout dates.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_calculate_payment(
    p_type_id INT,
    p_checkin DATE,
    p_checkout DATE
)
RETURNS NUMERIC AS $$
DECLARE
    v_days INT;
    v_base_rate NUMERIC(12, 2);
BEGIN
    SELECT base_rate INTO v_base_rate
    FROM room_type
    WHERE id = p_type_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Room type % does not exist', p_type_id;
    END IF;

    -- for same-day checkin and checkout, the minimum would be 1 night
    v_days := GREATEST(1, p_checkout - p_checkin);
    RETURN v_base_rate * v_days;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- FUNCTION fn_create_reservation_base creates the appropriate reservation row.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_reservation_base(
    p_guest_id BIGINT,
    p_checkin DATE,
    p_checkout DATE,
    p_type_id INT,
    p_payment_method payment_method,
    p_source reservation_source
)
RETURNS BIGINT
AS $$
DECLARE
    v_reservation_id BIGINT;
    v_payment_amount NUMERIC(12, 2);
BEGIN
    -- calculate total payment
    v_payment_amount := fn_calculate_payment(p_type_id, p_checkin, p_checkout);

    INSERT INTO reservation (
        guest_id,
        checkin_date,
        checkout_date,
        payment_amount,
        payment_method,
        source
    )
    VALUES (
        p_guest_id,
        p_checkin,
        p_checkout,
        v_payment_amount,
        p_payment_method,
        p_source
    )
    RETURNING id INTO v_reservation_id;

    RETURN v_reservation_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- FUNCTION fn_find_available_room returns the first row that satisfies availability
-- according to the date range. The lowest room number is selected first.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_find_available_room(
    p_hotel_id INT,
    p_type_id INT,
    p_checkin DATE,
    p_checkout DATE
)
RETURNS TABLE (hotel_id INT, room_number INT)
AS $$
BEGIN
    RETURN QUERY
    SELECT r.hotel_id, r.number
    FROM room r
    WHERE r.hotel_id = p_hotel_id
        AND r.room_type_id = p_type_id
        AND r.status_code = 'V/C' -- consider V/C as an available room
        AND NOT EXISTS (
            SELECT 1
            FROM registration reg
            JOIN reservation res ON res.id = reg.reservation_id
            -- also check if a room number is already registered and if checkin/checkout
            -- dates clash
            WHERE reg.hotel_id = r.hotel_id
                AND reg.room_number = r.number
                AND res.checkin_date < p_checkout
                AND res.checkout_date > p_checkin
    )
    ORDER BY r.number
    FOR UPDATE SKIP LOCKED -- concurrency check
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- FUNCTION fn_create_registration creates a registration entry.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_registration(
    p_reservation_id BIGINT,
    p_hotel_id INT,
    p_room_number INT
)
RETURNS VOID
AS $$
BEGIN
    INSERT INTO registration (
        reservation_id,
        hotel_id,
        room_number
    )
    VALUES (
        p_reservation_id,
        p_hotel_id,
        p_room_number
    );
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- FUNCTION create_reservation_workflow is the top-level workflow that creates 
-- reservations and other appropriate entries in the person, guest, and registration
-- tables. For the first room created for a reservation, pass in NULL for p_reservation_id.
-- Since it returns the reservation id, it can be reused in subsequent calls to
-- add to that same reservation.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_reservation_workflow(
    p_hotel_id INT,
    p_passport TEXT,
    p_contact_email TEXT,
    p_contact_phone TEXT,
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT,
    p_type_id INT,
    p_checkin DATE,
    p_checkout DATE,
    p_payment_method payment_method,
    p_source reservation_source,
    p_reservation_id BIGINT DEFAULT NULL
)
RETURNS BIGINT
AS $$
DECLARE
    v_guest_id BIGINT;
    v_reservation_id BIGINT;
    v_room RECORD;
BEGIN
    -- retrieve the id of an existing guest or a newly created one
    v_guest_id := fn_get_or_create_guest(
        p_passport,
        p_contact_email,
        p_contact_phone,
        p_name,
        p_gender,
        p_street,
        p_city,
        p_country
    );

    -- create a new reservation if not provided
    IF p_reservation_id IS NULL THEN
        v_reservation_id := fn_create_reservation_base(
            v_guest_id,
            p_checkin,
            p_checkout,
            p_type_id,
            p_payment_method,
            p_source
        );
    ELSE -- check if the reservation_id passed in belongs to the guest
        SELECT id INTO v_reservation_id
        FROM reservation
        WHERE id = p_reservation_id
            AND guest_id = v_guest_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION
                'Reservation % does not belong to guest %',
                p_reservation_id, v_guest_id;
        END IF;
    END IF;

    -- find available room
    SELECT * INTO v_room
    FROM fn_find_available_room(
        p_hotel_id,
        p_type_id,
        p_checkin,
        p_checkout
    );

    -- return an error in case there are no more available rooms of the requested type
    IF NOT FOUND THEN
        RAISE EXCEPTION 
            '[room-avail-err] No available room of type % for requested dates.',
            p_type_id;
    END IF;

    -- create the registration entry
    -- PERFORM is used since the return value is not needed
    PERFORM fn_create_registration(
        v_reservation_id,
        v_room.hotel_id,
        v_room.room_number
    );

    RETURN v_reservation_id;
END;
$$ LANGUAGE plpgsql;
