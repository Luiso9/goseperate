package utils

import "os"

// RemoveTempFiles deletes temp files after processing
func RemoveTempFiles(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}
