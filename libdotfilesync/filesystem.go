package libdotfilesync

import (
	"errors"
	//"fmt"
	"io"
	"os"

	//"github.com/jochemste/dotfile_sync/utils"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileSystem struct {
	FS billy.Filesystem
}

var Filesystem FileSystem

/*
var filesystem billy.Filesystem
var pFS *billy.Filesystem
*/

func NewFS() {
	Filesystem = FileSystem{}
	Filesystem.FS = memfs.New()
}

// Find in filesystem
func FindInFS(filename string) (string, error) {
	// Read root directory of memory file system
	files, err := Filesystem.FS.ReadDir("/")

	if err != nil {
		return "", errors.New("Could not find file: " + err.Error())
	}

	res, err := FindInFSDir(filename, files, "")
	if err != nil {
		return "", errors.New("Could not find file: " + err.Error())
	}

	return res, nil
}

// Recursive find in directory function (private)
func FindInFSDir(filename string, files []os.FileInfo, dirname string) (string, error) {

	for _, file := range files {
		if file.Name() == filename {
			return dirname + "/" + file.Name(), nil
		}

		if file.IsDir() == true {
			f, err := Filesystem.FS.ReadDir(file.Name())
			if err != nil {
				return "", err
			}
			res, err := FindInFSDir(filename, f, dirname+"/"+file.Name())
			if res != "" {
				return res, err
			}
		}
	}

	return "", errors.New("File was not found")
}

/*
func FileToFS(filename string) error {
	if !utils.FileExists(filename) {
		return errors.New("Could not write file to FS, File " + filename + " does not exist")
	}

	filesystem.Create(filename)

	return nil
  }
*/

func WriteToFS(filename string, content string) error {
	newfd, err := Filesystem.FS.Create(filename)
	if err != nil {
		errors.New("Could not write to FS: " + err.Error())
	}

	newfd.Write([]byte(content))
	newfd.Close()

	return nil
}

func GetFileDescriptorFS(filename string) (billy.File, error) {
	var fd billy.File
	path, err := FindInFS(filename)
	if path == "" && err != nil {
		fd, err = Filesystem.FS.Create(filename)
	}
	return fd, err
}

func GetContentFS(filename string) (string, error) {
	file, err := Filesystem.FS.Open(filename)
	if err != nil {
		return "", errors.New("Could not read file contents from FS: " + err.Error())
	}
	defer file.Close()

	const chunk = 1024
	buffer := make([]byte, chunk)
	content := ""

	for {
		length, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.New("Could not read file contents from FS: " + err.Error())
		}

		if length > 0 {
			content += string(buffer[:length])
		}
	}

	return content, nil
}

func ContentIsDiff(cont1 string, cont2 string) (bool, error) {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(cont1, cont2, false)
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			return true, nil
		}
	}
	return false, nil
}

func ExistsInFS(filename string) bool {
	path, err := FindInFS(filename)
	if path == "" && err != nil {
		return false
	}

	return true
}
