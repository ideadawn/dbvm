-- Deploy test:v1.7.0 to mysql


ALTER TABLE `test`
	ADD COLUMN `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '姓名' AFTER `id`,
	ADD COLUMN `age` tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '年龄' AFTER `name`;
