package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	driver "github.com/go-sql-driver/mysql"
	"github.com/ideadawn/dbvm/manager"
)

func inErrCodes(code uint16, arr []uint16) bool {
	for _, val := range arr {
		if val == code {
			return true
		}
	}
	return false
}

// Deploy 部署指定版本
func (m *MySQL) Deploy(plan *manager.Plan) error {
	if m == nil || m.db == nil {
		return errConnection
	}
	if m.table == `` {
		return errTableNotInit
	}

	blocks, err := parseSqlBlocks(plan.Deploy)
	if err != nil {
		return err
	}
	lth := len(blocks)
	if lth == 0 {
		return errDeployNothing
	}

	blocks = append(blocks, &sqlBlock{
		Items: []*sqlItem{
			&sqlItem{
				Line: 999999,
				SqlArr: []string{
					fmt.Sprintf(
						"INSERT INTO `%s` (`name`,`time`,`status`) VALUES ('%s', '%d', '%d');",
						m.table,
						plan.Name,
						plan.Time.Unix(),
						StatusDeployed,
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
				fmt.Println(`BEGIN deploy:`, err)
				continue
			}

			lastOffset := len(blk.Items) - 1
			for offset, itm := range blk.Items {
				exec := strings.Join(itm.SqlArr, "\n")
				_, err = tx.Exec(exec)
				if err == nil {
					continue
				}

				fmt.Println(`Deploy: file=`, plan.Deploy, ` line=`, itm.Line, `, err=`, err)
				fmt.Println(exec, "\n")

				//是否语法错误
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
			fmt.Println(`COMMIT deploy:`, err)
		}

		if err != nil {
			break
		}
	}

	return err
}
