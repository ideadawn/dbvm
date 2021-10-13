package manager

import (
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

// 默认配置
const (
	LogsTable = `sqitch_dbyouyou_logs` // 日志表名称

	ConfFile  = `sqitch.conf`
	PlanFile  = `sqitch.plan`
	DeployDir = `deploy`
	VerifyDir = `verify`
	RevertDir = `revert`
)

// Config 数据库配置
type Config struct {
	Engine    string // 优先使用dburi中的engine
	LogsTable string // [default= sqitch_dbyouyou_logs]
}

// 修正目录
func correctDir(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return ``, err
	}
	last := len(dir) - 1
	if dir[last] != '/' && dir[last] != '\\' {
		dir += `/`
	}
	return dir, nil
}

// 解析配置文件
func ParseConfig(dir string) (*Config, error) {
	dir, err := correctDir(dir)
	if err != nil {
		return nil, err
	}

	cnf, err := ini.Load(dir + ConfFile)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		LogsTable: LogsTable,
	}

	sec := cnf.Section(`core`)
	key, _ := sec.GetKey(`engine`)
	if key != nil {
		conf.Engine = key.String()
	}

	sec = cnf.Section(`core "variables"`)
	key, _ = sec.GetKey(`logsTable`)
	if key != nil {
		conf.LogsTable = key.String()
	}

	return conf, nil
}
