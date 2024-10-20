ALTER TABLE enrollment_results
DROP COLUMN capacity,
ADD COLUMN result BOOL NOT NULL;