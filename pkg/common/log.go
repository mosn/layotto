package common

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
)

func GetLogPath(fileName string) string {
	var logFolder string
	if u, err := user.Current(); err != nil {
		logFolder = "/home/admin/logs/mosn"
	} else if runtime.GOOS == "darwin" {
		logFolder = fmt.Sprintf(path.Join(u.HomeDir, "logs/mosn"))
	} else if runtime.GOOS == "windows" {
		logFolder = fmt.Sprintf(path.Join(u.HomeDir, "logs/mosn"))
	} else {
		logFolder = "/home/admin/logs/mosn"
	}

	logPath := logFolder + string(os.PathSeparator) + fileName
	return logPath
}
