package libdotfilesync

import (
	"testing"
	"time"

	"github.com/jochemste/dotfile_sync/utils"
)

func TestConfigStoredAndRetrieved(t *testing.T) {
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
