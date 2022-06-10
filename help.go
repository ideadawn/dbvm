package main

import (
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var (
	usage = strings.Join([]string{
		"",
		fmt.Sprintf("%s available commands:", NAME),
		"",
		"  add         Add a new change to the plan.",
		"  deploy      Deploy changes to a database.",
		fmt.Sprintf("  help        Display help information about %s commands.", NAME),
		"  init        Initialize a project.",
		"  print       Print the parsed sql script.",
		"  revert      Revert changes from a database.",
		"",
	}, "\n")

	usageAdd = strings.Join([]string{
		"",
		fmt.Sprintf("%s add [options]", NAME),
		"",
		"  --name      Name's format must be [v]1.0.0[.0]",
		"  --dir       Path to the deployments files.",
		"  --require   Name of change that is required by new change.",
		"  --user      Who make the change.",
		"  --note      A brief note describing the purpose of the change.",
		"",
	}, "\n")

	usageDeploy = strings.Join([]string{
		"",
		fmt.Sprintf("%s deploy [options]", NAME),
		"",
		"  --dir       Path to the deployments files.",
		"  --uri       Database params withing uri formed.",
		"  --to        Upgrade database version to.",
		"",
	}, "\n")

	usageHelp = strings.Join([]string{
		"",
		fmt.Sprintf("%s help [COMMAND]", NAME),
		"",
		"  COMMAND must be one of [add, deploy, help, init, revert, verify].",
		"",
	}, "\n")

	usageInit = strings.Join([]string{
		"",
		fmt.Sprintf("%s init [options]", NAME),
		"",
		"  --project   Project name, must be validated by [a-zA-Z0-9_]+",
		"  --uri       Optional URI to associate with the project.",
		"  --dir       Path to the deployments files.",
		"  --engine    Database driver.",
		"  --table     Log for deploy history.",
		"",
	}, "\n")

	usageRevert = strings.Join([]string{
		"",
		fmt.Sprintf("%s revert [options]", NAME),
		"",
		"  --dir       Path to the deployments files.",
		"  --uri       Database params withing uri formed.",
		"  --to        Revert database version to.",
		"",
	}, "\n")

	usagePrint = strings.Join([]string{
		"",
		fmt.Sprintf("%s print path/to/file.sql", NAME),
		"",
	}, "\n")

	helps = map[string]string{
		`add`:    usageAdd,
		`deploy`: usageDeploy,
		`help`:   usageHelp,
		`init`:   usageInit,
		`print`:  usagePrint,
		`revert`: usageRevert,
	}
)

// 帮助信息
func cmdHelp() error {
	opts := &optionsEmpty{}
	args, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}
	if len(args) < 2 {
		fmt.Print(usage)
		return nil
	}

	info, ok := helps[args[1]]
	if !ok {
		fmt.Print(usageHelp)
		return errCmdNotFound
	}

	fmt.Print(info)
	return nil
}
