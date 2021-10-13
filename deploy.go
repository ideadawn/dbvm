package main

import (
	"askc/tool/dbvm/manager"

	flags "github.com/jessevdk/go-flags"
)

// 部署参数
type optionsDeploy struct {
	URI string `long:"uri"`
	Dir string `long:"dir"`
	To  string `long:"to"`
}

// 部署更新
func cmdDeploy() error {
	opts := &optionsDeploy{}
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

	return mgr.Deploy(opts.To)
}
