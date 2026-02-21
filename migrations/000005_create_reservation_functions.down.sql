-- migrations/000005_create_reservation_functions.down.sql
-- Drops all functions for reservations and registrations in reverse order of their creation.

DROP FUNCTION IF EXISTS fn_create_reservation_workflow(
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

DROP FUNCTION IF EXISTS fn_create_reservation(
    BIGINT,
    DATE,
    DATE,
    payment_method,
    reservation_source,
    INT
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

DROP FUNCTION IF EXISTS fn_get_registrations(
    BIGINT
);

DROP FUNCTION IF EXISTS fn_create_registration(
    BIGINT,
    INT,
    INT
);