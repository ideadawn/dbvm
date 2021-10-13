package main

import (
	"askc/tool/dbvm/manager"

	flags "github.com/jessevdk/go-flags"
)

// 回退参数
type optionsRevert struct {
	URI string `long:"uri"`
	Dir string `long:"dir"`
	To  string `long:"to"`
}

// 回退部署
func cmdRevert() error {
	opts := &optionsRevert{}
	_, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	if opts.Dir == `` || opts.URI == `` || opts.To == `` {
		return errArgsNotEnough
	}

	mgr, err := manager.New(opts.Dir, opts.URI)
	if err != nil {
		return err
	}

	return mgr.Revert(opts.To)
}
