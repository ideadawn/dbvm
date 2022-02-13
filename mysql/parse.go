package mysql

import (
	"bytes"
	"os"
	"regexp"
	"strconv"

	"github.com/VividCortex/mysqlerr"
	"github.com/ideadawn/dbvm/manager"
)

var (
	reAlter        = regexp.MustCompile("(?is)(ALTER[ \t\n]+TABLE.*?)((?:ADD|CHANGE|MODIFY|DROP).*)")
	reAlterSub     = regexp.MustCompile("(?is)((?:ADD|CHANGE|MODIFY|DROP)[ \t\n]+(?:COLUMN|INDEX|KEY|PRIMARY|UNIQUE).*?(?:,[ \t\n]+|;))")
	reAddColumn    = regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+COLUMN")
	reAddPrimary   = regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+PRIMARY")
	reAddIndex     = regexp.MustCompile("^(?is)[ \t\n]*ADD[ \t\n]+(?:UNIQUE[ \t\n]+)?(?:INDEX|KEY)")
	reDropColumn   = regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+(?:COLUMN|INDEX|KEY|PRIMARY)")
	reChangeColumn = regexp.MustCompile("^(?is)[ \t\n]*CHANGE[ \t\n]+COLUMN")
	reModifyColumn = regexp.MustCompile("^(?is)[ \t\n]*MODIFY[ \t\n]+COLUMN")

	reCreateTable    = regexp.MustCompile("^(?is)[ \t\n]*CREATE[ \t\n]+TABLE")
	reCreateTableINE = regexp.MustCompile("^(?is)[ \t\n]*CREATE[ \t\n]+TABLE[ \t\n]+IF[ \t\n]+NOT[ \t\n]+EXISTS")
	reDropTable      = regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+TABLE")
	reDropTableIE    = regexp.MustCompile("^(?is)[ \t\n]*DROP[ \t\n]+TABLE[ \t\n]+IF[ \t\n]+EXISTS")
)

// SQL语句集
type sqlItem struct {
	line     int
	comments [][]byte
	sqlArr   [][]byte
}

// SQL事务块
type sqlBlock struct {
	noTrans bool
	ignores []uint16
	inBlock bool
	items   []*sqlItem
}

// SQL解析器
type sqlParser struct {
	file   string
	blocks []*sqlBlock

	line int
	sql  string
	err  error
}

// 解析SQL语句块
func (p *sqlParser) parseSqlBlocks() {
	var data []byte
	data, p.err = os.ReadFile(p.file)
	if p.err != nil {
		return
	}

	var (
		newLine   = []byte{'\n'}
		sqlEnd    = []byte{';'}
		commaGap  = []byte{','}
		delimiter = []byte(`DELIMITER`)
		empty     = []byte{}

		commentBegin  = []byte(`--`)
		blockBegin    = []byte(`BEGIN`)
		blockCommit   = []byte(`COMMIT`)
		blockRollback = []byte(`ROLLBACK`)
	)

	data = bytes.Replace(data, []byte("\r\n"), newLine, -1)
	data = bytes.Replace(data, []byte("\r"), newLine, -1)
	lines := bytes.Split(data, newLine)

	var (
		block = &sqlBlock{}
		item  = &sqlItem{}
	)

	for idx, line := range lines {
		data = bytes.TrimSpace(line)
		if len(data) == 0 {
			continue
		}

		if bytes.HasPrefix(data, blockBegin) {
			if len(item.sqlArr) > 0 {
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

		if bytes.HasPrefix(data, blockCommit) || bytes.HasPrefix(data, blockRollback) {
			if !block.inBlock {
				p.line = idx
				p.err = errBlockNoBegin
				return
			}

			if len(item.sqlArr) > 0 {
				block.items = append(block.items, item)
				p.analyzeBlock(block)
				if p.err != nil {
					return
				}
				block = &sqlBlock{}
				item = &sqlItem{}
			}

			block.inBlock = false
			continue
		}

		if bytes.HasPrefix(data, commentBegin) {
			data = bytes.TrimSpace(data[2:])
			if bytes.Compare(data, manager.MagicNoTrans) == 0 {
				block.noTrans = true
			} else if bytes.HasPrefix(data, manager.MagicIgnore) {
				lArr := bytes.Split(data[len(manager.MagicIgnore):], commaGap)
				for _, iData := range lArr {
					iData = bytes.TrimSpace(iData)
					if len(iData) == 0 {
						continue
					}
					u64, err := strconv.ParseUint(string(iData), 10, 16)
					if err != nil {
						p.line = idx
						p.sql = line
						p.err = err
						return
					}
					appendUint16Array(&block.ignores, uint16(u64))
				}
			} else {
				item.comments = append(item.comments, line)
			}
			continue
		}

		if bytes.HasPrefix(data, delimiter) {
			data = bytes.Replace(data, delimiter, empty, -1)
			sqlEnd = bytes.TrimSpace(data)
			continue
		}

		if bytes.HasSuffix(data, sqlEnd) {
			item.sqlArr = append(item.sqlArr, line)
			block.items = append(block.items, item)
			if !block.inBlock {
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
	items := blk.items
	blk.items = blk.items[0:0]
	for _, item := range items {
		sqlBytes := bytes.Join(item.sqlArr, []byte{'\n'})
		if reCreateTable.Match(sqlBytes) {
			if reCreateTableINE.Match(sqlBytes) {
				blk.items = append(blk.items, item)
				continue
			}
			p.line = item.line
			p.sql = string(sqlBytes)
			p.err = errCreateTableINE
			return
		}

		if reDropTable.Match(sqlBytes) {
			if reDropTableIE.Match(sqlBytes) {
				blk.items = append(blk.items, item)
				continue
			}
			p.line = item.line
			p.sql = string(sqlBytes)
			p.err = errDropTableIE
			return
		}

		alterArr := reAlter.FindSubmatch(sqlBytes)
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
	subArr := reAlterSub.FindAllSubmatch(alterArr[2], -1)
	comments := item.comments
	for idx, val := range subArr {
		var errNum uint16
		if idx > 0 {
			comments = comments[0:0]
		}

		if reAddColumn.Match(val[1]) {
			errNum = mysqlerr.ER_DUP_FIELDNAME
		} else if reAddIndex.Match(val[1]) {
			errNum = mysqlerr.ER_DUP_KEYNAME
		} else if reAddPrimary.Match(val[1]) {
			errNum = mysqlerr.ER_MULTIPLE_PRI_KEY
		} else if reChangeColumn.Match(val[1]) {
			errNum = mysqlerr.ER_BAD_FIELD_ERROR
		} else if reDropColumn.Match(val[1]) {
			errNum = mysqlerr.ER_CANT_DROP_FIELD_OR_KEY
		} else if reModifyColumn.Match(val[1]) {
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
		blk.items = append(blk.items, &sqlItem{
			line:     item.line + idx + 1,
			comments: comments,
			sqlArr: [][]byte{
				alterArr[1],
				[]byte{' '},
				bytes.TrimRight(val[1], ",; \t\n"),
				[]byte{';'},
			},
		})
	}
}
