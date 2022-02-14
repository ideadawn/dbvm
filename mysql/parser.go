package mysql

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/VividCortex/mysqlerr"
	"github.com/ideadawn/dbvm/manager"
)

// SQL语句集
type sqlItem struct {
	line     int
	comments [][]byte
	sqlArr   [][]byte
}

// SQL事务块
type sqlBlock struct {
	noTrans   bool
	ignores   []uint16
	inBlock   bool
	comments  [][]byte
	delimiter []byte
	items     []*sqlItem
}

// SQL解析器
type sqlParser struct {
	file   string
	blocks []*sqlBlock

	line int
	sql  string
	err  error
}

// Error implement error
func (p *sqlParser) Error() string {
	if p == nil {
		return ``
	}
	if p.sql != `` {
		return fmt.Sprintf("%s on line %d: %s\n\t%s", p.file, p.line, p.err, p.sql)
	}
	return fmt.Sprintf("%s on line %d: %s", p.file, p.line, p.err)
}

// reset parser
func (p *sqlParser) reset(file string) {
	p.file = file
	p.blocks = p.blocks[0:0]
	p.line = 0
	p.sql = ``
	p.err = nil
}

// 解析SQL语句块
func (p *sqlParser) parseSqlBlocks() {
	var data []byte
	data, p.err = os.ReadFile(p.file)
	if p.err != nil {
		return
	}

	data = bytes.Replace(data, []byte("\r\n"), myCnf.newLine, -1)
	data = bytes.Replace(data, []byte("\r"), myCnf.newLine, -1)
	lines := bytes.Split(data, myCnf.newLine)

	var (
		sqlEnd = myCnf.defaultEnd
		block  = &sqlBlock{}
		item   = &sqlItem{}
	)

	for idx, line := range lines {
		data = bytes.TrimSpace(line)
		if len(data) == 0 {
			if block.inBlock {
				item.comments = append(item.comments, myCnf.empty)
			} else {
				block.comments = append(block.comments, myCnf.empty)
			}
			continue
		}

		if bytes.HasPrefix(data, myCnf.blockBegin) {
			if len(item.sqlArr) > 0 || len(block.items) > 0 {
				if block.inBlock {
					p.line = idx
					p.err = errBlockNoEnd
					return
				}

				block.items = append(block.items, item)
				p.analyzeBlock(block)
				if p.err != nil {
					return
				}
				block = &sqlBlock{}
				item = &sqlItem{}
			}
			block.inBlock = true
			continue
		}

		if bytes.HasPrefix(data, myCnf.blockCommit) {
			if !block.inBlock {
				p.line = idx
				p.err = errBlockNoBegin
				return
			}

			if len(item.sqlArr) > 0 {
				block.items = append(block.items, item)
			}
			if len(block.items) > 0 {
				p.analyzeBlock(block)
				if p.err != nil {
					return
				}
			}

			block = &sqlBlock{
				inBlock: false,
			}
			item = &sqlItem{}
			continue
		}

		if bytes.HasPrefix(data, myCnf.commentBegin) {
			data = bytes.TrimSpace(data[2:])
			if bytes.Compare(data, manager.MagicNoTrans) == 0 {
				block.noTrans = true
			} else if bytes.HasPrefix(data, manager.MagicIgnore) {
				lArr := bytes.Split(data[len(manager.MagicIgnore):], myCnf.commaGap)
				for _, iData := range lArr {
					iData = bytes.TrimSpace(iData)
					if len(iData) == 0 {
						continue
					}
					u64, err := strconv.ParseUint(string(iData), 10, 16)
					if err != nil {
						p.line = idx
						p.sql = string(line)
						p.err = err
						return
					}
					appendUint16Array(&block.ignores, uint16(u64))
				}
			} else {
				if block.inBlock {
					item.comments = append(item.comments, line)
				} else {
					block.comments = append(block.comments, line)
				}
			}
			continue
		}

		if bytes.HasPrefix(data, myCnf.delimiter) {
			data = bytes.Replace(data, myCnf.delimiter, myCnf.empty, -1)
			sqlEnd = bytes.TrimSpace(data)
			if bytes.Compare(sqlEnd, myCnf.defaultEnd) != 0 {
				block.delimiter = sqlEnd
			}
			continue
		}

		if bytes.HasSuffix(data, sqlEnd) {
			item.sqlArr = append(item.sqlArr, line)
			block.items = append(block.items, item)
			if block.inBlock {
				item = &sqlItem{}
			} else {
				p.analyzeBlock(block)
				if p.err != nil {
					return
				}
				block = &sqlBlock{}
				item = &sqlItem{}
			}
		} else {
			if len(item.sqlArr) == 0 {
				item.line = idx
			}
			item.sqlArr = append(item.sqlArr, line)
		}
	}
}

