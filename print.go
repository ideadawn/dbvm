package main

import (
	"github.com/ideadawn/dbvm/mysql"
	flags "github.com/jessevdk/go-flags"
)

// 更新校验
func cmdPrint() error {
	args, err := flags.NewParser(nil, flags.Default|flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return errArgsNotEnough
	}

	return mysql.Print(args[1])
}
