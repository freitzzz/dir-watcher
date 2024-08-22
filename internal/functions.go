package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Parses a rules.json file.
func Parse(path string) (Rules, error) {
	var rules Rules
	bytes, err := os.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(bytes, &rules)
	}

	return rules, err
}

// Maps the move input from Rules struct to FileExtToDirMap.
func CacheMoveDirectories(move []MoveDir) FileExtToDirMap {
	cache := make(FileExtToDirMap)

	for _, dir := range move {
		for _, mappedExt := range dir.Ext {
			cache[mappedExt] = expand(dir.Path)
		}
	}

	return cache
}

func Watch(rules Rules, cache FileExtToDirMap, watcher *fsnotify.Watcher) error {
	var err error

	for _, dir := range rules.Watch {
		dirPath := expand(Path(dir))
		err = watcher.Add(dirPath)

		if err != nil {
			break
		}

		log.Printf("Watching %s directory\n", dirPath)
	}

	if err == nil {
		go onDirectoryChanged(rules, cache, watcher)
	}

	return err
}

func AutoCleanDir(rules Rules, cache FileExtToDirMap) error {
	var err error

	for _, dir := range rules.Watch {
		dirPath := expand(Path(dir))
		err = cleanDir(dirPath, cache, expand(rules.Unknown))

		if err != nil {
			break
		}

		log.Printf("Cleaned %s directory.\n", dirPath)
	}

	return err
}

func onDirectoryChanged(rules Rules, cache FileExtToDirMap, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			filePath := event.Name

			if event.Has(fsnotify.Chmod) && !shouldIgnoreFile(filePath) {
				moveDirPath := cache[ext(filePath)]

				if len(moveDirPath) == 0 {
					moveDirPath = expand(rules.Unknown)
				}

				err := move(filePath, moveDirPath)

				if err != nil {
					log.Printf("Failed to move file %s\n %v", filePath, err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func move(srcPath string, destPath string) error {
	var err error

	if !exists(destPath) {
		err = os.MkdirAll(destPath, os.ModePerm)
	}

	if err == nil {
		moveDestPath := filepath.Join(destPath, filepath.Base(srcPath))

		if exists(moveDestPath) {
			moveDestPath = uniqueFilePath(moveDestPath)
		}

		err = os.Rename(srcPath, moveDestPath)

		if err == nil {
			log.Printf("Moved file from %s to %s\n", srcPath, moveDestPath)
		}
	}

	return err
}

// Gets the next available file path for a file that already exists on the file system.
func uniqueFilePath(path string) string {
	ctr := 1

	uniqueFilePath := path

	for {
		filePath := filepath.Base(path)
		fileExt := filepath.Ext(filePath)
		dirName := strings.TrimSuffix(path, filePath)
		uniqueFilePath = filepath.Join(dirName, fmt.Sprintf("%s-%d%s", strings.TrimSuffix(filePath, fileExt), ctr, fileExt))

		if !exists(uniqueFilePath) {
			break
		}

		ctr += 1
	}

	return uniqueFilePath
}

// Checks if a file path exists on the filesystem.
func exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil || !os.IsNotExist(err)
}

// Expands a path if it contains short links like the ~ for home directory.
func expand(path Path) string {
	if path[0] != '~' {
		return string(path)
	}

	return os.ExpandEnv(fmt.Sprintf("%s%s", "$HOME", path[1:]))
}

// Extracts the file extension from a file path. The returning value does not include the dot (.).
func ext(path string) string {
	ext := filepath.Ext(path)

	if len(ext) > 1 {
		ext = ext[1:]
	}

	return strings.ToLower(ext)
}

func shouldIgnoreFile(path string) bool {
	fileName := filepath.Base(path)

	return strings.HasPrefix(fileName, ".com.google.Chrome") || strings.HasSuffix(fileName, ".crdownload")
}

// Cleans a directory by moving containing files to directories recognized by the files extension.
func cleanDir(dirPath string, cache FileExtToDirMap, fallbackDirPath string) error {
	entries, err := os.ReadDir(dirPath)

	if err == nil {
		for _, file := range entries {
			if !file.IsDir() || shouldIgnoreFile(file.Name()) {
				movingDirPath := cache[ext(file.Name())]

				if len(movingDirPath) == 0 {
					movingDirPath = fallbackDirPath
				}

				err = move(filepath.Join(dirPath, file.Name()), movingDirPath)

				if err != nil {
					break
				}
			}
		}
	}

	return err
}
