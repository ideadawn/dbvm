-- Deploy kc:v1.8.0 to mysql

-- IGNORE 1060
BEGIN;

ALTER TABLE `test`
	ADD COLUMN `age` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '年龄' AFTER `name`;

COMMIT;
