-- Seed admin and regular user accounts
-- Passwords: admin@cvmatcher.com / admin123, user@cvmatcher.com / user123

INSERT INTO users (id, first_name, last_name, email, password_hash, created_at, updated_at)
SELECT
    'a0000000-0000-0000-0000-000000000001',
    'Admin',
    'User',
    'admin@cvmatcher.com',
    '$2a$10$DKWAxFLoN1JjjwjLxQ503.z/YkZW58X4yJVWSnxY0GSyeb0JTnOXC',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@cvmatcher.com');

INSERT INTO users (id, first_name, last_name, email, password_hash, created_at, updated_at)
SELECT
    'a0000000-0000-0000-0000-000000000002',
    'Regular',
    'User',
    'user@cvmatcher.com',
    '$2a$10$HIfA2/GDBT0kCisl/7kpG.YDRhMD3Cgmyr55TUcZNCjBRqgvH4i7.',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'user@cvmatcher.com');

-- Assign admin role to admin user
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'admin@cvmatcher.com' AND r.name = 'admin'
ON CONFLICT (user_id) DO UPDATE SET role_id = (SELECT id FROM roles WHERE name = 'admin');

-- Assign user role to regular user
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'user@cvmatcher.com' AND r.name = 'user'
ON CONFLICT (user_id) DO UPDATE SET role_id = (SELECT id FROM roles WHERE name = 'user');

-- Fix existing users that might have been created without proper roles
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'admin@cvmatcher.com' AND r.name = 'admin'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'user@cvmatcher.com' AND r.name = 'user'
ON CONFLICT (user_id) DO NOTHING;
