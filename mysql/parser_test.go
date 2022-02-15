package mysql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParserIssue(t *testing.T) {
	dir := `../testdata/bad/`
	p := &sqlParser{
		file: dir + `create_table.sql`,
	}
	p.parseSqlBlocks()
	assert.Equal(t, errCreateTableINE, p.err)

	p.reset(dir + `drop_table.sql`)
	p.parseSqlBlocks()
	assert.Equal(t, errDropTableIE, p.err)

	p.reset(dir + `ignore.sql`)
	p.parseSqlBlocks()
	assert.Equal(t, true, p.err != nil)
	assert.Equal(t, true, strings.Contains(p.sql, `DUP_ENTRY`))

	p.reset(dir + `alter.sql`)
	p.parseSqlBlocks()
	assert.Equal(t, errAlterUnknown, p.err)
}

func Test_ParserDeploy(t *testing.T) {
	dir := `../testdata/deploy/`
	p := &sqlParser{
		file: dir + `v1.6.0.sql`,
	}
	p.parseSqlBlocks()
	assert.Equal(t, nil, p.err)
	fmt.Println("")
	p.print()
}

func Test_ParserRevert(t *testing.T) {
	dir := `../testdata/revert/`
	p := &sqlParser{
		file: dir + `v1.6.0.sql`,
	}
	p.parseSqlBlocks()
	assert.Equal(t, nil, p.err)
	fmt.Println("")
	p.print()
}
