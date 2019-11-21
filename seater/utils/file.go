package utils

import (
	"os"
)

// FileExists checks if file exists
func FileExists(path string) (exist bool, err error) {
	if _, err = os.Stat(path); err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
