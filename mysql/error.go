package mysql

import (
	"errors"
)

var (
	errConnection    = errors.New(`Database was not connected.`)
	errTableName     = errors.New(`Table name should be validated by "^[a-z0-9_]+$"`)
	errTableNotInit  = errors.New(`Log table was not initiated.`)
	errDeployNothing = errors.New(`There is nothing to deploy.`)

	errSyntaxError = errors.New(`Syntax error, Multi bytes charactor in sql.`)
	// errNoNotNull  = errors.New(`"NOT NULL" is necessary.`)
	// errNoDefault  = errors.New(`"DEFAULT" is necessary.`)

	errSqlNoEnd = errors.New(`Sql statement without end.`)

	errCreateTableINE = errors.New(`"CREATE TABLE" must follow " IF NOT EXISTS".`)
	errDropTableIE    = errors.New(`"DROP TABLE" must follow " IF EXISTS".`)
	errAlterUnknown   = errors.New(`Unknown sub command for alter.`)
)
