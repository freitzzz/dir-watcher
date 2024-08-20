package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Parse(path string) (Rules, error) {
	var rules Rules
	bytes, err := os.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(bytes, &rules)
	}

	return rules, err
}

func Cache(move []MoveDir) map[string]string {
	cache := make(map[string]string)

	for _, m := range move {
		for _, ext := range m.Ext {
			cache[ext] = Expand(m.Path)
		}
	}

	return cache
}

func Move(src string, dest string) error {
	var err error

	if !Exists(dest) {
		err = os.MkdirAll(dest, os.ModePerm)
	}

	if err == nil {
		dp := filepath.Join(dest, filepath.Base(src))

		if Exists(dp) {
			dp = FindNextAvailableFilepath(dp)
		}

		println(src)
		println(dp)

		err = os.Rename(src, dp)
	}

	return err
}

func FindNextAvailableFilepath(path string) string {
	ctr := 1

	npath := path

	for {
		filePath := filepath.Base(path)
		fileExt := filepath.Ext(filePath)
		dirName := strings.TrimSuffix(path, filePath)
		npath = filepath.Join(dirName, fmt.Sprintf("%s-%d%s", strings.TrimSuffix(filePath, fileExt), ctr, fileExt))

		if !Exists(npath) {
			break
		}
	}

	return npath
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil || !os.IsNotExist(err)
}

func Expand(gp GlobPath) string {
	if gp[0] != '~' {
		return string(gp)
	}

	return os.ExpandEnv(fmt.Sprintf("%s%s", "$HOME", gp[1:]))
}

func Ext(path string) string {
	ext := filepath.Ext(path)

	if len(ext) > 1 {
		ext = ext[1:]
	}

	return strings.ToLower(ext)
}

func ShouldIgnoreFile(path string) bool {
	fileName := filepath.Base(path)

	return strings.HasPrefix(fileName, ".com.google.Chrome") || strings.HasSuffix(fileName, ".crdownload")
}

func CleanDir(dirPath string, cache map[string]string, unknown string) error {
	entries, err := os.ReadDir(dirPath)

	if err == nil {
		for _, f := range entries {
			if !f.IsDir() || ShouldIgnoreFile(f.Name()) {
				mp := cache[Ext(f.Name())]

				if len(mp) == 0 {
					mp = unknown
				}

				err = Move(filepath.Join(dirPath, f.Name()), mp)

				if err != nil {
					break
				}
			}
		}
	}

	return err
}
