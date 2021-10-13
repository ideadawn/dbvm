package main

// CmdHandler 命令调度器
type CmdHandler func() error

var cmds = map[string]CmdHandler{
	`add`:    cmdAdd,
	`deploy`: cmdDeploy,
	`help`:   cmdHelp,
	`init`:   cmdInit,
	`revert`: cmdRevert,
	`verify`: cmdVerify,
}
