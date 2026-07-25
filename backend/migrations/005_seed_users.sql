-- Seed admin and regular user accounts
-- Passwords are bcrypt-hashed: "admin123" and "user123"
-- Admin user gets admin role, regular user gets user role

INSERT INTO users (id, first_name, last_name, email, password_hash, role, created_at, updated_at)
SELECT
    'a0000000-0000-0000-0000-000000000001',
    'Admin',
    'User',
    'admin@cvmatcher.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
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
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'user',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'user@cvmatcher.com');
