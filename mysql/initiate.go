package mysql

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/VividCortex/mysqlerr"
	driver "github.com/go-sql-driver/mysql"
)

// Initiate 初始化日志表
func (m *MySQL) Initiate(table string) error {
	if m == nil || m.db == nil {
		return errConnection
	}

	reName := regexp.MustCompile("^[a-z0-9_]+$")
	if !reName.Match([]byte(table)) {
		return errTableName
	}
	m.table = table

	query := "SELECT `name` FROM `" + table + "` LIMIT 1"
	row := m.db.QueryRow(query)
	var name string
	err := row.Scan(&name)
	if err == nil || err == sql.ErrNoRows {
		return nil
	}

	myerr, ok := err.(*driver.MySQLError)
	if !ok {
		return err
	}
	if myerr.Number != mysqlerr.ER_NO_SUCH_TABLE {
		return err
	}

	create := strings.Join([]string{
		"CREATE TABLE IF NOT EXISTS `" + table + "` (",
		"	`id` INT NOT NULL AUTO_INCREMENT COMMENT '日志ID',",
		"	`name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '版本名称' COLLATE 'ascii_bin',",
		"	`time` BIGINT NOT NULL DEFAULT 0 COMMENT '部署时间',",
		"	`status` TINYINT NOT NULL DEFAULT 0 COMMENT '部署状态',",
		"	PRIMARY KEY (`id`),",
		"	UNIQUE INDEX `name` (`name`)",
		")",
		"COMMENT='数据库版本更新日志表'",
		"COLLATE='utf8_bin';",
	}, "\n")
	_, err = m.db.Exec(create)
	return err
}
