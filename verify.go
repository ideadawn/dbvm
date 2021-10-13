package main

import (
	"askc/tool/dbvm/manager"

	flags "github.com/jessevdk/go-flags"
)

// 校验参数
type optionsVerify struct {
	URI  string `long:"uri"`
	Dir  string `long:"dir"`
	From string `long:"from"`
	To   string `long:"to"`
}

// 更新校验
func cmdVerify() error {
	opts := &optionsVerify{}
	_, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	if opts.Dir == `` || opts.URI == `` {
		return errArgsNotEnough
	}
	if opts.From == `` && opts.To == `` {
		return errArgsNotEnough
	}

	_, plans, err := manager.ParsePlan(opts.Dir)
	if err != nil {
		return err
	}

	fromPos, toPos := -1, -1
	for idx, plan := range plans {
		if plan.Name == opts.From {
			fromPos = idx
		}
		if plan.Name == opts.To {
			toPos = idx
		}
	}

	var arr []*manager.Plan
	if opts.From != `` {
		if opts.To != `` {
			if fromPos == -1 || toPos == -1 {
				return errDeployNotFound
			}
			if fromPos > toPos {
				return errFromGtTo
			}
			arr = plans[fromPos : toPos+1]
		} else {
			if fromPos == -1 {
				return errDeployNotFound
			}
			arr = plans[fromPos : fromPos+1]
		}
	} else {
		if toPos == -1 {
			return errDeployNotFound
		}
		arr = plans[toPos : toPos+1]
	}

	mgr, err := manager.New(opts.Dir, opts.URI)
	if err != nil {
		return err
	}

	for _, plan := range arr {
		err = mgr.Verify(plan.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
