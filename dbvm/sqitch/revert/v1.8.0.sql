-- Revert kc:v1.8.0 from mysql

-- IGNORE 1091
BEGIN;

ALTER TABLE `test`
	DROP COLUMN `age`;

COMMIT;
