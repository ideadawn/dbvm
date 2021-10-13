package mysql

import (
	"bytes"
	"os"
	"strconv"
)

// SQL语句集
type sqlItem struct {
	Line   int
	SqlArr []string
}

// SQL事务块
type sqlBlock struct {
	Ignore []uint16
	Items  []*sqlItem
}

// 解析SQL语句块
func parseSqlBlocks(file string) (blocks []*sqlBlock, err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		newLine       = []byte{'\n'}
		sqlEnd        = []byte{';'}
		commaGap      = []byte{','}
		blockBegin    = []byte(`BEGIN`)
		blockIgnore   = []byte(`IGNORE`)
		blockCommit   = []byte(`COMMIT`)
		blockRollback = []byte(`ROLLBACK`)
	)

	lines := bytes.Split(data, newLine)
	var blk *sqlBlock
	var ignores []uint16
	var itm *sqlItem

	for idx, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if len(line) > 1 && line[0] == '-' && line[1] == '-' {
			line = bytes.TrimSpace(line[2:])
			if !bytes.HasPrefix(line, blockIgnore) {
				continue
			}
			lArr := bytes.Split(line[len(blockIgnore):], commaGap)
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
			continue
		}
		if bytes.HasPrefix(line, blockBegin) {
			blk = &sqlBlock{
				Ignore: ignores,
			}
			continue
		}
		if bytes.HasPrefix(line, blockCommit) || bytes.HasPrefix(line, blockRollback) {
			if blk == nil {
				continue
			}
			if len(blk.Items) > 0 {
				blocks = append(blocks, blk)
			}
			blk = nil
			ignores = []uint16{}
			continue
		}
		if blk != nil {
			if itm == nil {
				itm = &sqlItem{
					Line: idx + 1,
				}
			}
			itm.SqlArr = append(itm.SqlArr, string(line))
			if bytes.HasSuffix(line, sqlEnd) {
				blk.Items = append(blk.Items, itm)
				itm = nil
			}
		}
	}

	return blocks, nil
}
