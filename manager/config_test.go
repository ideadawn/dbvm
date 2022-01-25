package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig(t *testing.T) {
	conf, err := ParseConfig(`../testdata`)
	assert.Equal(t, err, nil)
	assert.Equal(t, conf.Engine, `mysql`)
}
