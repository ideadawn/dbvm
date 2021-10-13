-- Revert kc:v1.6.0 from mysql

BEGIN;

DROP TABLE IF EXISTS `test`;

COMMIT;
