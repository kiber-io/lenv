package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type config struct {
	JdksDir           string    `yaml:"jdksDir"`
	InstalledVersions []Version `yaml:"installedJdks"`
	GlobalVersion     string    `yaml:"globalJdk"`
}

var Config config
var CurrentJdkDir string
var configFilePath string
var rootDir string

func LoadConfig() {
	Config = config{}
	rootDir = os.Getenv("JAVAENV_HOME")
	if rootDir == "" {
		fmt.Println("JAVAENV_HOME is not set")
		os.Exit(1)
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err := os.Mkdir(rootDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	Config = config{}
	configFilePath = filepath.Join(rootDir, "config.yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		data, err := yaml.Marshal(Config)
		if err != nil {
			log.Fatalf("Failed to initialize config data: %v", err)
		}
		err = os.WriteFile(configFilePath, data, 0644)
		if err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	if Config.JdksDir == "" {
		Config.JdksDir = filepath.Join(rootDir, "jdks")
	} else {
		Config.JdksDir = strings.Replace(filepath.FromSlash(Config.JdksDir), "${JAVAENV_HOME}", rootDir, 1)
	}

	CurrentJdkDir = filepath.Join(rootDir, "currentjdk")

	if Config.GlobalVersion != "" {
		Config.GlobalVersion = strings.Replace(filepath.FromSlash(Config.GlobalVersion), "${jdksDir}", Config.JdksDir, 1)
		if _, err := os.Stat(Config.GlobalVersion); os.IsNotExist(err) {
			fmt.Printf("Not found global version files at %s, unset global version\n", Config.GlobalVersion)
			Config.GlobalVersion = ""
		} else if _, err := os.Stat(CurrentJdkDir); os.IsNotExist(err) {
			fmt.Printf("Not found current jdk directory at %s, unset global version\n", CurrentJdkDir)
			Config.GlobalVersion = ""
		}
		files, err := os.ReadDir(CurrentJdkDir)
		if err != nil {
			fmt.Printf("Failed to read current jdk directory: %v, unset global version\n", err)
			Config.GlobalVersion = ""
		} else if len(files) == 0 {
			fmt.Printf("No files found in current jdk directory at %s, unset global version\n", CurrentJdkDir)
			Config.GlobalVersion = ""
		}
	}

	if _, err := os.Stat(Config.JdksDir); os.IsNotExist(err) {
		err := os.Mkdir(Config.JdksDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create jdks directory: %v", err)
		}
	}

	installedVersions := []Version{}
	for _, v := range Config.InstalledVersions {
		v.Path = filepath.FromSlash(strings.Replace(v.Path, "${jdksDir}", Config.JdksDir, 1))
		if _, err := os.Stat(v.Path); err == nil {
			installedVersions = append(installedVersions, v)
		} else {
			fmt.Printf("Version %s-%s not found at %s, removing from config", v.Version, v.Vendor, v.Path)
		}
	}
	Config.InstalledVersions = installedVersions
}

func SaveConfig() {
	configForSave := Config
	configForSave.JdksDir = filepath.ToSlash(strings.Replace(Config.JdksDir, rootDir, "${JAVAENV_HOME}", 1))
	for i, v := range configForSave.InstalledVersions {
		configForSave.InstalledVersions[i].Path = filepath.ToSlash(strings.Replace(v.Path, Config.JdksDir, "${jdksDir}", 1))
	}
	configForSave.GlobalVersion = filepath.ToSlash(strings.Replace(Config.GlobalVersion, Config.JdksDir, "${jdksDir}", 1))
	data, err := yaml.Marshal(configForSave)
	if err != nil {
		log.Fatalf("Failed to marshal config data: %v", err)
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}

func AddVersion(version Version) {
	Config.InstalledVersions = append(Config.InstalledVersions, version)
}

func RemoveVersion(version Version) bool {
	for i, v := range Config.InstalledVersions {
		if v.Version == version.Version && v.Vendor == version.Vendor {
			Config.InstalledVersions = append(Config.InstalledVersions[:i], Config.InstalledVersions[i+1:]...)
			return true
		}
	}
	return false
}
