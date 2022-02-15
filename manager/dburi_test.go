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
	assert.Equal(t, nil, err)
	params, err := ParseDbUri(`db:mysql://root:qwe:123@127.0.0.1:3306/test?charset=utf8mb4_bin`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `qwe:123`, params.Password)
	params, err = ParseDbUri(`db:mysql://root:123@456@127.0.0.1:3306/test?charset=utf8mb4_bin#hash`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `123@456`, params.Password)
	DbUri2Dsn(params)

	params, err = ParseDbUri(`db:mysql://root@127.0.0.1:3306/test#hash`)
	assert.Equal(t, nil, err)
	assert.Equal(t, ``, params.Password)
	params, err = ParseDbUri(`db:mysql:///test`)
	assert.Equal(t, nil, err)
	assert.Equal(t, ``, params.Username)
	assert.Equal(t, `test`, params.Database)
	assert.Equal(t, ``, params.Host)
}
