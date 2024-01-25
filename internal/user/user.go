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
	return err == nil && currUser.Username == "root"
}

func isAdministrator() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err != nil
}
