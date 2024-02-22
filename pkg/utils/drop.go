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
	"os/user"
	"strconv"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slog"
)

/*
This function drops privileges to the user that was
specified as an argument.
*/
func DropPrivileges(userToSwitchTo string, logger *slog.Logger, p *tea.Program) error {

	// Lookup the user
	userInfo, err := user.Lookup(userToSwitchTo)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Get the groupID
	gid, err := strconv.Atoi(userInfo.Gid)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Get the userID
	uid, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Unset supplementray group ids
	err = syscall.Setgroups([]int{})
	if err != nil {
		logger.Error("Failed to unset supplementary group IDs: " + err.Error())
		return err
	}

	// Set group ID
	err = syscall.Setgid(gid)
	if err != nil {
		logger.Error("Failed to set group ID: " + err.Error())
		return err
	}

	// Set the user id
	err = syscall.Setuid(uid)
	if err != nil {
		logger.Error("Failed to set user ID:" + err.Error())
		return err
	}

	return nil
}
