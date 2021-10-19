package main

import (
	"os"
	"testing"

	"github.com/ideadawn/dbvm/manager"
	"github.com/nbio/st"
)

type myEngine struct{}

func (m *myEngine) Connect(*manager.Params) error {
	return nil
}
func (m *myEngine) Close() {}
func (m *myEngine) Initiate(string) error {
	return nil
}
func (m *myEngine) ListLogs() ([]*manager.Log, error) {
	return []*manager.Log{}, nil
}
func (m *myEngine) Deploy(*manager.Plan) error {
	return nil
}
func (m *myEngine) Verify(*manager.Plan) error {
	return nil
}
func (m *myEngine) Revert(*manager.Plan) error {
	return nil
}

type myEngine2 struct {
	myEngine
}

func (m *myEngine2) ListLogs() ([]*manager.Log, error) {
	return []*manager.Log{
		&manager.Log{
			Name: `v1.0.0`,
		},
	}, nil
}

func Test_DBVM(t *testing.T) {
	manager.RegisterEngine(`mysql`, &myEngine{})

	dir := `./tmp-project`
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			t.Fatal(err)
		}
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	os.Args = []string{
		`dbvm`,
		`init`,
		`--project`, `test`,
		`--engine`, `mysql`,
		`--dir`, dir,
		`--set`, `logsTable=logs`,
	}
	err = cmdInit()
	st.Assert(t, err, nil)

	os.Args = []string{
		`dbvm`,
		`add`,
		`--name`, `v1.0.0`,
		`--dir`, dir,
		`--note`, `for test`,
		`--user`, `test`,
	}
	err = cmdAdd()
	st.Assert(t, err, nil)

	os.Args = []string{
		`dbvm`,
		`deploy`,
		`--dir`, dir,
		`--uri`, `db:mysql://root:123@127.0.0.1:3306/test`,
		`--to`, `latest`,
	}
	err = cmdDeploy()
	st.Assert(t, err, nil)

	manager.RegisterEngine(`mysql`, &myEngine2{})

	os.Args = []string{
		`dbvm`,
		`verify`,
		`--dir`, dir,
		`--uri`, `db:mysql://root:123@127.0.0.1:3306/test`,
		`--from`, `v1.0.0`,
	}
	err = cmdVerify()
	st.Assert(t, err, nil)

	os.Args = []string{
		`dbvm`,
		`verify`,
		`--dir`, dir,
		`--uri`, `db:mysql://root:123@127.0.0.1:3306/test`,
		`--to`, `v1.0.0`,
	}
	err = cmdVerify()
	st.Assert(t, err, nil)

	os.Args = []string{
		`dbvm`,
		`verify`,
		`--dir`, dir,
		`--uri`, `db:mysql://root:123@127.0.0.1:3306/test`,
		`--from`, `v1.0.0`,
		`--to`, `v1.0.0`,
	}
	err = cmdVerify()
	st.Assert(t, err, nil)

	os.Args = []string{
		`dbvm`,
		`revert`,
		`--dir`, dir,
		`--uri`, `db:mysql://root:123@127.0.0.1:3306/test`,
		`--to`, `v1.0.0`,
	}
	err = cmdRevert()
	st.Assert(t, err, nil)
}
