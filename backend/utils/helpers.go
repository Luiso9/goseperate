package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	fileTimestamps = make(map[string]time.Time)
	fileMapMutex   sync.Mutex
)

func RemoveTempFiles(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}

func FindExistingFile(dir, filename string) (string, error) {
	pattern := fmt.Sprintf("%s/%s.*", dir, filename)
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("file not found")
	}
	return matches[0], nil
}

func TrackFileForCleanup(filePath string) {
	fileMapMutex.Lock()
	defer fileMapMutex.Unlock()
	
	fileTimestamps[filePath] = time.Now()
}

func StartCleanupRoutine(cleanupInterval, fileLifetime time.Duration) {
	if cleanupInterval <= 0 {
		cleanupInterval = 1 * time.Minute
	}
	if fileLifetime <= 0 {
		fileLifetime = 10 * time.Minute
	}
	
	go func() {
		for {
			time.Sleep(cleanupInterval)
			CleanupOldFiles(fileLifetime)
		}
	}()
}

func CleanupOldFiles(maxAge time.Duration) {
	fileMapMutex.Lock()
	defer fileMapMutex.Unlock()
	
	now := time.Now()
	cutoffTime := now.Add(-maxAge)
	
	for path, timestamp := range fileTimestamps {
		if timestamp.Before(cutoffTime) {
			err := os.Remove(path)
			if err == nil {
				delete(fileTimestamps, path)
			}
		}
	}
}

func CleanupDirectory(dirPath string, maxAge time.Duration) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	
	now := time.Now()
	cutoffTime := now.Add(-maxAge)
	
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			if info.ModTime().Before(cutoffTime) {
				os.Remove(filePath)
			}
		}
	}
	
	return nil
}