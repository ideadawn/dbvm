package manager

import (
	"testing"

	"github.com/nbio/st"
)

func Test_DbUri(t *testing.T) {
	uri := `db:sqlite:`
	_, err := ParseDbUri(uri)
	if err == nil {
		t.Fatal(uri)
	}

	_, err = ParseDbUri(`db:sqlite:/var/sqlite/my.db`)
	st.Assert(t, err, nil)
	params, err := ParseDbUri(`db:mysql://root:qwe:123@127.0.0.1:3306/test?charset=utf8mb4_bin`)
	st.Assert(t, err, nil)
	st.Assert(t, params.Password, `qwe:123`)
	params, err = ParseDbUri(`db:mysql://root:123@456@127.0.0.1:3306/test?charset=utf8mb4_bin#hash`)
	st.Assert(t, err, nil)
	st.Assert(t, params.Password, `123@456`)
	DbUri2Dsn(params)

	params, err = ParseDbUri(`db:mysql://root@127.0.0.1:3306/test#hash`)
	st.Assert(t, err, nil)
	st.Assert(t, params.Password, ``)
	params, err = ParseDbUri(`db:mysql:///test`)
	st.Assert(t, err, nil)
	st.Assert(t, params.Username, ``)
	st.Assert(t, params.Database, `test`)
	st.Assert(t, params.Host, ``)
}
