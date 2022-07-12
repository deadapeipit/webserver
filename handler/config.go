package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	TesConnectionString string `json:"tes_connection_string"`
	SecureUser          string `json:"secure_user"`
	SecurePassword      string `json:"secure_password"`
}

var Config = Configuration{}

func GetConfig() error {
	var configPath string
	wd, _ := os.Getwd()

	//check is debug true
	if strings.Contains(wd, "cmd") {
		os.Chdir("..")
		os.Chdir("..")
	}
	configPath = "config/config.json"
	jsonFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("could not read config file: %v", err.Error())
		return err
	}
	err = json.Unmarshal(jsonFile, &Config)
	if err != nil {
		fmt.Printf("could not read config file: %v", err.Error())
		return err
	}

	return nil
}
