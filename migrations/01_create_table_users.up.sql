-- Add SQL query to create a users table with ID name email created_at and is_verified column, where id will be primary key.
CREATE EXTENSION citext;
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email CITEXT UNIQUE NOT NULL,
	password bytea NOT NULL,
    created_on TIMESTAMP DEFAULT NOW(),
    is_verified BOOLEAN DEFAULT FALSE
);
