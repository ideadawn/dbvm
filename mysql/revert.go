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
		items: []*sqlItem{
			&sqlItem{
				line: 999999,
				sqlArr: [][]byte{
					[]byte(fmt.Sprintf(
						"DELETE FROM `%s` WHERE `name` = '%s' LIMIT 1;",
						m.table,
						plan.Name,
					)),
				},
			},
		},
	})

	newLine := []byte{'\n'}
	var tx *sql.Tx
	var err error
	for _, blk := range parser.blocks {
		for tries := 0; tries < retry; tries++ {
			tx, err = m.db.Begin()
			if err != nil {
				fmt.Println(`Revert BEGIN:`, err)
				continue
			}

			lastOffset := len(blk.items) - 1
			for offset, itm := range blk.items {
				exec := string(bytes.Join(itm.sqlArr, newLine))
				_, err = tx.Exec(exec)
				if err == nil {
					continue
				}

				fmt.Println(`Revert `, plan.Revert, `on line`, itm.line, `:`, err)
				fmt.Println(exec)
				fmt.Println("")

				//语法错误
				myerr, ok := err.(*driver.MySQLError)
				if !ok {
					break
				}
				if inUint16Array(blk.ignores, myerr.Number) {
					err = nil
					if offset == lastOffset {
						_ = tx.Rollback()
						tx = nil
						break
					}
					_ = tx.Rollback()
					tx, err = m.db.Begin()
					if err == nil {
						continue
					}
					break
				}
				_ = tx.Rollback()
				return err
			}

			if err != nil {
				if tx != nil {
					_ = tx.Rollback()
				}
				continue
			}

			if tx != nil {
				err = tx.Commit()
			}
			if err == nil {
				break
			}
			fmt.Println(`Revert COMMIT:`, err)
		}

		if err != nil {
			break
		}
	}

	return err
}
