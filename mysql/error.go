package mysql

import (
	"errors"
)

var (
	errConnection    = errors.New(`Database was not connected.`)
	errTableName     = errors.New(`Table name should be validated by "^[a-z0-9_]+$"`)
	errTableNotInit  = errors.New(`Log table was not initiated.`)
	errDeployNothing = errors.New(`There is nothing to deploy.`)

	errMultiBytes    = errors.New(`Syntax Error: Multi bytes charactor in sql.`)
	errNoIfNotExists = errors.New(`"IF NOT EXISTS" is necessary.`)
	errNoIfExists    = errors.New(`"IF EXISTS" is necessary.`)
	errNoNotNull     = errors.New(`"NOT NULL" is necessary.`)
	errNoDefault     = errors.New(`"DEFAULT" is necessary.`)

	errNoRevert = errors.New(`Any deployment must has a pair of revertion, except "-- NO-REVERT" statement.`)

	errSqlNotEnd = errors.New(`Sql statement without end.`)
)
