package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myEngine struct{}

func (m *myEngine) Connect(*Params) error {
	return nil
}
func (m *myEngine) Close() {}
func (m *myEngine) Initiate(string) error {
	return nil
}
func (m *myEngine) ListLogs() ([]*Log, error) {
	return []*Log{
		&Log{
			Name: `v1.6.0`,
		},
	}, nil
}
func (m *myEngine) Deploy(*Plan) error {
	return nil
}
func (m *myEngine) Revert(*Plan) error {
	return nil
}

func Test_Manager(t *testing.T) {
	RegisterEngine(`mysql`, &myEngine{})

	mgr, err := New(`../testdata`, `db:mysql://root:qwe123@127.0.0.1:3306/test?charset=utf8mb4`)
	assert.Equal(t, nil, err)

	_ = mgr.GetLogsTable()
	err = mgr.Deploy(`latest`)
	assert.Equal(t, nil, err)

	err = mgr.Revert(`v1.7.0`)
	assert.Equal(t, nil, err)

	mgr.Close()
}
