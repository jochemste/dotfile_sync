package libdotfilesync

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/jochemste/dotfile_sync/utils"
)

type FileMap struct {
	Origin      string
	FSPath      string
	Changed     bool
	Message     string
	OriginTime  time.Time
	FSTime      time.Time
	NeedsCommit bool
}

func NewFileMap() *FileMap {
	fm := FileMap{}
	fm.Changed = false
	fm.Message = ""
	fm.NeedsCommit = false
	return &fm
}

func (fm *FileMap) GetOriginFilename() string {
	filename := path.Base(fm.Origin)
	return filename
}

func (fm *FileMap) GetFSFilename() string {
	filename := path.Base(fm.FSPath)
	return filename
}

func (fm *FileMap) Refresh() {
	fm.HasChanged()
	fm.OriginTime, _ = fm.GetModTimeOrigin()
}

// Determine if a file has changed. Returns true if it has or the file does not exist in the local filesystem and/or the memory filesystem
func (fm *FileMap) HasChanged() (bool, error) {
	// If has been run before and file has changed, return the prev result.
	if fm.Changed == true {
		return true, nil
	}

	if fm.ExistsLocal() == false {
		return true, nil
	} else if fm.ExistsInFS() == false {
		return true, nil
	}

	// Get file content from local file
	content1, err := utils.GetFileContent(fm.Origin)
	if err != nil {
		return false, errors.New("Could not get file content " + fm.Origin + ": " + err.Error())
	}
	// Get file content from memory filesystem
	content2, err := GetContentFS(fm.FSPath)
	if err != nil {
		return false, errors.New("Could not get FS file content " + fm.FSPath + ": " + err.Error())
	}

	// Determine if there are differences
	isdiff, err := ContentIsDiff(content1, content2)
	if err != nil {
		return false, errors.New("Could not determine diff " + fm.Origin + " " + fm.FSPath + ": " + err.Error())
	}

	fm.Changed = isdiff

	return isdiff, err
}

func (fm *FileMap) GetModTimeOrigin() (time.Time, error) {
	return utils.GetLastModified(fm.Origin)
}

func (fm *FileMap) ExistsLocal() bool {
	return utils.FileExists(fm.Origin)
}

func (fm *FileMap) ExistsInFS() bool {
	return ExistsInFS(fm.GetOriginFilename())
}

func (fm *FileMap) FMIsMoreRecent() bool {
	if fm.LocalIsMoreRecent() {
		return false
	} else if utils.IsSameTime(fm.OriginTime, fm.FSTime) {
		return false
	}

	return true
}

func (fm *FileMap) LocalIsMoreRecent() bool {
	if utils.IsMoreRecentTime(fm.OriginTime, fm.FSTime) {
		return true
	} else if utils.IsSameTime(fm.OriginTime, fm.FSTime) {
		return false
	}

	return false
}

// Update the Files in the file map
func (fm *FileMap) Update(mode ...string) error {
	// Variables to track push/pull mode
	var push bool
	var pull bool

	if len(mode) == 0 {
		mode = append(mode, "push")
	}

	switch mode[0] {
	case "all":
		fmt.Println("Pushing and Pulling enabled")
		push = true
		pull = true
	case "pull":
		fmt.Println("Pulling enabled")
		pull = true
		push = false
	case "push":
		fmt.Println("Pushing enabled")
		push = true
		pull = false
	default:
		fmt.Println("Mode not recognised, Pulling enabled")
		push = false
		pull = true
	}

	if push == true && fm.ExistsInFS() == false && fm.ExistsLocal() == true {
		// If FS does not contain file, create it in FS (PUSH action)
		fmt.Printf("Update: Writing file in FS: %s\n", fm.FSPath)
		content, err := utils.GetFileContent(fm.Origin)
		if err != nil {
			return errors.New("Could not update " + fm.FSPath + ": " + err.Error())
		}

		fm.NeedsCommit = true
		fm.Message = "Created file " + fm.GetOriginFilename()

		return WriteToFS(fm.GetOriginFilename(), content)
	} else if pull == true && fm.ExistsInFS() == true && fm.ExistsLocal() == false {
		// If local machine does not contain file, create it locally (PULL action)
		fmt.Printf("Update: Writing file in local: %s\n", fm.Origin)
		content, err := GetContentFS(fm.FSPath)
		if err != nil {
			return errors.New("Could not update " + fm.Origin + ": " + err.Error())
		}

		err = utils.CreateFile(fm.Origin)
		if err != nil {
			return errors.New("Could not update " + fm.Origin + ": " + err.Error())
		}

		err = utils.WriteToFile(fm.Origin, content)
		if err != nil {
			return errors.New("Could not update " + fm.Origin + ": " + err.Error())
		}
		return nil
	}

	// Check difference between content
	cOrigin, err := utils.GetFileContent(fm.Origin)
	if err != nil {
		return errors.New("Could not update " + fm.GetOriginFilename() + ": " + err.Error())
	}
	cFS, err := GetContentFS(fm.FSPath)
	isdiff, err := ContentIsDiff(cOrigin, cFS)
	if err != nil {
		return errors.New("Could not update " + fm.GetOriginFilename() + ": " + err.Error())
	}

	// If no difference, no need to update
	if isdiff == false {
		return nil
	}

	fmt.Println(fm.OriginTime)
	fmt.Println(fm.FSTime)

	ogIsMoreRecent := utils.IsMoreRecentTime(fm.OriginTime, fm.FSTime)
	if push == true && ogIsMoreRecent == true {
		// If local file is more recent, update the FS file (PUSH action)
		fmt.Printf("Update: Local is more recent: %s\n", fm.Origin)
		WriteToFS(fm.FSPath, cOrigin)
		fm.NeedsCommit = true
		fm.Message = "Changed " + fm.GetOriginFilename()
	} else if pull == true && !ogIsMoreRecent {
		// If FS file is more recent, update local file
		fmt.Printf("Update: FS is more recent: %s\n", fm.FSPath)
		utils.CopyFile(fm.Origin, fm.Origin+".bak")
		utils.WriteToFile(fm.Origin, cFS)
	}

	return errors.New("Could not update " + fm.GetOriginFilename() + " due to some unexpected state")
}
