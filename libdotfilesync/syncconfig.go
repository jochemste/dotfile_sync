package libdotfilesync

import (
	"errors"
	"fmt"
	//	"os"
	"time"

	//"github.com/jochemste/dotfile_sync/utils"

	//"github.com/go-git/go-billy/v5"

	"github.com/pelletier/go-toml/v2"
)

var Syncfilepath string = ".syncinfo"
var SyncChanged bool = false

type SyncFile struct {
	Name     string
	LastSync time.Time
}

type SyncConfig struct {
	LastSync time.Time
	Files    []SyncFile
}

func (sc *SyncConfig) ToFS() error {
	fd, err := Filesystem.FS.Create(Syncfilepath)
	if err != nil {
		return errors.New("Could not create file " + Syncfilepath)
	}

	if err := toml.NewEncoder(fd).Encode(sc); err != nil {
		return err
	}

	if err := fd.Close(); err != nil {
		return err
	}

	return nil
}

func (sc *SyncConfig) FromFS() error {
	_, err := GetFileDescriptorFS(Syncfilepath)
	content, err := GetContentFS(Syncfilepath)

	if err != nil {
		return err
	}

	if err := toml.Unmarshal([]byte(content), &sc); err != nil {
		return err
	}

	return nil
}

func (sc *SyncConfig) Print() {
	fmt.Printf("LastSync: %s\n", sc.LastSync)
	for _, file := range sc.Files {
		fmt.Printf("\t%s -> Last Sync: %s\n", file.Name, file.LastSync)
	}
}

func (sc *SyncConfig) SetLastSync(t time.Time) {
	sc.LastSync = t
	SyncChanged = true
}

func (sc *SyncConfig) AddFile(s SyncFile) {
	sc.Files = append(sc.Files, s)
	SyncChanged = true
}

func (sc *SyncConfig) AddNewFile(filename string) {
	t := time.Now()
	s := SyncFile{}
	s.Name = filename
	s.LastSync = t
	sc.AddFile(s)
}

func (sc *SyncConfig) IsInSync(filename string) bool {
	// Loop over files to check for existence
	for _, sc := range sc.Files {
		if filename == sc.Name {
			return true
		}
	}
	return false
}

func NewSyncConfig() *SyncConfig {
	return &SyncConfig{}
}
