package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	InstalledVersions []Version
	VersionsDir       string
	CurrentVersionDir string
	GlobalVersion     Version
}

var Config config
var rootDir string
var languageDir string

func GetRoot() string {
	dir := os.Getenv("LENV_HOME")
	if dir == "" {
		fmt.Println("LENV_HOME is not set")
		os.Exit(1)
	}
	return dir
}

func LoadConfig(language string) {
	rootDir = GetRoot()
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		log.Fatal("LENV_HOME directory not found")
	}
	if language != "java" && language != "python" {
		log.Fatalf("Unknown language: %s", language)
	}
	languageDir = filepath.Join(rootDir, strings.ToLower(language))
	if _, err := os.Stat(languageDir); os.IsNotExist(err) {
		err := os.Mkdir(languageDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}
	versionsDir := filepath.Join(languageDir, "versions")
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		err := os.Mkdir(versionsDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create versions directory: %v", err)
		}
	}
	Config.VersionsDir = versionsDir
	folders, err := os.ReadDir(versionsDir)
	if err != nil {
		log.Fatalf("Failed to read language directory: %v", err)
	}
	for _, folder := range folders {
		if folder.IsDir() {
			parts := strings.Split(folder.Name(), "-")
			version := Version{
				Version: parts[0],
				Vendor:  parts[1],
				Path:    filepath.Join(versionsDir, folder.Name()),
			}
			Config.InstalledVersions = append(Config.InstalledVersions, version)
		} else {
			fmt.Printf("Unexpected file found in versions directory: %s", folder.Name())
		}
	}
	Config.CurrentVersionDir = filepath.Join(languageDir, "current")
	globalVersionFile := filepath.Join(languageDir, "global")
	if _, err := os.Stat(globalVersionFile); os.IsNotExist(err) {
		err := os.WriteFile(globalVersionFile, []byte(""), 0644)
		if err != nil {
			log.Fatalf("Failed to create global version file: %v", err)
		}
	}
	data, err := os.ReadFile(globalVersionFile)
	if err != nil {
		log.Fatalf("Failed to read global version file: %v", err)
	}
	version := string(data)
	for _, v := range Config.InstalledVersions {
		if v.Name() == version {
			Config.GlobalVersion = v
			break
		}
	}
	if Config.GlobalVersion == (Version{}) {
		err := os.WriteFile(globalVersionFile, []byte(""), 0644)
		if err != nil {
			log.Fatalf("Failed to write global version file: %v", err)
		}
	}
}

func SetGlobalVersion(version Version) {
	globalVerionFile := filepath.Join(languageDir, "global")
	err := os.WriteFile(globalVerionFile, []byte(version.Name()), 0644)
	if err != nil {
		log.Fatalf("Failed to write global version file: %v", err)
	}
	Config.GlobalVersion = version
}
