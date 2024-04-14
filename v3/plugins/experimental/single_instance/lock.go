package single_instance

import (
	"os"
	"strconv"
)

func GetLockFilePid(filename string) (pid int, err error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	pid, err = strconv.Atoi(string(contents))
	return
}
