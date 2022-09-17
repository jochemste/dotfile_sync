package libdotfilesync

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jochemste/dotfile_sync/utils"

	"github.com/pelletier/go-toml/v2"
)

type Configuration interface {
	ToFile(file ...string) error
	FromFile(file string) error
	IsWritable() bool
	Equals(other *Config) (bool, error)
	IsSameLastCheck(other *Config) (bool, error)
	LastCheckIsAfter(other *Config) (bool, error)
	LastCheckIsBefore(other *Config) (bool, error)
	UpdateLastCheck() error
	AddFile(file string) error
	RemoveFile(file string) error
	String() string
	Print()
}

// Configuration that is managed by the user
type UserConfig struct {
	Origin      string
	Files       []string
	AccessToken string
	Mode        string
}

// Information that is used and adjusted automatically
type DontChangeConfig struct {
	LastCheck time.Time
	NrSync    int
	Store     string
	File      string
}

// Configuration struct
type Config struct {
	UserSettings UserConfig
	DoNotChange  DontChangeConfig
}

// Write the Config object to a file
func (config *Config) ToFile(file ...string) error {
	var fp string
	if len(file) == 0 {
		fp = config.DoNotChange.File
	} else {
		fp = file[0]
		config.DoNotChange.File = file[0]
	}

	fd, err := os.Create(fp)
	if err != nil {
		return errors.New("Could not create file " + fp)
	}

	if err := toml.NewEncoder(fd).Encode(config); err != nil {
		return err
	}

	if err := fd.Close(); err != nil {
		return err
	}

	return nil
}

// Set a Config from a toml file, if the file exists
func (config *Config) FromFile(file string) error {
	if !utils.FileExists(file) {
		return errors.New("File " + file + " does not exist")
	}

	content, err := utils.GetFileContent(file)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal([]byte(content), &config); err != nil {
		return err
	}

	config.DoNotChange.File = file

	writable := utils.IsWritable(file)
	if writable != true {
		return errors.New("File " + file + " is not writable for the current user")
	}
	return nil
}

func (config *Config) IsWritable() bool {
	return utils.IsWritable(config.DoNotChange.File)
}

func (config *Config) Equals(other *Config) (bool, error) {
	isSame := true
	var err error
	err = nil

	if config == other {
		isSame = true
	}

	// Check UserSettings
	if isSame &&
		(config.UserSettings.Origin != other.UserSettings.Origin ||
			config.UserSettings.AccessToken != other.UserSettings.AccessToken ||
			config.UserSettings.Mode != other.UserSettings.Mode) {
		isSame = false
		err = errors.New("UserSettings members are not equal" +
			config.String() + "\n" + other.String())
	}

	// Check files is UserSettings to be the same, including the order
	if isSame &&
		(len(config.UserSettings.Files) != len(other.UserSettings.Files)) {
		isSame = false
		err = errors.New(fmt.Sprintf("File lists are not of equal length: %d, %d",
			len(config.UserSettings.Files), len(other.UserSettings.Files)))
	}
	if isSame {
		for i, _ := range config.UserSettings.Files {
			if config.UserSettings.Files[i] != other.UserSettings.Files[i] {
				isSame = false
				err = errors.New("UserSettings files are not equal: " +
					config.UserSettings.Files[i] + " " + other.UserSettings.Files[i])
			}
		}
	}

	// Check DoNotChange members
	if isSame &&
		(config.DoNotChange.NrSync != other.DoNotChange.NrSync ||
			config.DoNotChange.Store != other.DoNotChange.Store ||
			config.DoNotChange.File != other.DoNotChange.File) {
		isSame = false
		err = errors.New("DoNotChange members are not equal\n" +
			config.String() + "\n" + other.String())
	}

	if isSame {
		isSame, err = config.IsSameLastCheck(other)
		if err != nil {
			err = errors.New("LastCheck is different: " + err.Error())
		}
	}

	return isSame, err
}

