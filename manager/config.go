package manager

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// 默认配置
const (
	LogsTable = `dbvm_logs` // 日志表名称

	ConfFile  = `dbvm.yaml`
	PlanFile  = `dbvm.plan`
	DeployDir = `deploy`
	RevertDir = `revert`
)

// 魔法注释
var (
	NoRevert      = []byte(`NO-REVERT`)      // 不需要回退
	NoTransaction = []byte(`NO-TRANSACTION`) // 不需要启动事务
)

// Database 数据库规则
type Database struct {
	Create bool `yaml:"create"`
	Drop   bool `yaml:"drop"`
}

// Field 字段规则
type Field struct {
	NotNull bool     `yaml:"notNull"`
	Excepts []string `yaml:"excepts"`
}

// Rule 规则
type Rule struct {
	Database *Database `yaml:"database"`
	Field    *Field    `yaml:"field"`

	NoRevert      bool `yaml:"-"`
	NoTransaction bool `yaml:"-"`
}

// Config 数据库配置
type Config struct {
	Engine    string `yaml:"engine"`    // 优先使用dburi中的engine
	FromTable string `yaml:"fromTable"` // 重命名日志表
	LogsTable string `yaml:"logsTable"` // [default= dbvm_logs]

	Rule *Rule `yaml:"rule"`
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

	data, err := os.ReadFile(dir + ConfFile)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	if conf.Rule == nil {
		conf.Rule = &Rule{}
	}
	if conf.Rule.Database == nil {
		conf.Rule.Database = &Database{
			Create: false,
			Drop:   false,
		}
	}
	if conf.Rule.Field == nil {
		conf.Rule.Field = &Field{
			NotNull: true,
			Excepts: []string{
				"TEXT",
				"BLOB",
			},
		}
	}

	return conf, err
}
