-- Rename the default "user" role to "professional" for all existing accounts
UPDATE users SET role = 'professional' WHERE role = 'user';
