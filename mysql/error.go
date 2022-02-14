package mysql

import (
	"errors"
)

var (
	errConnection    = errors.New(`Database was not connected.`)
	errTableName     = errors.New(`Table name should be validated by "^[a-z0-9_]+$"`)
	errTableNotInit  = errors.New(`Log table was not initiated.`)
	errDeployNothing = errors.New(`There is nothing to deploy.`)

	errMultiBytes = errors.New(`Syntax Error: Multi bytes charactor in sql.`)
	errNoNotNull  = errors.New(`"NOT NULL" is necessary.`)
	errNoDefault  = errors.New(`"DEFAULT" is necessary.`)

	errBlockNoBegin = errors.New(`Block must be begin with "BEGIN;".`)
	errBlockNoEnd   = errors.New(`Block must be end with "COMMIT;".`)
	errSqlNoEnd     = errors.New(`Sql statement without end.`)

	errCreateTableINE = errors.New(`"CREATE TABLE" must follow " IF NOT EXISTS".`)
	errDropTableIE    = errors.New(`"DROP TABLE" must follow " IF EXISTS".`)
	errAlterUnknown   = errors.New(`Unknown sub command for alter.`)
)
