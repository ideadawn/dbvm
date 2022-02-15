package mysql

import (
	"bytes"
	"database/sql"
	"fmt"

	driver "github.com/go-sql-driver/mysql"
	"github.com/ideadawn/dbvm/manager"
)

// Revert 回退指定版本
func (m *MySQL) Revert(plan *manager.Plan) error {
	if m == nil || m.db == nil {
		return errConnection
	}
	if m.table == `` {
		return errTableNotInit
	}

	parser := &sqlParser{
		file: plan.Revert,
	}
	parser.parseSqlBlocks()
	if parser.err != nil {
		return parser
	}

	parser.blocks = append(parser.blocks, &sqlBlock{
		line: 999999,
		sqlArr: [][]byte{
			[]byte(fmt.Sprintf(
				"DELETE FROM `%s` WHERE `name` = '%s' LIMIT 1;",
				m.table,
				plan.Name,
			)),
		},
	})

	var tx *sql.Tx
	var err error
	for _, block := range parser.blocks {
		for tries := 0; tries < retry; tries++ {
			if !block.noTrans {
				tx, err = m.db.Begin()
				if err != nil {
					fmt.Println(`Revert BEGIN:`, err)
					continue
				}
			} else {
				tx = nil
			}

			exec := string(bytes.Join(block.sqlArr, myCnf.newLine))
			if tx == nil {
				_, err = m.db.Exec(exec)
			} else {
				_, err = tx.Exec(exec)
			}
			if err == nil {
				if block.noTrans {
					break
				}
				err = tx.Commit()
				if err == nil {
					break
				}
				fmt.Println(`Revert COMMIT:`, err)
				continue
			} else {
				fmt.Println(`Revert`, plan.Revert, `on line`, block.line, `:`, err)
				fmt.Println(exec)
				fmt.Println("")
			}

			if tx != nil {
				_ = tx.Rollback()
			}

			//是否语法错误，是否可以忽略
			myerr, ok := err.(*driver.MySQLError)
			if ok {
				if inUint16Array(block.ignores, myerr.Number) {
					err = nil
					break
				}
				return err
			}
		}

		if err != nil {
			break
		}
	}

	return err
}
