package mysql

import (
	"errors"
)

var (
	errConnection    = errors.New(`Database was not connected.`)
	errTableName     = errors.New(`Table name should be validated by ^[a-z0-9_]+$`)
	errTableNotInit  = errors.New(`Log table was not initiated.`)
	errDeployNothing = errors.New(`There is nothing to deploy.`)
)
