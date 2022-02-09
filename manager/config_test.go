package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig(t *testing.T) {
	conf, err := ParseConfig(`../testdata`)
	assert.Equal(t, err, nil)
	assert.Equal(t, conf.Rule.Field.Excepts[0], "TEXT")
}
