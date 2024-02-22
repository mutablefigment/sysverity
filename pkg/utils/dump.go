/*
sysverity
Copyright (C) 2024  mutablefigment

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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
