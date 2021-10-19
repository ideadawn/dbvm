package mysql

import (
	"os"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VividCortex/mysqlerr"
	driver "github.com/go-sql-driver/mysql"
	"github.com/ideadawn/dbvm/manager"
	"github.com/nbio/st"
)

func Test_MySQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	st.Assert(t, err, nil)
	my := New()
	my.db = db
	my.table = `dbvm`

	plan := &manager.Plan{
		Name:     `v1.7.0`,
		Requires: []string{},
		Time:     time.Now(),
		Deploy:   `./tmp-deploy-v1.7.0.sql`,
		Verify:   `./tmp-verify-v1.7.0.sql`,
		Revert:   `./tmp-revert-v1.7.0.sql`,
	}

	mock.ExpectQuery("^SELECT `name` FROM `" + my.table + "` LIMIT 1$").
		WillReturnError(&driver.MySQLError{
			Number:  mysqlerr.ER_NO_SUCH_TABLE,
			Message: `No such table.`,
		})
	mock.ExpectExec("^CREATE TABLE IF NOT EXISTS `" + my.table + "` .*").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery("^SELECT `id`,`name`,`time`,`status` FROM `"+my.table+"` WHERE `id` > \\? ORDER BY `id` ASC LIMIT \\?$").
		WithArgs(int64(0), 1000).
		WillReturnRows(sqlmock.
			NewRows([]string{`id`, `name`, `time`, `status`}).
			AddRow(int64(1), `v1.6.0`, int64(162), int8(2)),
		)

	mock.ExpectBegin()
	mock.ExpectExec("^ALTER TABLE `test` ADD COLUMN `name` VARCHAR\\(32\\) NOT NULL DEFAULT '' AFTER `id`;$").
		WillReturnError(&driver.MySQLError{
			Number: mysqlerr.ER_DUP_FIELDNAME,
		})
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `" + my.table + "` \\(`name`,`time`,`status`\\) VALUES .*$").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test` \\(`name`\\) VALUES \\('test'\\);$").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	mock.ExpectExec("^UPDATE `"+my.table+"` SET `status` = \\? WHERE `name` = \\? LIMIT 1$").
		WithArgs(StatusVerified, plan.Name).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectBegin()
	mock.ExpectExec("^ALTER TABLE `test` DROP COLUMN `name`;$").
		WillReturnError(&driver.MySQLError{
			Number: mysqlerr.ER_CANT_DROP_FIELD_OR_KEY,
		})
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM `" + my.table + "` WHERE `name` = 'v1.7.0' LIMIT 1;$").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	err = my.Initiate(my.table)
	st.Assert(t, err, nil)

	logs, err := my.ListLogs()
	st.Assert(t, err, nil)
	st.Assert(t, logs[0].ID, int64(1))

	err = tmpDeployFile(plan.Deploy)
	st.Assert(t, err, nil)
	err = tmpRevertFile(plan.Revert)
	st.Assert(t, err, nil)
	err = tmpVerifyFile(plan.Verify)
	st.Assert(t, err, nil)

	err = my.Deploy(plan)
	st.Assert(t, err, nil)
	err = my.Verify(plan)
	st.Assert(t, err, nil)
	err = my.Revert(plan)
	st.Assert(t, err, nil)

	my.Close()

	_ = os.Remove(plan.Deploy)
	_ = os.Remove(plan.Verify)
	_ = os.Remove(plan.Revert)
}

func tmpDeployFile(file string) error {
	data := []byte(strings.Join([]string{
		"-- Deploy kc:v1.7.0 to mysql",
		"",
		"-- IGNORE 1060",
		"BEGIN;",
		"",
		"ALTER TABLE `test` ADD COLUMN `name` VARCHAR(32) NOT NULL DEFAULT '' AFTER `id`;",
		"",
		"COMMIT;",
		"",
	}, "\n"))

	return os.WriteFile(file, data, os.ModePerm)
}

func tmpRevertFile(file string) error {
	data := []byte(strings.Join([]string{
		"-- Deploy kc:v1.7.0 to mysql",
		"",
		"-- IGNORE 1091",
		"BEGIN;",
		"",
		"ALTER TABLE `test` DROP COLUMN `name`;",
		"",
		"COMMIT;",
		"",
	}, "\n"))

	return os.WriteFile(file, data, os.ModePerm)
}

func tmpVerifyFile(file string) error {
	data := []byte(strings.Join([]string{
		"-- Deploy kc:v1.7.0 to mysql",
		"",
		"BEGIN;",
		"",
		"INSERT INTO `test` (`name`) VALUES ('test');",
		"",
		"ROLLBACK;",
		"",
	}, "\n"))

	return os.WriteFile(file, data, os.ModePerm)
}
