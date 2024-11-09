package common

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	JdksDir           string    `json:"jdksDir"`
	InstalledVersions []Version `json:"installedJdks"`
	GlobalVersion     string    `json:"globalJdk"`
}

var Config config
var CurrentJdkDir string
var configFilePath string

func LoadConfig() {
	Config = config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	rootDir := filepath.Join(homeDir, ".javaenv")
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err := os.Mkdir(rootDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	Config = config{}
	configFilePath = filepath.Join(rootDir, "config.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		data, err := json.MarshalIndent(Config, "", "  ")
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
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}
	if Config.JdksDir == "" {
		Config.JdksDir = filepath.Join(rootDir, "jdks")
	}
	CurrentJdkDir = filepath.Join(rootDir, "currentjdk")
	if _, err := os.Stat(Config.JdksDir); os.IsNotExist(err) {
		err := os.Mkdir(Config.JdksDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create jdks directory: %v", err)
		}
	}
}

func SaveConfig() {
	data, err := json.MarshalIndent(Config, "", "  ")
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
			// Удаляем элемент, сдвигая остальные элементы
			Config.InstalledVersions = append(Config.InstalledVersions[:i], Config.InstalledVersions[i+1:]...)
			return true
		}
	}
	return false
}
