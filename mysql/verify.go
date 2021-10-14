package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"dbvm/manager"

	driver "github.com/go-sql-driver/mysql"
)

// Verify 检查版本部署
func (m *MySQL) Verify(plan *manager.Plan) error {
	if m == nil || m.db == nil {
		return errConnection
	}
	if m.table == `` {
		return errTableNotInit
	}

	blocks, err := parseSqlBlocks(plan.Verify)
	if err != nil {
		return err
	}
	if len(blocks) == 0 {
		blocks = append(blocks, &sqlBlock{
			Items: []*sqlItem{
				&sqlItem{
					Line: 999999,
					SqlArr: []string{
						fmt.Sprintf(
							"SELECT `id` FROM `%s` WHERE `name` = '%s' LIMIT 1;",
							m.table,
							plan.Name,
						),
					},
				},
			},
		})
	}

	var tx *sql.Tx
	for _, blk := range blocks {
		for tries := 0; tries < retry; tries++ {
			tx, err = m.db.Begin()
			if err != nil {
				fmt.Println(`BEGIN verify:`, err)
				continue
			}

			lastOffset := len(blk.Items) - 1
			for offset, itm := range blk.Items {
				exec := strings.Join(itm.SqlArr, "\n")
				if strings.HasPrefix(exec, `SELECT`) || strings.HasPrefix(exec, `select`) {
					var rows *sql.Rows
					rows, err = tx.Query(exec)
					if rows != nil {
						_ = rows.Close()
					}
				} else {
					_, err = tx.Exec(exec)
				}
				if err == nil {
					continue
				}

				fmt.Println(`EXEC verify: line=`, itm.Line, `, err=`, err)
				fmt.Println(exec)

				//语法错误
				myerr, ok := err.(*driver.MySQLError)
				if !ok {
					break
				}
				if inErrCodes(myerr.Number, blk.Ignore) {
					err = nil
					if offset == lastOffset {
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

			if tx != nil {
				_ = tx.Rollback()
			}
			if err == nil {
				break
			}
			fmt.Println(`ROLLBACK verify:`, err)
		}

		if err != nil {
			return err
		}
	}

	_, err = m.db.Exec(
		"UPDATE `"+m.table+"` SET `status` = ? WHERE `name` = ? LIMIT 1",
		StatusVerified,
		plan.Name,
	)
	return err
}
