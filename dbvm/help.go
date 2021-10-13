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
		"  revert      Revert changes from a database.",
		"  verify      Verify changes to a database.",
		"",
	}, "\n")

	usageAdd = strings.Join([]string{
		"",
		fmt.Sprintf("%s add [options]", NAME),
		"",
		"  --name      Name's format must be [v]1.0.0",
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
		"  --project   Project's format must be \\w+",
		"  --uri       Optional URI to associate with the project.",
		"  --engine    Database driver.",
		"  --dir       Path to the deployments files.",
		"  --set       Set a variable name and value.",
		"              The format must be \"name=value\".",
		"              Variables are set in \"core.variables\".",
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

	usageVerify = strings.Join([]string{
		"",
		fmt.Sprintf("%s verify [options]", NAME),
		"",
		"  --dir       Path to the deployments files.",
		"  --uri       Database params withing uri formed.",
		"  --from      Verify deployment from.",
		"  --to        Verify deployment to.",
		"",
	}, "\n")

	helps = map[string]string{
		`add`:    usageAdd,
		`deploy`: usageDeploy,
		`help`:   usageHelp,
		`init`:   usageInit,
		`revert`: usageRevert,
		`verify`: usageVerify,
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
