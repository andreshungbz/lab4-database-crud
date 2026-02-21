-- migrations/000003_seed_data.up.sql
-- Inserts example data to work with.

-- ====================================================================================
-- HOTEL & DEPARTMENT
-- ====================================================================================

INSERT INTO hotel (name, street, city, state, country, phone) VALUES
    ('Grand Ocean View', '1234 Tailwind St', 'San Pedro', 'Belize', 'Belize', '501-111-0001'),
    ('Secret Grove', '5678 Shallow Rd', 'San Ignacio', 'Cayo', 'Belize', '501-222-0002');

INSERT INTO department (dept_name, budget) VALUES
    ('Hotel Operations', 100000),
    ('Guest Services & Front Desk', 50000),
    ('Housekeeping & Maintenance', 30000);

-- ====================================================================================
-- PERSONS
-- ====================================================================================

-- Hotel Operations

INSERT INTO person (name, gender, street, city, country) VALUES
    ('Angus Garcia', 'M', '78 Maple St', 'Belize City', 'Belize');

INSERT INTO employee (id, hotel_id, department, reports_to, salary, ssn, work_email, work_phone, password_hash, employed)
    SELECT id, 1, 'Hotel Operations', NULL, 60000, '123-45-6789', 'angus@grandoceanview.com', '501-111-1001', digest('hotel_password','sha256'), TRUE
    FROM person
    WHERE name='Angus Garcia';

INSERT INTO operations_manager (id, hotel_owner)
    SELECT id, TRUE
    FROM employee
    WHERE work_email='angus@grandoceanview.com';

-- Guest Services & Front Desk

INSERT INTO person (name, gender, street, city, country) VALUES
    ('Bea Sierra', 'F', '90 Cedar St', 'San Pedro', 'Belize');

INSERT INTO employee (id, hotel_id, department, reports_to, salary, ssn, work_email, work_phone, password_hash, employed)
    SELECT id, 1, 'Guest Services & Front Desk', 1, 40000, '987-65-4321', 'bea@grandoceanview.com', '501-222-1002', digest('hotel_password','sha256'), TRUE
    FROM person WHERE name='Bea Sierra';

INSERT INTO front_desk (id, shift)
    SELECT id, 'day'
    FROM employee
    WHERE work_email='bea@grandoceanview.com';

-- Housekeeping & Maintenance

INSERT INTO person (name, gender, street, city, country) VALUES
    ('Clara Mendoza', 'F', '120 Sea Rd', 'San Pedro', 'Belize');

INSERT INTO employee (id, hotel_id, department, reports_to, salary, ssn, work_email, work_phone, password_hash, employed)
    SELECT id, 1, 'Housekeeping & Maintenance', 1, 35000, '555-66-7777', 'clara@grandoceanview.com', '501-111-1003', digest('hotel_password','sha256'), TRUE
    FROM person
    WHERE name='Clara Mendoza';

INSERT INTO housekeeper (id, shift)
    SELECT id, 'day'
    FROM employee
    WHERE work_email='clara@grandoceanview.com';

-- Guests

INSERT INTO person (name, gender, street, city, country) VALUES
    ('Mae Smith', 'F', '12 Oak St', 'Belize City', 'Belize'),
    ('Greg Jones', 'M', '34 Pine St', 'San Pedro', 'Belize'),
    ('Lara Bennett', 'F', '56 Palm St', 'San Pedro', 'Belize'),
    ('Tom Rivera', 'M', '78 Mango St', 'San Pedro', 'Belize'),
    ('Nina Patel', 'F', '90 Oak St', 'San Pedro', 'Belize');

INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    SELECT id, 'A1234567', 'mae@example.com', '501-111-1234'
    FROM person
    WHERE name='Mae Smith';

INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    SELECT id, 'B9876543', 'greg@example.com', '501-222-5678'
    FROM person
    WHERE name='Greg Jones';

INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    SELECT id, 'C1122334', 'lara@example.com', '501-111-2345'
    FROM person
    WHERE name='Lara Bennett';

INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    SELECT id, 'D5566778', 'tom@example.com', '501-111-3456'
    FROM person
    WHERE name='Tom Rivera';

INSERT INTO guest (id, passport_number, contact_email, contact_phone)
    SELECT id, 'E9988776', 'nina@example.com', '501-111-4567'
    FROM person
    WHERE name='Nina Patel';

