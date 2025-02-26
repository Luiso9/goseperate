package utils

import "os"

func RemoveTempFiles(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}
