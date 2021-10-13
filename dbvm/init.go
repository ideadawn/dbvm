package main

import (
	"fmt"
	"regexp"

	"askc/tool/dbvm/manager"

	flags "github.com/jessevdk/go-flags"
)

// 初始化参数
type optionsInit struct {
	Project string   `long:"project"`
	URI     string   `long:"uri"`
	Engine  string   `long:"engine"`
	Dir     string   `long:"dir"`
	Set     []string `long:"set"`
}

// 初始化项目
func cmdInit() error {
	opts := &optionsInit{}
	_, err := flags.Parse(opts)
	if err != nil {
		return err
	}

	reProject := regexp.MustCompile(`^\w+$`)
	if !reProject.Match([]byte(opts.Project)) {
		fmt.Print(usageInit)
		return errArgInvalid
	}

	if opts.Engine == `` || opts.Dir == `` {
		fmt.Print(usageInit)
		return errArgsNotEnough
	}

	project := &manager.ProjectInfo{
		Version: VERSION,
		Project: opts.Project,
		URI:     opts.URI,
		Engine:  opts.Engine,
		Dir:     opts.Dir,
		Set:     opts.Set,
	}
	return manager.InitProject(project)
}
