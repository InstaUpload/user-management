-- Create the roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO roles (name) VALUES ('regular') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('admin') ON CONFLICT (name) DO NOTHING;

-- Add the role_id column to the users table
ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES roles(id) ON DELETE SET DEFAULT DEFAULT 1;

-- Optionally, update existing users to have the 'viewer' role if role_id is NULL
UPDATE users
SET role_id = (SELECT id FROM roles WHERE name = 'regular')
WHERE role_id IS NULL;
