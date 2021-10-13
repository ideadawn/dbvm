-- Verify kc:v1.8.0 on mysql

BEGIN;

INSERT INTO `test` (`name`, `age`) VALUES ('test', '18');

ROLLBACK;
