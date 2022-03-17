package libdotfilesync

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jochemste/dotfile_sync/utils"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Origin   string
	Method   string
	Files    []string
	User     string
	Password string
}

func Init(file_path ...string) error {
	var origin string
	var method string
	var fp string
	if len(file_path) == 0 {
		fp = "/etc/dotfile_sync/config.toml"
	} else {
		fp = file_path[0]
	}

	if utils.FileExists(fp) {
		var choice string
		fmt.Println("File exists. Do you want to overwrite it? [y/n]")
		fmt.Scanln(&choice)
		if choice == "y" || choice == "Y" {
			//continue
		} else {
			fmt.Println("Okay, exiting")
			os.Exit(1)
		}
	}

	fmt.Println("Initialising")

	fmt.Println("Enter repository origin")
	fmt.Scanln(&origin)
	if strings.HasPrefix(origin, "git@") {
		method = "ssh"
	} else if strings.HasPrefix(origin, "http") {
		method = "http"
	} else {
		fmt.Println("Could not determine method, please enter method for the repository: [ssh, http]")
		fmt.Scanln(&method)
	}

	if method == "http" {
		log.Fatalln("Method http is not implemented yet")
	}

	fmt.Println(os.Getenv("HOME"))

	config := Config{
		Origin: origin,
		Method: method,
		Files: []string{os.Getenv("HOME") + "/" + ".bashrc",
			os.Getenv("HOME") + "/" + ".emacs"},
		User:     "",
		Password: "",
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	fmt.Println("New configuration file created here: "+fp, "\nMake sure to enter the correct paths to the files you would like to synchronise in there")

	return nil
}
