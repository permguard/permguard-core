// Copyright 2024 Nitro Agility S.r.l.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

// CheckFileIfExists checks if a file exists.
func CheckFileIfExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

// CreateFileIfNotExists creates a file if it does not exist.
func CreateFileIfNotExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return false, nil
	} else if os.IsNotExist(err) {
		dir := filepath.Dir(name)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return false, errors.New("core: failed to create directory")
		}
		file, err := os.Create(name)
		if err != nil {
			return false, errors.New("core: failed to create file")
		}
		defer file.Close()
	} else if os.IsExist(err) {
		return false, nil
	} else {
		return false, errors.New("core: failed to stat file")
	}
	return true, nil
}

// CreateDirIfNotExists creates a directory if it does not exist.
func CreateDirIfNotExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return false, nil
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(name, 0755)
		if err != nil {
			return false, errors.New("core: failed to create directory")
		}
	} else {
		return false, errors.New("core: failed to stat directory")
	}
	return true, nil
}

// WriteFileIfNotExists writes a file if it does not exist.
func WriteFileIfNotExists(name string, data []byte, perm os.FileMode) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return false, nil
	} else if os.IsExist(err) {
		return false, nil
	} else if os.IsNotExist(err) {
		return WriteFile(name, data, perm)
	} else {
		return false, errors.New("core: failed to stat file")
	}
}

// WriteFile writes a file.
func WriteFile(name string, data []byte, perm os.FileMode) (bool, error) {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return false, errors.New("core: failed to write file")
	}
	return true, nil
}

// AppendToFile appends to a file.
func AppendToFile(name string, data []byte) (bool, error) {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, errors.New("core: failed to open file")
	}
	defer file.Close()
	if _, err := file.WriteString(string(data)); err != nil {
		return false, errors.New("core: failed to write to file")
	}
	return true, nil
}

// ReadTOMLFile reads a TOML file.
func ReadTOMLFile(name string, v any) error {
	file, err := os.Open(name)
	if err != nil {
		return errors.New("core: failed to open file")
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return errors.New("core: failed to read file")
	}
	err = toml.Unmarshal(b, v)
	if err != nil {
		return errors.New("core: failed to unmarshal TOML")
	}
	return nil
}

// IsInsideDir checks if a directory is inside another directory.
func IsInsideDir(name string) (bool, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return false, errors.New("core: failed to get current directory")
	}
	for {
		if filepath.Base(currentDir) == name {
			return true, nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return false, nil
}

// ReadIgnoreFile reads an ignore file.
func ReadIgnoreFile(name string) ([]string, error) {
	var ignorePatterns []string
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignorePatterns = append(ignorePatterns, line)
	}
	return ignorePatterns, nil
}

// ShouldIgnore checks if a file should be ignored.
func ShouldIgnore(path string, ignorePatterns []string) bool {
	ignored := false
	for _, pattern := range ignorePatterns {
		isNegation := strings.HasPrefix(pattern, "!")
		pattern = strings.TrimPrefix(pattern, "!")
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			if match == path || strings.HasPrefix(path, match) {
				if isNegation {
					ignored = false
				} else {
					ignored = true
				}
			}
		}
		if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "/")) {
			if isNegation {
				ignored = false
			} else {
				ignored = true
			}
		}
	}
	return ignored
}

// ScanAndFilterFiles scans and filters files.
func ScanAndFilterFiles(rootDir string, exts []string, ignorePatterns []string) ([]string, []string, error) {
	var files []string
	var ignoredFiles []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ShouldIgnore(path, ignorePatterns) {
			ignoredFiles = append(ignoredFiles, path)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			if len(exts) > 0 {
				matched := false
				for _, ext := range exts {
					if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(ext)) {
						matched = true
						break
					}
				}
				if !matched {
					ignoredFiles = append(ignoredFiles, path)
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return files, ignoredFiles, nil
}
