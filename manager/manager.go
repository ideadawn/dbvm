package manager

import (
	"errors"
)

// ConfParser 配置文件解析函数
type ConfParser func(dir string) (*Config, error)

// PlanParser 部署计划解析函数
type PlanParser func(dir string) (map[string]string, []*Plan, error)

// Manager 管理器，非线程安全
type Manager struct {
	engine   Engine
	params   *Params
	config   *Config
	plans    []*Plan
	deployed map[string]bool

	confParser ConfParser
	planParser PlanParser
}

// GetLogsTable 获取日志数据库的表名
func (m *Manager) GetLogsTable() string {
	if m == nil || m.config == nil {
		return ``
	}
	return m.config.LogsTable
}

// Close 关闭管理器
func (m *Manager) Close() {
	if m.engine != nil {
		m.engine.Close()
		m.engine = nil
	}
}

// New 初始化版本管理器
func New(dir, dburi string) (mgr *Manager, err error) {
	return NewWithParser(dir, dburi, nil, nil)
}

// NewWithParser 自定义配置解析
func NewWithParser(dir, dburi string, confParser ConfParser, planParser PlanParser) (mgr *Manager, err error) {
	mgr = &Manager{
		deployed: make(map[string]bool),

		confParser: confParser,
		planParser: planParser,
	}
	if mgr.confParser == nil {
		mgr.confParser = ParseConfig
	}
	if mgr.planParser == nil {
		mgr.planParser = ParsePlan
	}

	mgr.config, err = mgr.confParser(dir)
	if err != nil {
		return nil, err
	}
	if mgr.config.LogsTable == `` {
		return nil, errors.New(`LogsTable not set.`)
	}

	_, mgr.plans, err = mgr.planParser(dir)
	if err != nil {
		return nil, err
	}

	mgr.params, err = ParseDbUri(dburi)
	if err != nil {
		return nil, err
	}

	mgr.engine = GetEngine(mgr.params.Engine)
	if mgr.engine == nil {
		return nil, errors.New(`Database engine not found.`)
	}

	err = mgr.engine.Connect(mgr.params)
	if err != nil {
		return nil, err
	}

	err = mgr.engine.Initiate(mgr.config.LogsTable)
	if err != nil {
		return nil, err
	}

	return mgr, err
}
