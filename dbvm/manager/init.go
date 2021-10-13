package manager

import (
	"fmt"
	"os"
	"strings"
)

// ProjectInfo 项目信息
type ProjectInfo struct {
	Version string
	Project string
	URI     string
	Engine  string
	Dir     string
	Set     []string
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
	err = os.Mkdir(project.Dir+VerifyDir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
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

	data = strings.Join([]string{
		"[core]",
		"	engine = " + project.Engine,
		"",
	}, "\n")
	if len(project.Set) > 0 {
		data += "[core \"variables\"]\n"
		for _, set := range project.Set {
			pos := strings.Index(set, "=")
			if pos > 0 {
				data += fmt.Sprintf("	%s = %s\n", set[0:pos], set[pos+1:])
			}
		}
	}

	err = os.WriteFile(confPath, []byte(data), os.ModePerm)
	return err
}
