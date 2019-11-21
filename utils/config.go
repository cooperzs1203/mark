/**
* @Author: Cooper
* @Date: 2019/11/19 22:04
 */

package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	GlobalConfig = &Config{}
)

type Config struct {
	Name         string `json:"name"`
	NetType      string `json:"netType"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	MaxConnCount uint32 `json:"maxConnCount"`
	MMReqCS      uint32 `json:"mMReqCS"`   // MManager request chan space
	MMRepCS      uint32 `json:"mMRepCS"`   // MManager reply chan space
	ConnRepCS    uint32 `json:"connRepCS"` // Connection reply chan space
}

func LoadConfigFile(path string) error {
	if path == "" {
		path = defaultConfigFilePath()
	}

	if path == "" {
		return errors.New("config file path was empty and get default config file path fail")
	}

	data , err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data , GlobalConfig)
	if err != nil {
		return err
	}

	log.Printf("[SUCCESS] load global config success : %+v" , GlobalConfig)

	return nil
}

func defaultConfigFilePath() string {
	projectPath , err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Join(projectPath , "conf" , "defaultConfig.json")
}