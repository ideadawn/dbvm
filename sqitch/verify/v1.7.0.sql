-- Verify kc:v1.7.0 on mysql

BEGIN;

INSERT INTO `test` (`name`) VALUES ('test');

ROLLBACK;
