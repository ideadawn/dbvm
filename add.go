package main

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/ideadawn/dbvm/manager"
	flags "github.com/jessevdk/go-flags"
)

// 添加参数
type optionsAdd struct {
	Name    string   `long:"name"`
	Dir     string   `long:"dir"`
	Require []string `long:"require"`
	Note    string   `long:"note"`
	User    string   `long:"user"`
}

// 添加部署计划
func cmdAdd() error {
	opts := &optionsAdd{}
	_, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	if opts.Name == `` || opts.Dir == `` {
		fmt.Print(usageAdd)
		return errArgsNotEnough
	}

	reVersion := regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)
	if !reVersion.Match([]byte(opts.Name)) {
		fmt.Print(usageAdd)
		return errArgInvalid
	}

	for _, req := range opts.Require {
		if !reVersion.Match([]byte(req)) {
			fmt.Print(usageAdd)
			return errArgInvalid
		}
	}

	info, err := os.Stat(opts.Dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errDirCheckFailed
	}

	env, plans, err := manager.ParsePlan(opts.Dir)
	if err != nil {
		return err
	}

	for _, plan := range plans {
		if plan.Name == opts.Name {
			return nil
		}
	}

	checked := 0
	for _, req := range opts.Require {
		for _, plan := range plans {
			if plan.Name == req {
				checked++
				break
			}
		}
	}
	if checked != len(opts.Require) {
		return errRequireNotFound
	}

	plan := &manager.Plan{
		Name:     opts.Name,
		Time:     time.Now(),
		Note:     opts.Note,
		User:     opts.User,
		Requires: opts.Require,
	}

	osu, err := user.Current()
	var osUser string
	if err == nil {
		osUser = osu.Username
	} else {
		osUser = `unkown`
	}
	if plan.User == `` {
		plan.User = osUser
	}
	plan.Hostname, _ = os.Hostname()
	plan.Hostname = strings.Join([]string{
		`<`,
		osUser,
		`@`,
		plan.Hostname,
		`>`,
	}, ``)

	conf, err := manager.ParseConfig(opts.Dir)
	if err != nil {
		return err
	}

	return manager.AddPlan(opts.Dir, plan, env[`project`], conf.Engine)
}
