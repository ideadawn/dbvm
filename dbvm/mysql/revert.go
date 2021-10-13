package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"askc/tool/dbvm/manager"

	driver "github.com/go-sql-driver/mysql"
)

// Revert 回退指定版本
func (m *MySQL) Revert(plan *manager.Plan) error {
	if m == nil || m.db == nil {
		return errConnection
	}
	if m.table == `` {
		return errTableNotInit
	}

	blocks, err := parseSqlBlocks(plan.Revert)
	if err != nil {
		return err
	}

	blocks = append(blocks, &sqlBlock{
		Items: []*sqlItem{
			&sqlItem{
				Line: 999999,
				SqlArr: []string{
					fmt.Sprintf(
						"DELETE FROM `%s` WHERE `name` = '%s' LIMIT 1;",
						m.table,
						plan.Name,
					),
				},
			},
		},
	})

	var tx *sql.Tx
	for _, blk := range blocks {
		for tries := 0; tries < retry; tries++ {
			tx, err = m.db.Begin()
			if err != nil {
				fmt.Println(`BEGIN revert:`, err)
				continue
			}

			lastOffset := len(blk.Items) - 1
			for offset, itm := range blk.Items {
				exec := strings.Join(itm.SqlArr, "\n")
				_, err = tx.Exec(exec)
				if err == nil {
					continue
				}

				fmt.Println(`EXEC revert: line=`, itm.Line, `, err=`, err)
				fmt.Println(exec)

				//语法错误
				myerr, ok := err.(*driver.MySQLError)
				if !ok {
					break
				}
				if inErrCodes(myerr.Number, blk.Ignore) {
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
			fmt.Println(`COMMIT revert:`, err)
		}

		if err != nil {
			break
		}
	}

	return err
}
