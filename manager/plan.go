package manager

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Plan 执行计划
type Plan struct {
	Name     string
	Requires []string
	Time     time.Time
	User     string
	Hostname string
	Note     string
	Deploy   string
	Verify   string
	Revert   string

	deployed bool
}

// 解析执行计划列表
func ParsePlan(dir string) (map[string]string, []*Plan, error) {
	dir, err := correctDir(dir)
	if err != nil {
		return nil, nil, err
	}

	data, err := os.ReadFile(dir + PlanFile)
	if err != nil {
		return nil, nil, err
	}

	var list []*Plan
	env := map[string]string{}

	re_plan := regexp.MustCompile(`^([^ ]+) (?:\[(.+)\] )?([0-9\-T:Z]{20}) \w+ <.*> #(.*)$`)
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		line = bytes.Trim(line, " \r\n\t")
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		if line[0] == '%' {
			pos := bytes.Index(line, []byte{'='})
			if pos > 0 {
				env[string(line[1:pos])] = string(line[pos+1:])
			}
			continue
		}

		matches := re_plan.FindAllSubmatch(line, -1)
		if len(matches) == 0 {
			continue
		}

		name := string(matches[0][1])
		plan := &Plan{
			Name:   name,
			Note:   string(matches[0][4]),
			Deploy: strings.Join([]string{dir, DeployDir, `/`, name, `.sql`}, ``),
			Revert: strings.Join([]string{dir, RevertDir, `/`, name, `.sql`}, ``),
		}
		plan.Time, err = time.ParseInLocation(`2006-01-02T15:04:05Z`, string(matches[0][3]), time.UTC)
		if err == nil {
			plan.Time = plan.Time.Local()
		} else {
			plan.Time = time.Now()
		}
		if len(matches[0][2]) > 0 {
			plan.Requires = strings.Split(string(matches[0][2]), ` `)
		}

		list = append(list, plan)
	}

	return env, list, nil
}

// AddPlan 添加部署计划
func AddPlan(dir string, plan *Plan, project, engine string) error {
	dir, err := correctDir(dir)
	if err != nil {
		return err
	}

	plan.Deploy = strings.Join([]string{dir, DeployDir, `/`, plan.Name, `.sql`}, ``)
	plan.Revert = strings.Join([]string{dir, RevertDir, `/`, plan.Name, `.sql`}, ``)

	data := strings.Join([]string{
		fmt.Sprintf("-- Deploy %s:%s to %s", project, plan.Name, engine),
		"",
		"BEGIN;",
		"",
		"-- add deploy sql at here...",
		"",
		"COMMIT;",
		"",
	}, "\n")
	err = os.WriteFile(plan.Deploy, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}

	data = strings.Join([]string{
		fmt.Sprintf("-- Revert %s:%s from %s", project, plan.Name, engine),
		"",
		"BEGIN;",
		"",
		"-- add revert sql at here...",
		"",
		"COMMIT;",
		"",
	}, "\n")
	err = os.WriteFile(plan.Revert, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}

	planItem := plan.Name
	if len(plan.Requires) > 0 {
		planItem += strings.Join([]string{
			` [`,
			strings.Join(plan.Requires, ` `),
			`]`,
		}, ``)
	}
	planItem += ` `
	planItem += strings.Join([]string{
		plan.Time.UTC().Format(`2006-01-02T15:04:05Z`),
		plan.User,
		plan.Hostname,
		`#` + plan.Note,
	}, ` `)

	planPath := dir + PlanFile
	content, err := os.ReadFile(planPath)
	if err != nil {
		return err
	}

	content = bytes.TrimSpace(content)
	cttArr := bytes.Split(content, []byte{'\n'})
	last := len(cttArr) - 1
	if last < 0 {
		return errProjectNotInit
	}

	line := bytes.TrimSpace(cttArr[last])
	if len(line) < 1 {
		return errProjectNotInit
	}
	if line[0] == '%' {
		content = append(content, '\n', '\n')
	} else {
		content = append(content, '\n')
	}
	content = append(content, []byte(planItem)...)
	content = append(content, '\n')

	err = os.WriteFile(planPath, content, os.ModePerm)
	return err
}
