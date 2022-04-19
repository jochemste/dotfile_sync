package libdotfilesync

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jochemste/dotfile_sync/utils"

	"github.com/pelletier/go-toml/v2"
)

// Configuration that is managed by the user
type UserConfig struct {
	Origin      string
	Files       []string
	AccessToken string
}

// Information that is used and adjusted automatically
type DontChangeConfig struct {
	LastCheck time.Time
	NrSync    int
	Store     string
}

// Configuration struct
type Config struct {
	UserSettings UserConfig
	DoNotChange  DontChangeConfig
}

// Write the Config object to a file
func (config *Config) ToFile(file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return errors.New("Could not create file " + file)
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

	return nil
}

func (config *Config) Print() {
	fmt.Printf("Configuration:\n\t%s:\t%s\n", "Origin", config.UserSettings.Origin)
	fmt.Printf("\tFiles:\n")
	for _, file := range config.UserSettings.Files {
		fmt.Printf("\t\t%s\n", file)
	}
}

// Get a new Config object
func NewConfig() *Config {
	return &Config{}
}
