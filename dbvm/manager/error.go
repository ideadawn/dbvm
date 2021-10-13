package manager

import (
	"errors"
)

var (
	errProjectExists  = errors.New(`Project has exists.`)
	errProjectNotInit = errors.New(`Project not initiated.`)
)
