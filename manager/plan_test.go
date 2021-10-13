package manager

import (
	"os"
	"testing"
	"time"

	"github.com/nbio/st"
)

func Test_AddPlan(t *testing.T) {
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

	project := &ProjectInfo{
		Version: `1.0.0`,
		Project: `test`,
		URI:     ``,
		Engine:  `mysql`,
		Dir:     dir,
		Set: []string{
			`logsTable=dbvm`,
		},
	}
	err = InitProject(project)
	st.Assert(t, err, nil)

	plan := &Plan{
		Name:     `v1.6.0`,
		Requires: []string{},
		Time:     time.Now(),
		User:     `test`,
		Hostname: `test`,
		Note:     `For Test`,
	}
	err = AddPlan(dir, plan, project.Project, project.Engine)
	st.Assert(t, err, nil)
}