// 分析事务块
func (p *sqlParser) analyzeBlock(blk *sqlBlock) {
	items := make([]*sqlItem, 0, len(blk.items))
	items = append(items, blk.items...)
	blk.items = blk.items[0:0]

	for _, item := range items {
		sqlBytes := bytes.Join(item.sqlArr, myCnf.newLine)
		if myCnf.reCreateTable.Match(sqlBytes) {
			if myCnf.reCreateTableINE.Match(sqlBytes) {
				blk.items = append(blk.items, item)
				continue
			}
			p.line = item.line
			p.sql = string(sqlBytes)
			p.err = errCreateTableINE
			return
		}

		if myCnf.reDropTable.Match(sqlBytes) {
			if myCnf.reDropTableIE.Match(sqlBytes) {
				blk.items = append(blk.items, item)
				continue
			}
			p.line = item.line
			p.sql = string(sqlBytes)
			p.err = errDropTableIE
			return
		}

		alterArr := myCnf.reAlter.FindSubmatch(sqlBytes)
		if len(alterArr) == 3 {
			p.splitAlter(blk, item, alterArr)
			if p.err != nil {
				return
			}
		} else {
			blk.items = append(blk.items, item)
		}
	}

	p.blocks = append(p.blocks, blk)
}

// 拆分ALTER
func (p *sqlParser) splitAlter(blk *sqlBlock, item *sqlItem, alterArr [][]byte) {
	subArr := myCnf.reAlterSub.FindAllSubmatch(alterArr[2], -1)
	comments := item.comments
	for idx, val := range subArr {
		var errNum uint16
		if idx > 0 {
			comments = comments[0:0]
		}

		if myCnf.reAddColumn.Match(val[1]) {
			errNum = mysqlerr.ER_DUP_FIELDNAME
		} else if myCnf.reAddIndex.Match(val[1]) {
			errNum = mysqlerr.ER_DUP_KEYNAME
		} else if myCnf.reAddPrimary.Match(val[1]) {
			errNum = mysqlerr.ER_MULTIPLE_PRI_KEY
		} else if myCnf.reChangeColumn.Match(val[1]) {
			errNum = mysqlerr.ER_BAD_FIELD_ERROR
		} else if myCnf.reDropColumn.Match(val[1]) {
			errNum = mysqlerr.ER_CANT_DROP_FIELD_OR_KEY
		} else if myCnf.reModifyColumn.Match(val[1]) {
			errNum = 1
		}
		if errNum == 0 {
			p.line = item.line + idx + 1
			p.sql = string(val[1])
			p.err = errAlterUnknown
			return
		}

		if errNum > 1 {
			appendUint16Array(&blk.ignores, errNum)
		}
		var delimiter []byte
		if len(blk.delimiter) > 0 {
			delimiter = blk.delimiter
		} else {
			delimiter = myCnf.defaultEnd
		}
		blk.items = append(blk.items, &sqlItem{
			line:     item.line + idx + 1,
			comments: comments,
			sqlArr: [][]byte{
				bytes.Join([][]byte{
					alterArr[1],
					myCnf.space,
					bytes.TrimRight(val[1], ", \t\n"+string(delimiter)),
					delimiter,
				}, myCnf.empty),
			},
		})
	}
}
