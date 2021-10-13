package manager

import (
	"errors"
	"fmt"
	"os"
)

// Revert 回退版本
func (m *Manager) Revert(to string) error {
	logs, err := m.engine.ListLogs()
	if err != nil {
		return err
	}

	var arr []*Plan

	for idx, plan := range m.plans {
		for _, log := range logs {
			if log.Name == plan.Name {
				plan.deployed = true
				m.deployed[log.Name] = true
				break
			}
		}

		if to == plan.Name {
			if !plan.deployed {
				return fmt.Errorf("Version( %s ) not deployed.\n", to)
			}

			arr = m.plans[idx:]
		}
	}

	lth := len(arr)
	if lth == 0 {
		return errors.New(`Version not found: ` + to)
	}

	for lth -= 1; lth > -1; lth-- {
		if !arr[lth].deployed {
			continue
		}

		info, err := os.Stat(arr[lth].Revert)
		if err != nil || info.IsDir() {
			return fmt.Errorf("Revert-File ( %s ) check failed: %s\n", arr[lth].Revert, err.Error())
		}
		err = m.engine.Revert(arr[lth])
		if err != nil {
			return err
		}

		fmt.Println(arr[lth].Name, `reverted.`)
	}

	return nil
}
