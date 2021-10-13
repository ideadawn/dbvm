package manager

import (
	"errors"
	"fmt"
	"os"
)

// Deploy 部署到指定版本
func (m *Manager) Deploy(to string) error {
	logs, err := m.engine.ListLogs()
	if err != nil {
		return err
	}

	if to == `latest` {
		lth := len(m.plans)
		if lth > 0 {
			lth--
			to = m.plans[lth].Name
		}
	}
	notFound := true

	for _, plan := range m.plans {
		for _, log := range logs {
			if log.Name == plan.Name {
				plan.deployed = true
				m.deployed[log.Name] = true
				break
			}
		}

		if to == plan.Name {
			if plan.deployed {
				fmt.Printf("Version( %s ) was deployed.\n", to)
				return nil
			}
			notFound = false
		}
	}

	if notFound {
		return errors.New(`Version not found: ` + to)
	}

	for _, plan := range m.plans {
		if plan.deployed {
			continue
		}
		for _, req := range plan.Requires {
			if _, ok := m.deployed[req]; !ok {
				return fmt.Errorf("Deploy %s need %s , check your plans.\n", plan.Name, req)
			}
		}

		info, err := os.Stat(plan.Deploy)
		if err != nil || info.IsDir() {
			return fmt.Errorf("Deploy-File ( %s ) check failed: %s\n", plan.Deploy, err.Error())
		}
		info, err = os.Stat(plan.Verify)
		if err != nil || info.IsDir() {
			return fmt.Errorf("Verify-File ( %s ) check failed: %s\n", plan.Verify, err.Error())
		}
		info, err = os.Stat(plan.Revert)
		if err != nil || info.IsDir() {
			return fmt.Errorf("Revert-File ( %s ) check failed: %s\n", plan.Revert, err.Error())
		}

		err = m.engine.Deploy(plan)
		if err != nil {
			e2 := m.engine.Revert(plan)
			if e2 != nil {
				fmt.Println(`revert error:`, e2)
			}
			return err
		}

		err = m.engine.Verify(plan)
		if err != nil {
			e2 := m.engine.Revert(plan)
			if e2 != nil {
				fmt.Println(`revert error:`, e2)
			}
			return err
		}

		plan.deployed = true
		m.deployed[plan.Name] = true
		fmt.Println(plan.Name, `deployed.`)

		if to == plan.Name {
			break
		}
	}

	return nil
}
