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
