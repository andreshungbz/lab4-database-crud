-- migrations/000004_create_guest_functions.down.sql
-- Drops all functions for guests in reverse order of their creation.

DROP FUNCTION IF EXISTS fn_delete_guest(
    TEXT
);

DROP FUNCTION IF EXISTS fn_update_guest(
    TEXT,
    CITEXT,
    TEXT,
    TEXT,
    TEXT,
    TEXT,
    TEXT,
    TEXT
);

DROP FUNCTION IF EXISTS fn_get_guest(
    TEXT
);

DROP FUNCTION IF EXISTS fn_create_guest(
    TEXT, CITEXT, TEXT,
    TEXT, TEXT, TEXT, TEXT, TEXT
);
