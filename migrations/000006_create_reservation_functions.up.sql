-- migrations/000006_create_reservation_functions.up.sql
-- Creates functions for reservations and registrations.

-- ====================================================================================
-- CREATE FUNCTION fn_create_registration creates a registration entry.
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
-- READ FUNCTION fn_get_registrations returns all registrations for a
-- given reservation id with full room and room type details.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_registrations(
    p_reservation_id BIGINT
)
RETURNS TABLE (
    hotel_id BIGINT,
    room_number INT,
    floor INT,
    status_code room_status,
    modified_at TIMESTAMP(0) WITH TIME ZONE,
    room_type_id INT,
    title TEXT,
    base_rate NUMERIC(12, 2),
    max_occupancy INT,
    bed_count INT,
    has_balcony BOOLEAN
)
AS $$
BEGIN
    RETURN QUERY
    SELECT
        reg.hotel_id,
        reg.room_number,
        r.floor,
        r.status_code,
        r.modified_at,
        rt.id,
        rt.title,
        rt.base_rate,
        rt.max_occupancy,
        rt.bed_count,
        rt.has_balcony
    FROM registration reg
    JOIN room r
        ON r.hotel_id = reg.hotel_id
        AND r.number = reg.room_number
    JOIN room_type rt
        ON rt.id = r.room_type_id
    WHERE reg.reservation_id = p_reservation_id
    ORDER BY reg.hotel_id, reg.room_number;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- HELPER FUNCTION fn_find_available_room returns the first row that satisfies availability
-- according to the date range. The lowest room number is selected first.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_find_available_room(
    p_hotel_id INT,
    p_room_type_id INT,
    p_checkin DATE,
    p_checkout DATE
)
RETURNS TABLE (
    hotel_id INT,
    room_number INT
)
AS $$
BEGIN
    RETURN QUERY
    SELECT r.hotel_id, r.number
    FROM room r
    WHERE r.hotel_id = p_hotel_id
        AND r.room_type_id = p_room_type_id
        AND r.status_code = 'V/C' -- only consider vacant/clean rooms
        AND NOT EXISTS (
            SELECT 1
            FROM registration reg
            JOIN reservation res ON res.id = reg.reservation_id
            WHERE reg.hotel_id = r.hotel_id
              AND reg.room_number = r.number
              AND res.canceled = FALSE -- ignore canceled reservations
              AND res.checkout_date > p_checkin -- overlapping check
              AND res.checkin_date < p_checkout
        )
    ORDER BY r.number
    FOR UPDATE SKIP LOCKED  -- safe for concurrent allocations
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- HELPER FUNCTION fn_calculate_payment takes a room types base rate and multiplies it
-- by the number of nights defined from the checkin and checkout dates.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_calculate_payment(
    -- room attribute
    p_room_type_id INT,
    -- reservation attributes
    p_checkin DATE,
    p_checkout DATE
)
RETURNS NUMERIC AS $$
DECLARE
    v_days INT;
    v_base_rate NUMERIC(12, 2);
BEGIN
    -- get room type's base rate
    SELECT base_rate INTO v_base_rate
    FROM room_type
    WHERE id = p_room_type_id;

    -- handle non-existent room types
    IF NOT FOUND THEN
        RAISE EXCEPTION '[nonexistent-room-type] Room type % does not exist', p_room_type_id;
    END IF;

    -- calculate days
    -- checkout_date > checkin_date business rule in schema ensures this is at least 1
    v_days := p_checkout - p_checkin;
    RETURN v_base_rate * v_days;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- CREATE FUNCTION fn_create_reservation creates a reservation entry.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_reservation(
    -- reservation attributes
    p_guest_id BIGINT,
    p_checkin DATE,
    p_checkout DATE,
    p_payment_method payment_method,
    p_source reservation_source,
    -- room attribute
    p_room_type_id INT
)
RETURNS BIGINT
AS $$
DECLARE
    v_reservation_id BIGINT;
    v_payment_amount NUMERIC(12, 2);
BEGIN
    -- calculate total payment
    v_payment_amount := fn_calculate_payment(p_room_type_id, p_checkin, p_checkout);

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
-- CREATE FUNCTION fn_create_reservation_workflow creates a reservation and its associated
-- registration. It also adds a registration to an existing reservation if passed
-- an existing reservation id.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_reservation_workflow(
    -- reservation attributes
    p_guest_id BIGINT,
    p_checkin DATE,
    p_checkout DATE,
    p_payment_method payment_method,
    p_source reservation_source,
    -- registration attribute
    p_hotel_id INT,
    -- room attribute
    p_room_type_id INT,
    -- for multiple registrations to one reservation
    p_reservation_id BIGINT DEFAULT NULL
)
RETURNS BIGINT
AS $$
DECLARE
    v_reservation_id BIGINT;
    v_room RECORD;
BEGIN
    -- create a new reservation if p_reservation_id was not provided
    IF p_reservation_id IS NULL THEN
        v_reservation_id := fn_create_reservation(
            p_guest_id,
            p_checkin,
            p_checkout,
            p_payment_method,
            p_source,
            p_room_type_id
        );
    ELSE -- if p_reservation_id was provided
        -- check if the reservation id passed corresponds to the guest id passed
        SELECT id INTO v_reservation_id
        FROM reservation
        WHERE id = p_reservation_id
            AND guest_id = p_guest_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION
                '[reservation-guest-mismatch] Reservation % does not belong to guest %',
                p_reservation_id, p_guest_id;
        END IF;
    END IF;

    -- get available room of the specified room type that does not have conflicting dates
    SELECT * INTO v_room
    FROM fn_find_available_room(
        p_hotel_id,
        p_room_type_id,
        p_checkin,
        p_checkout
    );

    -- return an error if there are no available rooms
    IF NOT FOUND THEN
        RAISE EXCEPTION 
            '[no-available-room] No available room of type % for requested dates',
            p_room_type_id;
    END IF;

    -- create the registration entry (PERFORM is used since the return value is not needed)
    PERFORM fn_create_registration(
        v_reservation_id,
        v_room.hotel_id,
        v_room.room_number
    );

    RETURN v_reservation_id;
END;
$$ LANGUAGE plpgsql;
