package user

import (
	"os"
	"os/user"
	"runtime"
)

func IsPrivileged() bool {
	if isLinux() {
		return isRoot()
	}

	return isAdministrator()

}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func isRoot() bool {
	currUser, err := user.Current()
	if err != nil {
		return false
	}

	return currUser.Username == "root"
}

func isAdministrator() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}
