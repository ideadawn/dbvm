package mysql

import (
	"bytes"
	"os"
	"strconv"

	"github.com/ideadawn/dbvm/manager"
)

// SQL事务块
type sqlBlock struct {
	Ignore uint16
	Line   int
	Origin string
	Sql    string
}

// 解析SQL语句块
func parseSqlBlocks(file string, rule *manager.Rule) (blocks []*sqlBlock, err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		newLine   = []byte{'\n'}
		sqlEnd    = []byte{';'}
		delimiter = []byte(`DELIMITER`)
		empty     = []byte{}

		blockBegin    = []byte(`BEGIN`)
		blockCommit   = []byte(`COMMIT`)
		blockRollback = []byte(`ROLLBACK`)
	)

	data = bytes.Replace(data, []byte{'\r'}, newLine, -1)
	lines := bytes.Split(data, newLine)
	sqlArr := make([][]byte, 0, 64)

	for idx, line := range lines {
		data = bytes.TrimSpace(line)
		if len(data) == 0 {
			continue
		}
		if bytes.HasPrefix(data, blockBegin) || bytes.HasPrefix(data, blockCommit) || bytes.HasPrefix(data, blockRollback) {
			continue
		}
		if bytes.HasPrefix(data, delimiter) {
			data = bytes.Replace(data, delimiter, empty, -1)
			sqlEnd = bytes.TrimSpace(data)
			continue
		}
		if bytes.HasSuffix(data, sqlEnd) {
			sqlArr = append(sqlArr, line)
			blkArr, err := parseStatement(sqlArr)
			if err != nil {
				return nil, err
			}
			if len(blkArr) > 0 {
				blocks = append(blocks, blkArr...)
			}
			sqlArr = sqlArr[0:0]
		} else {
			sqlArr = append(sqlArr, line)
		}
	}

	return blocks, nil
}

// 解析SQL执行语句
// CREATE ALTER DROP TRUNCATE SELECT INSERT DELETE UPDATE
// parseAlterTable ??
func parseStatement(sqlArr [][]byte) ([]*sqlBlock, error) {
	subEnd := []byte{','}
	var (
		command []byte
		target  []byte
	)
}
