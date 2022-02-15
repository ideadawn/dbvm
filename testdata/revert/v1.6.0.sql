-- Revert test:v1.6.0 from mysql


ALTER TABLE `test`
	DROP COLUMN `not_exists`,
	DROP KEY `phone`;

DROP TABLE IF EXISTS `test`;

DROP PROCEDURE IF EXISTS `delete_test`;
