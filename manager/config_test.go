package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig(t *testing.T) {
	conf, err := ParseConfig(`../testdata`)
	assert.Equal(t, nil, err)
	assert.Equal(t, "TEXT", conf.Rule.Field.Excepts[0])
}
