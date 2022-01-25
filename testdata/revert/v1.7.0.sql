-- Revert kc:v1.7.0 from mysql


ALTER TABLE `test`
	DROP COLUMN `name`,
	DROP COLUMN `age`;
