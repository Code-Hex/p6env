package util

import "os"

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
