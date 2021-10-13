package mysql

import (
	"os"
	"strings"
	"testing"

	"github.com/nbio/st"
)

func Test_ParseSqlBlock(t *testing.T) {
	data := []byte(strings.Join([]string{
		"-- Deploy kc:v1.6.0 to mysql",
		"",
		"-- IGNORE 1091",
		"BEGIN;",
		"",
		"ALTER TABLE `test`",
		"	ADD COLUMN `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '姓名' AFTER `id`;",
		"",
		"ALTER TABLE `test` DROP INDEX `not_exists`;",
		"",
		"ALTER TABLE `test` DROP COLUMN `name`;",
		"",
		"COMMIT;",
		"",
		"BEGIN;",
		"",
		"CREATE TABLE IF NOT EXISTS `test` (",
		"	`id` INT(10) NOT NULL AUTO_INCREMENT COMMENT '自增ID',",
		"	PRIMARY KEY (`id`) USING BTREE",
		")",
		"COMMENT='测试'",
		"COLLATE='utf8_bin'",
		"ENGINE=InnoDB;",
		"",
		"COMMIT;",
		"",
	}, "\n"))
	file := `./temp-plan.sql`

	err := os.WriteFile(file, data, os.ModePerm)
	st.Assert(t, err, nil)

	blocks, err := parseSqlBlocks(file)
	st.Assert(t, err, nil)
	st.Assert(t, len(blocks), 2)
	st.Assert(t, len(blocks[0].Items), 3)
	st.Assert(t, blocks[0].Ignore, []uint16{1091})
	st.Assert(t, blocks[1].Ignore, []uint16{})

	_ = os.Remove(file)
}
