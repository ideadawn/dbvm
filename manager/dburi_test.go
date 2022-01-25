package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DbUri(t *testing.T) {
	uri := `db:sqlite:`
	_, err := ParseDbUri(uri)
	if err == nil {
		t.Fatal(uri)
	}

	_, err = ParseDbUri(`db:sqlite:/var/sqlite/my.db`)
	assert.Equal(t, err, nil)
	params, err := ParseDbUri(`db:mysql://root:qwe:123@127.0.0.1:3306/test?charset=utf8mb4_bin`)
	assert.Equal(t, err, nil)
	assert.Equal(t, params.Password, `qwe:123`)
	params, err = ParseDbUri(`db:mysql://root:123@456@127.0.0.1:3306/test?charset=utf8mb4_bin#hash`)
	assert.Equal(t, err, nil)
	assert.Equal(t, params.Password, `123@456`)
	DbUri2Dsn(params)

	params, err = ParseDbUri(`db:mysql://root@127.0.0.1:3306/test#hash`)
	assert.Equal(t, err, nil)
	assert.Equal(t, params.Password, ``)
	params, err = ParseDbUri(`db:mysql:///test`)
	assert.Equal(t, err, nil)
	assert.Equal(t, params.Username, ``)
	assert.Equal(t, params.Database, `test`)
	assert.Equal(t, params.Host, ``)
}
