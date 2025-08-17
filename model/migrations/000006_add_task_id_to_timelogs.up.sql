-- timelogshû task_idWµå/û¡sT
ALTER TABLE timelogs ADD COLUMN task_id INTEGER;

-- û .¦_
-- èSQLite/ô¥û .¦_ÙÌÅ\:‡cô
-- FOREIGN KEY (task_id) REFERENCES tasks(id)

-- ú"ÐØåâ'ý
CREATE INDEX idx_timelogs_task_id ON timelogs(task_id);