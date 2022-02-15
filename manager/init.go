package manager

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectInfo 项目信息
type ProjectInfo struct {
	Version string
	Project string
	URI     string
	Dir     string
	Engine  string
	Table   string
}

// InitProject 初始化项目
func InitProject(project *ProjectInfo) error {
	var err error
	project.Dir, err = correctDir(project.Dir)
	if err != nil {
		return err
	}

	confPath := project.Dir + ConfFile
	planPath := project.Dir + PlanFile

	_, err = os.Stat(confPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return errProjectExists
	}

	err = os.Mkdir(project.Dir+DeployDir, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
	}
	err = os.Mkdir(project.Dir+RevertDir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	data := strings.Join([]string{
		"%syntax-version=" + project.Version,
		"%project=" + project.Project,
		"%uri=" + project.URI,
		"",
	}, "\n")
	err = os.WriteFile(planPath, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}

	bin, err := yaml.Marshal(&Config{
		Engine:    project.Engine,
		FromTable: project.Table,
		LogsTable: project.Table,

		Rule: &Rule{
			Database: &Database{},
			Field:    &Field{},
		},
	})
	if err == nil {
		err = os.WriteFile(confPath, bin, os.ModePerm)
	}
	return err
}
