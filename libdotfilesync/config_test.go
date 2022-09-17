package libdotfilesync

import (
	"testing"
	"time"

	"github.com/jochemste/dotfile_sync/utils"
)

// TestConfigStoredAndRetrievedGW
// Brief: Test if the configuration is stored and retrieved correctly
// Good weather: Create config and store it in file. Get other config from same
//               file and ensure they are equal
// Pass if:
//        - Configuration 1 is stored in file
//        - Configuration 2 is retrieved from file
//        - Configurations are equal
//        - No errors occur
func TestConfigStoredAndRetrievedGW(t *testing.T) {
	config := NewConfig()
	configfile := "/tmp/.temptest"

	config.UserSettings.Origin = "test.com"
	config.UserSettings.Files = []string{"test1", "test2", "test3"}
	config.UserSettings.AccessToken = "123456"
	config.UserSettings.Mode = ""

	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 124
	config.DoNotChange.Store = "./etete"
	config.DoNotChange.File = "slslslsl"

	config.ToFile(configfile)

	config2 := NewConfig()
	config2.FromFile(configfile)

	// Cleanup already
	utils.RemoveFile(configfile)

	isSame, err := config.Equals(config2)

	if !isSame && err == nil {
		t.Errorf("Stored configuration is not equal to retrieved" +
			"configuration, but no error message was provided")
	}

	if err != nil {
		t.Errorf("Stored configuration is not equal to retrieved configuration: %s",
			err)
	}
}

// TestConfigIsDifferentTimeGW
// Brief: Test if the configuration is stored and retrieved correctly
// Good weather: Create config and store it in file. Get other config from same
//               file, change the time and ensure they are not equal
// Pass if:
//        - Configuration 1 is stored in file
//        - Configuration 2 is retrieved from file
//        - Configurations are not equal
//        - No errors occur
func TestConfigIsDifferentTimeGW(t *testing.T) {
	config := NewConfig()
	configfile := "/tmp/.temptest"

	config.UserSettings.Origin = "test.com"
	config.UserSettings.Files = []string{"test1", "test2", "test3"}
	config.UserSettings.AccessToken = "123456"
	config.UserSettings.Mode = ""

	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 124
	config.DoNotChange.Store = "./etete"
	config.DoNotChange.File = "slslslsl"

	config.ToFile(configfile)

	config2 := NewConfig()
	config2.FromFile(configfile)
	config2.DoNotChange.LastCheck = config2.DoNotChange.LastCheck.Add(time.Hour * 1)
	// Cleanup already
	utils.RemoveFile(configfile)

	isSame, _ := config.Equals(config2)

	if isSame {
		t.Errorf("Configurations should not be equal, since config2 was adjusted\n%s\n%s",
			config.String(), config2.String())
	}
}

// TestConfigIsDifferentFilesGW
// Brief: Test if the configuration is stored and retrieved correctly
// Good weather: Create config and store it in file. Get other config from same
//               file, add a file and ensure they are not equal
// Pass if:
//        - Configuration 1 is stored in file
//        - Configuration 2 is retrieved from file
//        - Configurations are not equal
//        - No errors occur
func TestConfigIsDifferentFilesGW(t *testing.T) {
	config := NewConfig()
	configfile := "/tmp/.temptest"

	config.UserSettings.Origin = "test.com"
	config.UserSettings.Files = []string{"test1", "test2", "test3"}
	config.UserSettings.AccessToken = "123456"
	config.UserSettings.Mode = ""

	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 124
	config.DoNotChange.Store = "./etete"
	config.DoNotChange.File = "slslslsl"

	config.ToFile(configfile)

	config2 := NewConfig()
	config2.FromFile(configfile)
	config2.AddFile("file123.txt")

	// Cleanup already
	utils.RemoveFile(configfile)

	isSame, _ := config.Equals(config2)

	if isSame {
		t.Errorf("Configurations should not be equal, since config2 was adjusted\n%s\n%s",
			config.String(), config2.String())
	}
}

// TestConfigIsDifferentFilesGW2
// Brief: Test if the configuration is stored and retrieved correctly
// Good weather: Create config and store it in file. Get other config from same
//               file, remove file and ensure they are not equal
// Pass if:
//        - Configuration 1 is stored in file
//        - Configuration 2 is retrieved from file
//        - Configurations are not equal
//        - No errors occur
func TestConfigIsDifferentFilesGW2(t *testing.T) {
	config := NewConfig()
	configfile := "/tmp/.temptest"

	config.UserSettings.Origin = "test.com"
	config.UserSettings.Files = []string{"test1", "test2", "test3"}
	config.UserSettings.AccessToken = "123456"
	config.UserSettings.Mode = ""

	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 124
	config.DoNotChange.Store = "./etete"
	config.DoNotChange.File = "slslslsl"

	config.ToFile(configfile)

	config2 := NewConfig()
	config2.FromFile(configfile)
	config2.RemoveFile("test2")

	// Cleanup already
	utils.RemoveFile(configfile)

	isSame, _ := config.Equals(config2)

	if isSame {
		t.Errorf("Configurations should not be equal, since config2 was adjusted\n%s\n%s",
			config.String(), config2.String())
	}
}
