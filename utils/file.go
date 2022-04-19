package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type PathType int64

const (
	NONE    PathType = -1
	DIR     PathType = 0
	FILE    PathType = 1
	CONFDIR string   = "/etc/dotfilesync"
)

func SetDefPermission(path string) error {
	exists, _, err := PathExists(path)
	if err != nil {
		return err
	}
	if exists == false {
		return errors.New("Path does not exist: " + path)
	}
	os.Chmod(path, 0755)
	return nil
}

func Mkdir(path string) error {
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		return errors.New("Could not create dir " + path + ": " + err.Error())
	}

	return nil
}

//Get only name of a file
func DirFileSplit(path string) (string, string) {
	dir, file := filepath.Split(path)

	return dir, file
}

func PathExists(path string) (bool, PathType, error) {
	info, err := os.Stat(path)
	exists := !errors.Is(err, os.ErrNotExist)

	if !exists {
		return exists, NONE, err
	} else {
		if info.IsDir() {
			return true, DIR, nil
		} else {
			return true, FILE, nil
		}
	}
	return false, NONE, err
}

// Check if a file exist
func FileExists(path string) bool {
	exists, ptype, _ := PathExists(path)

	if exists == true && ptype == FILE {
		return true
	} else {
		return false
	}
}

// Check if a dir exists
func DirExists(path string) bool {
	exists, ptype, _ := PathExists(path)

	if exists == true && ptype == DIR {
		return true
	} else {
		return false
	}
}

// Get the contents of a file (non-buffered, so not ideal for large files)
func GetFileContent(path string) (string, error) {
	if !FileExists(path) {
		return "", errors.New("File does not exist")
	}

	content, err := ioutil.ReadFile(path)
	CheckIfError(err)

	return string(content), nil
}

// Get the ModTime of a file
func GetLastModified(path string) (time.Time, error) {
	if !FileExists(path) {
		return time.Now(), errors.New("File does not exist")
	}
	fStat, err := os.Stat(path)
	CheckIfError(err)
	return fStat.ModTime(), nil
}

// Compare the most recent ModTime of two files and determine if "path1" is more recent or not
func IsMoreRecent(path1 string, path2 string) (bool, error) {
	if !FileExists(path1) || !FileExists(path2) {
		return false, errors.New("File does not exist")
	}

	t1, _ := GetLastModified(path1)
	t2, _ := GetLastModified(path2)

	delta := t1.Sub(t2)

	if delta > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Compare two files and calculate the difference. If a difference exists, return true,
// otherwise return false
func IsDiff(path1 string, path2 string) (bool, error) {
	if !FileExists(path1) || !FileExists(path2) {
		return true, errors.New("File does not exist")
	}

	dmp := diffmatchpatch.New()
	t1, err := GetFileContent(path1)
	CheckIfError(err)
	t2, err := GetFileContent(path2)
	CheckIfError(err)

	diffs := dmp.DiffMain(t1, t2, false)
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			return true, nil
		}
	}
	return false, nil
}
