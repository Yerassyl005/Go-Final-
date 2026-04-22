ALTER TABLE users
ADD COLUMN priority_category TEXT NOT NULL DEFAULT 'none'
    CHECK (priority_category IN ('none', 'pregnant', 'elderly', 'disabled'));

ALTER TABLE tickets
ADD COLUMN is_priority BOOLEAN NOT NULL DEFAULT false;
