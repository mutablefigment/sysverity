package utils

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slog"
)

const mbrDumpPath = "/dev/shm/mbr.bak"

func DumpMBR(logger *slog.Logger, p *tea.Program) (string, error) {

	if _, err := os.Stat(mbrDumpPath); err == nil {
		return "", err
	} else if os.IsNotExist(err) {
		// Backup the mbr
		// FIXME: how can I use another command and also a config file here to make disks more dynamic
		cmd := exec.Command("dd", "if=/dev/sda", fmt.Sprintf("of=%s", mbrDumpPath), "bs=1024k", "count=1") // "status=none"
		stdout, err := cmd.Output()
		if err != nil {
			//p.Send(resultMsg{msg: err.Error(), level: Error})
			// logger.Error(err.Error())
			return "", err
		}

		return string(stdout), err
	} else {
		return "", fmt.Errorf("failed to check if file exits! %s", mbrDumpPath)
	}

}
