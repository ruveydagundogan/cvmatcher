-- Fix bcrypt hashes for seed users
-- admin@cvmatcher.com / admin123
UPDATE users SET password_hash = '$2a$10$DKWAxFLoN1JjjwjLxQ503.z/YkZW58X4yJVWSnxY0GSyeb0JTnOXC'
WHERE email = 'admin@cvmatcher.com' AND password_hash != '$2a$10$DKWAxFLoN1JjjwjLxQ503.z/YkZW58X4yJVWSnxY0GSyeb0JTnOXC';

-- user@cvmatcher.com / user123
UPDATE users SET password_hash = '$2a$10$HIfA2/GDBT0kCisl/7kpG.YDRhMD3Cgmyr55TUcZNCjBRqgvH4i7.'
WHERE email = 'user@cvmatcher.com' AND password_hash != '$2a$10$HIfA2/GDBT0kCisl/7kpG.YDRhMD3Cgmyr55TUcZNCjBRqgvH4i7.';
