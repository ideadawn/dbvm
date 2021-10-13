-- Deploy kc:v1.7.0 to mysql

-- IGNORE 1060
BEGIN;

ALTER TABLE `test`
	ADD COLUMN `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '姓名' AFTER `id`;

COMMIT;
