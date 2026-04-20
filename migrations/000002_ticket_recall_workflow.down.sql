DROP INDEX IF EXISTS tickets_one_active_per_queue_idx;

ALTER TABLE tickets
DROP COLUMN IF EXISTS recall_count;
