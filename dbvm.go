package main

import (
	"fmt"
	"os"

	_ "github.com/ideadawn/dbvm/mysql"
	flags "github.com/jessevdk/go-flags"
)

// 空参数
type optionsEmpty struct{}

// 帮助
type optionsHelp struct {
	Help bool `short:"h" long:"help"`
}

func main() {
	opts := &optionsHelp{}
	args, err := flags.NewParser(opts, flags.PrintErrors|flags.PassDoubleDash|flags.IgnoreUnknown).Parse()
	if err != nil {
		fmt.Println(err)
		fmt.Println("")
		fmt.Print(usage)
		os.Exit(-1)
	}

	if opts.Help {
		fmt.Print(usage)
		return
	}

	if len(args) == 0 {
		fmt.Print(usage)
		os.Exit(-2)
	}

	handler, ok := cmds[args[0]]
	if !ok {
		fmt.Print(usage)
		os.Exit(-3)
	}

	err = handler()
	if err != nil {
		fmt.Println("")
		fmt.Println(err)
		os.Exit(-4)
	}
}
