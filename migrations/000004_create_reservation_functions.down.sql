-- migrations/000004_create_reservation_functions.down.sql
-- Drops all functions created for reservations in reverse order of their creation.

DROP FUNCTION IF EXISTS fn_create_reservation_workflow(
    INT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT,
    INT, DATE, DATE, payment_method, reservation_source, BIGINT
);

DROP FUNCTION IF EXISTS fn_create_registration(BIGINT, INT, INT);
DROP FUNCTION IF EXISTS fn_find_available_room(INT, INT, DATE, DATE);
DROP FUNCTION IF EXISTS fn_create_reservation_base(BIGINT, DATE, DATE, INT, payment_method, reservation_source);
DROP FUNCTION IF EXISTS fn_get_or_create_guest(TEXT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT);
