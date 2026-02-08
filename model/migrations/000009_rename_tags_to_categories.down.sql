-- Revert categories back to tags
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#3B82F6',
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Migrate categories back to tags (only root level)
INSERT INTO tags (id, name, color, description, created_at, updated_at)
SELECT id, name, color, description, created_at, updated_at
FROM categories
WHERE level = 0;

-- Rename timelogs.category_id back to tag_id
ALTER TABLE timelogs RENAME COLUMN category_id TO tag_id;

-- Drop categories table
DROP TABLE IF EXISTS categories;