func (config *Config) IsSameLastCheck(other *Config) (bool, error) {
	var err error
	err = nil
	isSame := true

	if config.DoNotChange.LastCheck.Equal(other.DoNotChange.LastCheck) == false {
		isSame = false
		err = errors.New("DoNotChange.LastCheck is not equal:\n" +
			fmt.Sprintf("%s\n%s\n", config.DoNotChange.LastCheck,
				other.DoNotChange.LastCheck))
	}

	return isSame, err
}

func (config *Config) LastCheckIsAfter(other *Config) (bool, error) {
	var err error
	err = nil
	isAfter := true

	if config.DoNotChange.LastCheck.After(other.DoNotChange.LastCheck) == false {
		isAfter = false
		err = errors.New("DoNotChange.LastCheck is not after:\n" +
			fmt.Sprintf("%s\n%s\n", config.DoNotChange.LastCheck,
				other.DoNotChange.LastCheck))
	}

	return isAfter, err
}

func (config *Config) LastCheckIsBefore(other *Config) (bool, error) {
	var err error
	err = nil
	isBefore := true

	if config.DoNotChange.LastCheck.Before(other.DoNotChange.LastCheck) == false {
		isBefore = false
		err = errors.New("DoNotChange.LastCheck is not before:\n" +
			fmt.Sprintf("%s\n%s\n", config.DoNotChange.LastCheck,
				other.DoNotChange.LastCheck))
	}

	return isBefore, err
}

func (config *Config) UpdateLastCheck() error {
	var err error
	err = nil
	config.DoNotChange.LastCheck = time.Now()

	return err
}

func (config *Config) FileInConfig(file string) bool {
	exists := false
	for i, _ := range config.UserSettings.Files {
		if file == config.UserSettings.Files[i] {
			exists = true
		}
	}
	return exists
}

func (config *Config) AddFile(file string) error {
	var err error
	err = nil
	exists := config.FileInConfig(file)

	if !exists {
		config.UserSettings.Files = append(config.UserSettings.Files, file)
	} else {
		err = errors.New("File already exists in config")
	}

	return err
}

func (config *Config) RemoveFile(file string) error {
	var err error
	err = nil
	exists := config.FileInConfig(file)

	if exists {
		files := config.UserSettings.Files
		config.UserSettings.Files = []string{}
		for i, _ := range files {
			if files[i] != file {
				config.AddFile(files[i])
			}
		}
	} else {
		err = errors.New("File already exists in config")
	}

	return err
}

func (config *Config) String() string {
	var str string
	str = "Config\n"
	str += "UserSettings:\n"
	str += fmt.Sprintf("Mode: %s\n", config.UserSettings.Mode)
	str += fmt.Sprintf("Origin: %s\n", config.UserSettings.Origin)
	str += fmt.Sprintf("AccessToken: %s\n", config.UserSettings.AccessToken)
	for i, _ := range config.UserSettings.Files {
		str += fmt.Sprintf("File %d: %s\n", i, config.UserSettings.Files[i])
	}
	str += "DoNotChange:\n"
	str += fmt.Sprintf("LastCheck: %s\n", config.DoNotChange.LastCheck)
	str += fmt.Sprintf("NrSync: %d\n", config.DoNotChange.NrSync)
	str += fmt.Sprintf("Store: %s\n", config.DoNotChange.Store)
	str += fmt.Sprintf("File: %s\n", config.DoNotChange.File)

	return str
}

func (config *Config) Print() {
	fmt.Printf("Configuration:\n")
	fmt.Printf("\tLocation: %s\n", config.DoNotChange.File)
	fmt.Printf("\tOrigin: %s\n", config.UserSettings.Origin)
	fmt.Printf("\tMode: %s\n", config.UserSettings.Mode)
	fmt.Printf("\tFiles:\n")
	for _, file := range config.UserSettings.Files {
		fmt.Printf("\t\t%s\n", file)
	}
	fmt.Printf("\tNr Synchronisations: %d\n", config.DoNotChange.NrSync)
	fmt.Printf("\tLastCheck: %s\n", config.DoNotChange.LastCheck)
}

// Get a new Config object
func NewConfig() *Config {
	return &Config{}
}
