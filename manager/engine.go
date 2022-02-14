package manager

import (
	"fmt"
)

// Engine 数据库驱动引擎
type Engine interface {
	// Connect 连接到数据库
	Connect(*Params) error
	// Close 关闭数据库连接
	Close()

	// Initiate 初始化日志表
	Initiate(string) error

	// ListLogs 获取版本更新历史
	ListLogs() ([]*Log, error)

	// 不提供单独保存日志和删除日志的接口，而是在部署及回退时一起执行

	// Deploy 部署指定版本
	Deploy(*Plan) error
	// Revert 回退指定版本
	Revert(*Plan) error
}

var engines = map[string]Engine{}

// RegisterEngine 注册数据库引擎
func RegisterEngine(name string, eng Engine) {
	_, ok := engines[name]
	if ok {
		fmt.Println(`数据库引擎被替换：`, name)
	}
	engines[name] = eng
}

// GetEngine 获取数据库引擎，可能为空
func GetEngine(name string) Engine {
	return engines[name]
}
