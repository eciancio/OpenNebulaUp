package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

type Config struct {
	OpenNebulaUsername string
	OpenNebulaToken    string
	OpenNebulaAPI      string
	SshKey             string
	ProjectsPath       string
	ConfigPath         string
}

var DefaultConfigPath = "/.config/OpenNebulaUp/opennebulaup.config"

func GetOpenNebulaConfig(conf *Config) goca.OneConfig {
	config := goca.NewConfig(conf.OpenNebulaUsername, conf.OpenNebulaToken, conf.OpenNebulaAPI)
	return config
}

func GetConfig(ConfigFilePath string) (*Config, error) {
	var conf Config
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		fmt.Printf("%s does not exists\n", ConfigFilePath)
		return &conf, err
	}
	if _, err := toml.DecodeFile(ConfigFilePath, &conf); err != nil {
		fmt.Println(err)
		return &conf, err
	}
	if !strings.HasSuffix(conf.ProjectsPath, "/") {
		conf.ProjectsPath = conf.ProjectsPath + "/"
	}
	return &conf, nil
}

func GetPublicKey(env Env) string {
	return env.config.SshKey
}

func GetUserHomeDir() string {
	return os.Getenv("HOME")
}

func GetAltProjectsPath(args []string) string {
	path := ""
	for num, arg := range args {
		if arg == "--ProjectsPath" {
			if num == len(args)-1 {
				fmt.Println("Please sepcify --ProjectsPath")
			} else {
				path = args[num+1]
			}
		}
	}
	return path
}

func GetConfigLocation(args []string) string {
	path := GetUserHomeDir() + DefaultConfigPath
	for num, arg := range args {
		if arg == "--config" {
			if num == len(args)-1 {
				fmt.Println("Please specify a path after --config. Using default config path")
			} else {
				path = args[num+1]
			}
		}
	}
	return path
}
