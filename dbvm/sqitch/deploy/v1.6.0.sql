-- Deploy kc:v1.6.0 to mysql

BEGIN;

CREATE TABLE IF NOT EXISTS `test` (
	`id` INT(10) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='测试'
COLLATE='utf8_bin'
ENGINE=InnoDB;

COMMIT;
