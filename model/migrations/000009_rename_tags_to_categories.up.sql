-- Rename tags table to categories and add hierarchical support
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '#3B82F6',
    description TEXT,
    parent_id INTEGER,
    level INTEGER DEFAULT 0 CHECK (level >= 0 AND level <= 2),
    sort_order INTEGER DEFAULT 0,
    path VARCHAR(255) DEFAULT '/',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE(name, parent_id)
);

-- Migrate existing tags to categories
INSERT INTO categories (id, name, color, description, level, path, created_at, updated_at)
SELECT id, name, color, description, 0, '/', created_at, updated_at
FROM tags;

-- Create index on parent_id
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Rename timelogs.tag_id to category_id
ALTER TABLE timelogs RENAME COLUMN tag_id TO category_id;

-- Update constraints table tag_id to category_id
-- Note: This is handled by application layer, constraints reference is soft

-- Drop old tags table
DROP TABLE IF EXISTS tags;
