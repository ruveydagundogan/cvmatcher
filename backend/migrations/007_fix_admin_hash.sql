-- Fix bcrypt hashes for seed users
-- Previous hash was incorrect (did not match 'admin123')
-- Verified hash for "admin123": $2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2

UPDATE users SET password_hash = '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2'
WHERE email = 'admin@cvmatcher.com' AND password_hash != '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2';

UPDATE users SET password_hash = '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2'
WHERE email = 'user@cvmatcher.com' AND password_hash != '$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2';
