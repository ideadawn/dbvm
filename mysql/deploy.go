package mysql

import (
	"bytes"
	"database/sql"
	"fmt"

	driver "github.com/go-sql-driver/mysql"
	"github.com/ideadawn/dbvm/manager"
)

// Deploy 部署指定版本
func (m *MySQL) Deploy(plan *manager.Plan) error {
	if m == nil || m.db == nil {
		return errConnection
	}
	if m.table == `` {
		return errTableNotInit
	}

	parser := &sqlParser{
		file: plan.Deploy,
	}

	parser.parseSqlBlocks()
	if parser.err != nil {
		return parser
	}
	if len(parser.blocks) == 0 {
		return errDeployNothing
	}

	parser.blocks = append(parser.blocks, &sqlBlock{
		items: []*sqlItem{
			&sqlItem{
				line: 999999,
				sqlArr: [][]byte{
					[]byte(fmt.Sprintf(
						"INSERT INTO `%s` (`name`,`time`,`status`) VALUES ('%s', '%d', '%d');",
						m.table,
						plan.Name,
						plan.Time.Unix(),
						StatusDeployed,
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
				fmt.Println(`Deploy BEGIN:`, err)
				continue
			}

			lastOffset := len(blk.items) - 1
			for offset, itm := range blk.items {
				exec := string(bytes.Join(itm.sqlArr, newLine))
				_, err = tx.Exec(exec)
				if err == nil {
					continue
				}

				fmt.Println(`Deploy`, plan.Deploy, `on line`, itm.line, `:`, err)
				fmt.Println(exec)
				fmt.Println("")

				//是否语法错误
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
			fmt.Println(`Deploy COMMIT:`, err)
		}

		if err != nil {
			break
		}
	}

	return err
}
