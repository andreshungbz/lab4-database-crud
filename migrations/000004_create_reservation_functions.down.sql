-- migrations/000004_create_reservation_functions.down.sql
-- Drops all functions created for reservations in reverse order of their creation.

DROP FUNCTION IF EXISTS fn_reservation_workflow(
    BIGINT,
    INT,
    DATE,
    DATE,
    payment_method,
    reservation_source,
    INT,
    INT,
    BIGINT
);

DROP FUNCTION IF EXISTS fn_create_registration(
    BIGINT,
    INT,
    INT
);

DROP FUNCTION IF EXISTS fn_create_reservation_base(
    BIGINT,
    DATE,
    DATE,
    payment_method,
    reservation_source,
    INT
);

DROP FUNCTION IF EXISTS fn_get_or_create_guest(
    TEXT, CITEXT, TEXT,
    TEXT, TEXT, TEXT, TEXT, TEXT
);

DROP FUNCTION IF EXISTS fn_calculate_payment(
    INT,
    DATE,
    DATE
);

DROP FUNCTION IF EXISTS fn_find_available_room(
    INT,
    INT,
    DATE,
    DATE
);