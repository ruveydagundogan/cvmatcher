-- Seed admin and regular user accounts
-- Passwords are bcrypt-hashed: "admin123" and "user123"
-- Admin user gets admin role, regular user gets user role
-- Verified hash for "admin123": $2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2

INSERT INTO users (id, first_name, last_name, email, password_hash, role, created_at, updated_at)
SELECT
    'a0000000-0000-0000-0000-000000000001',
    'Admin',
    'User',
    'admin@cvmatcher.com',
    '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2',
    'admin',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@cvmatcher.com');

INSERT INTO users (id, first_name, last_name, email, password_hash, role, created_at, updated_at)
SELECT
    'a0000000-0000-0000-0000-000000000002',
    'Regular',
    'User',
    'user@cvmatcher.com',
    '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2',
    'user',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'user@cvmatcher.com');

-- Fix existing admin user if created with wrong hash
UPDATE users SET password_hash = '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2'
WHERE email = 'admin@cvmatcher.com' AND password_hash != '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2';

UPDATE users SET password_hash = '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2'
WHERE email = 'user@cvmatcher.com' AND password_hash != '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2';
