package manager

import (
	"testing"

	"github.com/nbio/st"
)

func Test_ParseConfig(t *testing.T) {
	conf, err := ParseConfig(`../sqitch`)
	st.Assert(t, err, nil)
	st.Assert(t, conf.Engine, `mysql`)
}
