-- Revert kc:v1.7.0 from mysql

-- IGNORE 1091
BEGIN;

ALTER TABLE `test`
	DROP COLUMN `name`;

COMMIT;
