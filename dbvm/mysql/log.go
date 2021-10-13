package mysql

import (
	"database/sql"
	"fmt"
	"askc/tool/dbvm/manager"
)

// ListLogs 获取版本更新历史
func (m *MySQL) ListLogs() (logs []*manager.Log, err error) {
	if m == nil || m.db == nil {
		return nil, errConnection
	}
	if m.table == `` {
		return nil, errTableNotInit
	}

	var lastID int64
	limit := 1000
	query := "SELECT `id`,`name`,`time`,`status` FROM `" +
		m.table +
		"` WHERE `id` > ? ORDER BY `id` ASC LIMIT ?"

	for {
		count := 0
		for tries := 0; tries < retry; tries++ {
			rows, err := m.db.Query(query, lastID, limit)
			if err != nil {
				if err == sql.ErrNoRows {
					break
				}
				fmt.Println(`Query Logs:`, err)
				continue
			}
			for rows.Next() {
				log := &manager.Log{}
				err := rows.Scan(&log.ID, &log.Name, &log.Time, &log.Status)
				if err != nil {
					fmt.Println(`Scan Logs:`, err)
					continue
				}
				if log.ID > 0 {
					count++
					lastID = log.ID
					logs = append(logs, log)
				}
			}
			break
		}

		if count < limit {
			break
		}
	}

	return
}
