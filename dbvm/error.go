package main

import (
	"errors"
)

var (
	errArgsNotEnough = errors.New(`Arguments not enough`)
	errArgInvalid    = errors.New(`Argument validate failed`)

	errCmdNotFound = errors.New(`Command not found`)

	errDirCheckFailed = errors.New(`Dir check failed`)

	errRequireNotFound = errors.New(`Required deployment not found`)
	errDeployNotFound  = errors.New(`Deployment not found`)

	errFromGtTo = errors.New(`--from greater than --to`)
)
