-- timelogsh��task_idW��/��sT
ALTER TABLE timelogs ADD COLUMN task_id INTEGER;

-- ��.�_
-- �SQLite/����.�_���\:�c�
-- FOREIGN KEY (task_id) REFERENCES tasks(id)

-- �"����'�
CREATE INDEX idx_timelogs_task_id ON timelogs(task_id);