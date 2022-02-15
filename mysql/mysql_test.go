package mysql

import (
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VividCortex/mysqlerr"
	driver "github.com/go-sql-driver/mysql"
	"github.com/ideadawn/dbvm/manager"
	"github.com/stretchr/testify/assert"
)

func Test_MySQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Equal(t, nil, err)
	my := New()
	my.db = db
	my.table = `dbvm`

	// Initiate
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

	err = my.Initiate(my.table)
	assert.Equal(t, nil, err)

	logs, err := my.ListLogs()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), logs[0].ID)

	plan := &manager.Plan{
		Name:     `v1.6.0`,
		Requires: []string{},
		Time:     time.Now(),
		Deploy:   `../testdata/deploy/v1.6.0.sql`,
		Revert:   `../testdata/revert/v1.6.0.sql`,
	}

	// Deploy
	mockDeploy(mock)
	err = my.Deploy(plan)
	assert.Equal(t, nil, err)

	// Revert
	mockRevert(mock)
	err = my.Revert(plan)
	assert.Equal(t, nil, err)

	mock.ExpectClose()
	my.Close()
}

func mockDeploy(mock sqlmock.Sqlmock) {
	result := sqlmock.NewResult(0, 0)
	mock.ExpectBegin()
	mock.ExpectExec("^(?is)CREATE TABLE ").WillReturnError(errors.New(`Temp-DB-Error`))
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectExec("^(?is)CREATE TABLE ").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*ADD COLUMN `phone`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*ADD COLUMN `nickname`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*DROP COLUMN `not_exists`").
		WillReturnError(&driver.MySQLError{
			Number:  mysqlerr.ER_CANT_DROP_FIELD_OR_KEY,
			Message: `Column is not exists.`,
		})
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*ADD INDEX `phone`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*ADD PRIMARY").
		WillReturnError(&driver.MySQLError{
			Number:  mysqlerr.ER_MULTIPLE_PRI_KEY,
			Message: `Multi PRIMARY KEY.`,
		})
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*MODIFY COLUMN `id`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*CHANGE COLUMN `phone`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectExec("^CREATE DEFINER=`root`@`localhost` PROCEDURE `delTestById`").WillReturnResult(result)

	mock.ExpectBegin()
	mock.ExpectExec("^CREATE DEFINER=`root`@`localhost` PROCEDURE `delTestByPhone`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func mockRevert(mock sqlmock.Sqlmock) {
	result := sqlmock.NewResult(0, 0)
	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*DROP COLUMN `not_exists`").WillReturnError(errors.New(`Temp-DB-Error`))
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*DROP COLUMN `not_exists`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^(?is)ALTER TABLE `test`.*DROP KEY `phone`").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^DROP TABLE ").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^DROP PROCEDURE ").WillReturnResult(result)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
