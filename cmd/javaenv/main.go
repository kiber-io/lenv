package main

import (
	"fmt"
	"kiber-io/javaenv/common"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	common.LoadConfig()

	var showAll bool
	var rootCmd = &cobra.Command{
		Use:               "javaenv",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	var installCmd = &cobra.Command{
		Use:   "install [version]",
		Short: "Install specific Java version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			install(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	var uninstallCmd = &cobra.Command{
		Use:   "uninstall [version]",
		Short: "Uninstall specific Java version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uninstall(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List installed or available Java versions",
		Run: func(cmd *cobra.Command, args []string) {
			if showAll {
				listAvailable()
			} else {
				listInstalled()
			}
		},
	}
	var globalCmd = &cobra.Command{
		Use:   "global",
		Short: "Set global Java version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			setGlobal(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of javaenv",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("javaenv v" + version)
		},
	}

	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all available versions")
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(globalCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()

	common.SaveConfig()
}

func install(version string) {
	parts := strings.Split(version, "-")
	installed := common.FindVersion(common.Config.InstalledVersions, parts[0], parts[1])
	if installed != nil {
		fmt.Printf("Java version %s is already installed\n", version)
		return
	}
	fmt.Println("Downloading...")
	platformPrefix := common.GetPlatformPrefix(runtime.GOOS, runtime.GOARCH)
	url := fmt.Sprintf("https://github.com/kiber-io/jdks/releases/download/%s/%s-%s.zip", parts[0], platformPrefix, parts[1])
	filePath, err := common.DownloadFile(url)
	if err != nil {
		fmt.Println("Failed to download file: ", err)
		return
	}
	jdkDir := filepath.Join(common.Config.JdksDir, version)
	fmt.Println("Extracting...")
	common.Unzip(filePath, jdkDir)
	os.Remove(filePath)
	fmt.Printf("Java version %s installed\n", version)
	common.AddVersion(common.Version{
		Version: parts[0],
		Path:    jdkDir,
		Vendor:  parts[1],
	})
}

func listInstalled() {
	if len(common.Config.InstalledVersions) == 0 {
		fmt.Println("No versions installed")
		return
	}
	fmt.Println("Installed Versions:")
	for _, version := range common.Config.InstalledVersions {
		prefix := "    "
		if version.Path == common.Config.GlobalVersion {
			prefix = " -> "
		}
		fmt.Printf("%s%s-%s", prefix, version.Version, version.Vendor)
		fmt.Println()
	}
}

func listAvailable() {
	versions, err := common.FetchVersions(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		log.Fatalf("Error fetching versions: %v", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions available for your platform and architecture")
		return
	}

	fmt.Println("Available Versions:")
	for _, version := range versions {
		installed := common.FindVersion(common.Config.InstalledVersions, version.Version, version.Vendor)
		prefix := "    "
		if installed != nil {
			prefix = "  * "
			if installed.Path == common.Config.GlobalVersion {
				prefix = " -> "
			}
		}
		fmt.Printf("%s%s-%s", prefix, version.Version, version.Vendor)
		fmt.Println()
	}
}

func setGlobal(version string) {
	parts := strings.Split(version, "-")
	installed := common.FindVersion(common.Config.InstalledVersions, parts[0], parts[1])
	if installed == nil {
		log.Fatalf("Java version %s is not installed", version)
	}
	switch runtime.GOOS {
	case "windows":
		os.Remove(common.CurrentJdkDir)
		cmd := exec.Command("cmd", "/c", "mklink", "/J", common.CurrentJdkDir, installed.Path)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to set global version: %v", err)
		}
		setJavaHome()
		fmt.Printf("Java version %s set as global\n", version)
	default:
		log.Fatalf("Unknown operating system: %s", runtime.GOOS)
	}

	common.Config.GlobalVersion = installed.Path
}

func setJavaHome() {
	needSet := false
	value, exists := os.LookupEnv("JAVA_HOME")
	if !exists {
		needSet = true
	} else if value != common.CurrentJdkDir {
		needSet = true
	}
	if needSet {
		err := os.Setenv("JAVA_HOME", common.CurrentJdkDir)
		if err != nil {
			log.Fatalf("Failed to set JAVA_HOME: %v", err)
		}
		cmd := exec.Command("cmd", "/c", "setx", "JAVA_HOME", common.CurrentJdkDir)
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to set JAVA_HOME: %v", err)
		}
	}
}

func uninstall(version string) {
	parts := strings.Split(version, "-")
	installed := common.FindVersion(common.Config.InstalledVersions, parts[0], parts[1])
	if installed == nil {
		fmt.Printf("Java version %s is not installed\n", version)
		return
	}
	if installed.Path == common.Config.GlobalVersion {
		fmt.Printf("Java version %s is set as global, unset it first\n", version)
		return
	}
	err := os.RemoveAll(installed.Path)
	if err != nil {
		fmt.Printf("Failed to uninstall Java version %s: %v\n", version, err)
		return
	}
	fmt.Printf("Java version %s uninstalled", version)
	removed := common.RemoveVersion(*installed)
	if !removed {
		log.Fatalf("Failed to remove version %s from config", version)
	}
}
