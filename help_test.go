package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Help(t *testing.T) {
	os.Args = []string{
		`dbvm`,
		`-h`,
	}
	main()

	os.Args = []string{
		`dbvm`,
		`help`,
		`add`,
	}
	main()

	os.Args = []string{
		`dbvm`,
		`help`,
	}

	err := cmdHelp()
	assert.Equal(t, nil, err)

	os.Args = append(os.Args, `not_exists`)
	err = cmdHelp()
	assert.Equal(t, errCmdNotFound, err)
}
