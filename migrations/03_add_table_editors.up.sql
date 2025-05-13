-- Add editors table which will be used as many to many relation form users to users table
CREATE TABLE IF NOT EXISTS editors (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    editor_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (editor_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Add unique constraint to prevent duplicate entries
CREATE UNIQUE INDEX idx_user_editor ON editors (user_id, editor_id);
