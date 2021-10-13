package manager

// Log 数据库操作日志
type Log struct {
	ID     int64  // 日志ID， 可选
	Name   string // 版本名称，必须
	Time   int64  // 部署时间，可选
	Status int8   // 完成状态，必须（0未完成，1已完成）
}
