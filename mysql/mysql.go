package mysql

import (
	"database/sql"

	"askc/tool/dbvm/manager"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	manager.RegisterEngine(`mysql`, New())
}

var retry int = 3

// 部署状态
const (
	StatusDeploying int8 = 0
	StatusDeployed  int8 = 1
	StatusVerified  int8 = 2
)

// MySQL 驱动引擎
type MySQL struct {
	db    *sql.DB
	table string
}

// New 创建一个新的引擎
func New() *MySQL {
	return &MySQL{}
}

// Connect 连接数据库
func (m *MySQL) Connect(params *manager.Params) (err error) {
	m.db, err = sql.Open(`mysql`, manager.DbUri2Dsn(params))
	return
}

// Close 关闭连接
func (m *MySQL) Close() {
	if m.db != nil {
		_ = m.db.Close()
		m.db = nil
	}
}
