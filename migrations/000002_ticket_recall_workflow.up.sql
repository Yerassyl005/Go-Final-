ALTER TABLE tickets
ADD COLUMN recall_count INT NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX tickets_one_active_per_queue_idx
ON tickets(queue_id)
WHERE status = 'called';
