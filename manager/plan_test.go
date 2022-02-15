package manager

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		Dir:     dir,
		Engine:  `mysql`,
		Table:   `dbvm_logs`,
	}
	err = InitProject(project)
	assert.Equal(t, nil, err)

	plan := &Plan{
		Name:     `v1.6.0`,
		Requires: []string{},
		Time:     time.Now(),
		User:     `test`,
		Hostname: `test`,
		Note:     `For Test`,
	}
	err = AddPlan(dir, plan, project.Project, project.Engine)
	assert.Equal(t, nil, err)
}