-- ====================================================================================
-- ROOM TYPE & ROOM
-- ====================================================================================

INSERT INTO room_type (title, base_rate, max_occupancy, bed_count, has_balcony) VALUES
    ('Single', 100, 1, 1, FALSE),
    ('Double', 150, 2, 2, FALSE),
    ('Suite', 300, 4, 2, TRUE);

INSERT INTO room (hotel_id, number, room_type_id, floor, status_code) VALUES
    (1, 101, 1, 1, 'O/C'),
    (1, 102, 2, 1, 'O/C'),
    (1, 201, 3, 2, 'O/C'),
    (1, 202, 1, 2, 'O/D'),
    (1, 301, 2, 3, 'V/D'),
    (1, 302, 3, 3, 'V/C'),
    (1, 401, 1, 4, 'V/C'),
    (1, 402, 2, 4, 'V/C');

-- ====================================================================================
-- RESERVATIONS
-- ====================================================================================

-- Room 101: O/C (single-room reservation)

INSERT INTO reservation (guest_id, checkin_date, checkout_date, payment_amount, payment_method, source, canceled)
    SELECT id, '2026-03-01', '2026-03-05', 600, 'cash', 'direct', FALSE
    FROM guest
    WHERE passport_number='A1234567';

INSERT INTO registration (reservation_id, hotel_id, room_number)
    SELECT r.id, 1, 101
    FROM reservation r
    JOIN guest g ON r.guest_id = g.id
    WHERE g.passport_number='A1234567';

-- Room 102 & 201: O/C (multi-room reservation)

INSERT INTO reservation (guest_id, checkin_date, checkout_date, payment_amount, payment_method, source, canceled)
    SELECT id, '2026-03-02', '2026-03-06', 300, 'credit_card', 'direct', FALSE
    FROM guest
    WHERE passport_number='C1122334';

INSERT INTO registration (reservation_id, hotel_id, room_number)
    SELECT r.id, 1, 102
    FROM reservation r
    JOIN guest g ON r.guest_id = g.id
    WHERE g.passport_number='C1122334';

INSERT INTO registration (reservation_id, hotel_id, room_number)
    SELECT r.id, 1, 201
    FROM reservation r
    JOIN guest g ON r.guest_id = g.id
    WHERE g.passport_number='C1122334';

-- Room 202: O/D

INSERT INTO reservation (guest_id, checkin_date, checkout_date, payment_amount, payment_method, source, canceled)
    SELECT id, '2026-03-04', '2026-03-08', 100, 'debit_card', 'Expedia', FALSE
    FROM guest
    WHERE passport_number='D5566778';

INSERT INTO registration (reservation_id, hotel_id, room_number)
    SELECT r.id, 1, 202
    FROM reservation r
    JOIN guest g ON r.guest_id = g.id
    WHERE g.passport_number='D5566778';

-- Room 301: V/D

INSERT INTO reservation (guest_id, checkin_date, checkout_date, payment_amount, payment_method, source, canceled)
    SELECT id, '2026-03-12', '2026-03-13', 250, 'debit_card', 'Booking.com', FALSE
    FROM guest
    WHERE passport_number='E9988776';

INSERT INTO registration (reservation_id, hotel_id, room_number)
    SELECT r.id, 1, 301
    FROM reservation r
    JOIN guest g ON r.guest_id = g.id
    WHERE g.passport_number='E9988776';

-- ====================================================================================
-- HOUSEKEEPING ACTIVITIES
-- ====================================================================================

INSERT INTO housekeeping_task (hotel_id, room_number, housekeeper_id, task_type, completed_at) VALUES
    (1, 101, 3, 'bed', '2026-02-18 09:00:00'),
    (1, 102, 3, 'bathroom', NULL),
    (1, 201, 3, 'bed', NULL),
    (1, 202, 3, 'bathroom', '2026-02-19 14:30:00'),
    (1, 101, 3, 'dusting', '2026-02-18 15:45:00');

INSERT INTO maintenance_report (hotel_id, room_number, housekeeper_id, description, completed_at) VALUES
    (1, 201, 3, 'Air conditioning not working', NULL),
    (1, 102, 3, 'Leaky faucet', NULL),
    (1, 101, 3, 'Broken light fixture', '2026-02-17 11:20:00'),
    (1, 102, 3, 'Window lock broken', NULL);
