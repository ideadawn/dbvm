package manager

/*
import (
	"errors"
	"fmt"
	"os"
)

// Verify 校验部署
func (m *Manager) Verify(to string) error {
	logs, err := m.engine.ListLogs()
	if err != nil {
		return err
	}

	var dest *Plan

	for _, plan := range m.plans {
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
			dest = plan
			break
		}
	}

	if dest == nil {
		return errors.New(`Version not found: ` + to)
	}

	info, err := os.Stat(dest.Verify)
	if err != nil || info.IsDir() {
		return fmt.Errorf("Verify-File ( %s ) check failed: %s\n", dest.Verify, err.Error())
	}

	return m.engine.Verify(dest)
}
*/
