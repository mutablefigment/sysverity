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
	"os"
	"runtime"
)

const mbrSize = 1024

/* This function dumps the MBR */
func UniversalMBRDump() ([]byte, error) {

	var mbrBytes = make([]byte, mbrSize)
	var err error = nil
	var diskPath string

	// Set boot drive depending on os
	switch runtime.GOOS {
	case "darwin":
		// FIXME: how to dump MBR on MacOs?
		fallthrough
	case "linux":
		diskPath = "/dev/sda"
	case "windows":
		diskPath = "\\.\\PhysicalDrive0"
	}

	// Open the boot drive
	mbrHandle, err := os.Open(diskPath)
	if err != nil {
		return nil, err
	}

	// Read the first 1024 bytes
	_, err = mbrHandle.Read(mbrBytes)
	if err != nil {
		return nil, err
	}

	return mbrBytes, err
}
