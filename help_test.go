package main

import (
	"os"
	"testing"

	"github.com/nbio/st"
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
	st.Assert(t, err, nil)

	os.Args = append(os.Args, `not_exists`)
	err = cmdHelp()
	st.Assert(t, err, errCmdNotFound)
}
