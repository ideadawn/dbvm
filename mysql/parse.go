package mysql

import (
	"bytes"
	"os"
	"regexp"
	"strconv"

	"github.com/ideadawn/dbvm/manager"
)

// SQL语句集
type sqlItem struct {
	line   int
	sqlArr []string
}

// SQL事务块
type sqlBlock struct {
	ignore []uint16
	items  []*sqlItem
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
		commaGap  = []byte{','}
		delimiter = []byte(`DELIMITER`)
		empty     = []byte{}

		commentBegin  = []byte(`--`)
		blockIgnore   = []byte(`IGNORE`)
		blockBegin    = []byte(`BEGIN`)
		blockCommit   = []byte(`COMMIT`)
		blockRollback = []byte(`ROLLBACK`)

		reAlter = regexp.MustCompile("(?i)^ALTER[ \t\n]+TABLE")
	)

	data = bytes.Replace(data, []byte("\r\n"), newLine, -1)
	data = bytes.Replace(data, []byte("\r"), newLine, -1)
	lines := bytes.Split(data, newLine)
	sqlArr := make([][]byte, 0, 64)

	var (
		inBlock bool
		block   *sqlBlock
		ignores []uint16
		item    *sqlItem
	)

	for idx, line := range lines {
		data = bytes.TrimSpace(line)
		if len(data) == 0 {
			continue
		}
		if bytes.HasPrefix(data, blockBegin) {
			if len(sqlArr) > 0 {
				return nil, errSqlNotEnd
			}
			inBlock = true
			continue
		}
		if bytes.HasPrefix(data, blockCommit) || bytes.HasPrefix(data, blockRollback) {

			inBlock = false
			continue
		}
		if bytes.HasPrefix(data, commentBegin) {
			data = bytes.TrimSpace(data[2:])
			if bytes.Compare(data, manager.NoRevert) == 0 {
				rule.NoRevert = true
			} else if bytes.Compare(data, manager.NoTransaction) == 0 {
				rule.NoTransaction = true
			} else if bytes.HasPrefix(data, blockIgnore) {
				lArr := bytes.Split(data[len(blockIgnore):], commaGap)
				for _, iData := range lArr {
					iData = bytes.TrimSpace(iData)
					if len(iData) == 0 {
						continue
					}
					u64, err := strconv.ParseUint(string(iData), 10, 16)
					if err != nil {
						return nil, err
					}
					ignores = append(ignores, uint16(u64))
				}
			}
			continue
		}
		if bytes.HasPrefix(data, delimiter) {
			data = bytes.Replace(data, delimiter, empty, -1)
			sqlEnd = bytes.TrimSpace(data)
			continue
		}
		if bytes.HasSuffix(data, sqlEnd) {
			sqlArr = append(sqlArr, line)
			blkArr, err := analyseSqlBlock(sqlArr, idx+1, rule)
			if err != nil {
				return nil, err
			}
			if len(blkArr) > 0 {
				blocks = append(blocks, blkArr...)
			}
			sqlArr = sqlArr[0:0]
			rule.NoRevert = false
			rule.NoTransaction = false
		} else {
			sqlArr = append(sqlArr, line)
		}
	}

	return blocks, nil
}

// 拆分ALTER
func splitAlter(blk *sqlBlock, item *sqlItem) {

}
