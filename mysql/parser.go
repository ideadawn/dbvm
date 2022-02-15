package mysql

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/VividCortex/mysqlerr"
	"github.com/ideadawn/dbvm/manager"
)

// SQL事务块
type sqlBlock struct {
	comments  [][]byte
	noTrans   bool
	ignores   []uint16
	delimiter []byte
	line      int
	sqlArr    [][]byte
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
		block  = &sqlBlock{
			delimiter: myCnf.defaultEnd,
		}
	)

	for idx, line := range lines {
		data = bytes.TrimSpace(line)
		if len(data) == 0 {
			block.comments = append(block.comments, myCnf.empty)
			continue
		}
		if bytes.HasPrefix(data, myCnf.blockBegin) || bytes.HasPrefix(data, myCnf.blockCommit) {
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
				block.comments = append(block.comments, line)
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
			data = bytes.TrimRight(line, " \t\n")
			if bytes.Compare(sqlEnd, myCnf.defaultEnd) != 0 {
				block.sqlArr = append(block.sqlArr, data[0:len(data)-len(sqlEnd)])
			} else {
				block.sqlArr = append(block.sqlArr, data)
			}
			p.analyzeBlock(block)
			if p.err != nil {
				return
			}
			block = &sqlBlock{
				delimiter: sqlEnd,
			}
		} else {
			if len(block.sqlArr) == 0 {
				block.line = idx
			}
			block.sqlArr = append(block.sqlArr, line)
		}
	}
}

// 分析事务块
func (p *sqlParser) analyzeBlock(block *sqlBlock) {
	sqlBytes := bytes.Join(block.sqlArr, myCnf.newLine)
	if myCnf.reCreateTable.Match(sqlBytes) {
		if myCnf.reCreateTableINE.Match(sqlBytes) {
			p.analyzeSyntax(block)
		} else {
			p.line = block.line
			p.sql = string(sqlBytes)
			p.err = errCreateTableINE
		}
		return
	}

	if myCnf.reDropTable.Match(sqlBytes) {
		if myCnf.reDropTableIE.Match(sqlBytes) {
			p.analyzeSyntax(block)
		} else {
			p.line = block.line
			p.sql = string(sqlBytes)
			p.err = errDropTableIE
		}
		return
	}

	alterArr := myCnf.reAlter.FindSubmatch(sqlBytes)
	if len(alterArr) == 3 {
		p.splitAlter(block, alterArr)
		if p.err != nil {
			return
		}
	} else {
		p.analyzeSyntax(block)
	}
}

// 拆分ALTER
func (p *sqlParser) splitAlter(block *sqlBlock, alterArr [][]byte) {
	var end byte
	var backslash int = -1
	var subArr [][]byte
	var current []byte
	idx := 1
	for pos, v := range alterArr[2] {
		if v > 127 && (end == 0 || end == ')') {
			p.line = block.line + idx
			p.sql = string(alterArr[2])
			p.err = errSyntaxError
			return
		}
		current = append(current, v)
		switch v {
		case '\\':
			if backslash == -1 {
				backslash = pos
			} else if backslash == pos-1 {
				backslash = -1
			}
		case '\'':
			if end == 0 {
				end = '\''
			} else if backslash != -1 && backslash != pos-1 {
				end = 0
			} else if backslash == pos-1 {
				backslash = -1
			}
		case '"':
			if end == 0 {
				end = '"'
			} else if backslash != -1 && backslash != pos-1 {
				end = 0
			} else if backslash == pos-1 {
				backslash = -1
			}
		case '`':
			if end == 0 {
				end = '`'
			} else if end == '`' {
				end = 0
			}
		case '(':
			if end == 0 {
				end = ')'
			}
		case ')':
			if end == ')' {
				end = 0
			}
		case ',':
			fallthrough
		case ';':
			if end == 0 {
				tmp := bytes.TrimSpace(current)
				if len(tmp) > 0 {
					cpy := make([]byte, 0, len(current))
					tmp = bytes.TrimRight(current, " \t\n")
					tmp = bytes.TrimLeft(tmp, "\n")
					cpy = append(cpy, tmp...)
					subArr = append(subArr, cpy)
				}
				current = current[0:0]
			}
		case '\n':
			idx++
		}
	}
	tmp := bytes.TrimSpace(current)
	if len(tmp) > 0 {
		tmp = bytes.TrimRight(current, " \t\n")
		tmp = bytes.TrimLeft(tmp, "\n")
		subArr = append(subArr, tmp)
	}

	var errNum uint16
	for idx, sub := range subArr {
		errNum = 0
		if myCnf.reAddColumn.Match(sub) {
			errNum = mysqlerr.ER_DUP_FIELDNAME
		} else if myCnf.reAddIndex.Match(sub) {
			errNum = mysqlerr.ER_DUP_KEYNAME
		} else if myCnf.reAddPrimary.Match(sub) {
			errNum = mysqlerr.ER_MULTIPLE_PRI_KEY
		} else if myCnf.reChangeColumn.Match(sub) {
			errNum = mysqlerr.ER_BAD_FIELD_ERROR
		} else if myCnf.reDropColumn.Match(sub) {
			errNum = mysqlerr.ER_CANT_DROP_FIELD_OR_KEY
		} else if myCnf.reModifyColumn.Match(sub) {
			errNum = 1
		}
		if errNum == 0 {
			p.line = block.line + idx + 1
			p.sql = string(sub)
			p.err = errAlterUnknown
			return
		}

		newBlk := &sqlBlock{
			noTrans:   block.noTrans,
			delimiter: myCnf.defaultEnd,
			line:      block.line + idx + 1,
			sqlArr: [][]byte{
				bytes.Join([][]byte{
					alterArr[1],
					bytes.TrimRight(sub, ",; \t\n"),
					myCnf.defaultEnd,
				}, myCnf.empty),
			},
		}
		if idx == 0 {
			newBlk.comments = block.comments
		}
		if errNum > 1 {
			appendUint16Array(&newBlk.ignores, errNum)
		}
		p.blocks = append(p.blocks, newBlk)
	}
}

// 分析是否存在语法错误
func (p *sqlParser) analyzeSyntax(block *sqlBlock) {
	var end byte
	var backslash int = -1
	for idx, sub := range block.sqlArr {
		for pos, v := range sub {
			if v > 127 && end == 0 {
				p.line = block.line + idx + 1
				p.sql = string(sub)
				p.err = errSyntaxError
				return
			}
			switch v {
			case '\\':
				backslash = pos
			case '\'':
				if end == 0 {
					end = '\''
				} else if backslash != pos-1 {
					end = 0
				}
			case '"':
				if end == 0 {
					end = '"'
				} else if backslash != pos-1 {
					end = 0
				}
			}
		}
	}
	p.blocks = append(p.blocks, block)
}
