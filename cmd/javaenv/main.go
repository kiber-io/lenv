package main

import (
	"bufio"
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

var version = "0.1.1"

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
	case "linux":
		os.Remove(common.CurrentJdkDir)
		err := os.Symlink(installed.Path, common.CurrentJdkDir)
		if err != nil {
			log.Fatalf("Failed to set global version: %v", err)
		}
		cmd := exec.Command("chmod", "-R", "755", common.CurrentJdkDir)
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to set global version: %v", err)
		}
		setJavaHome()
	default:
		log.Fatalf("Unknown operating system: %s", runtime.GOOS)
	}

	common.Config.GlobalVersion = installed.Path
	fmt.Printf("Java version %s set as global. Restart the terminal if you need to use the updated JAVA_HOME variable\n", version)
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
		switch runtime.GOOS {
		case "windows":
			setWindowsJavaHome()
		case "linux":
			setLinuxJavaHome()
		default:
			log.Fatalf("Unknown operating system: %s", runtime.GOOS)
		}
	}
}

func setWindowsJavaHome() {
	cmd := exec.Command("cmd", "/c", "setx", "JAVA_HOME", common.CurrentJdkDir)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to set JAVA_HOME: %v", err)
	}
}

func setLinuxJavaHome() {
	profilePath := filepath.Join(os.Getenv("HOME"), ".profile")
	addJavaHomeToFile(profilePath)
	bashrcPath := filepath.Join(os.Getenv("HOME"), ".bashrc")
	addJavaHomeToFile(bashrcPath)
}

func addJavaHomeToFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			if err != nil {
				log.Fatalf("Failed to create %s: %v", path, err)
			}
		} else {
			log.Fatalf("Failed to open %s: %v", path, err)
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	javaHomeFound := false
	javaHomeLine := fmt.Sprintf("export JAVA_HOME=%s", common.CurrentJdkDir)

	for scanner.Scan() {
		line := scanner.Text()
		if line == javaHomeLine {
			javaHomeFound = true
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read %s: %v", path, err)
	}

	if !javaHomeFound {
		lines = append(lines, javaHomeLine)
		outputFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Failed to open %s for writing: %v", path, err)
		}
		defer outputFile.Close()

		for _, line := range lines {
			_, err = fmt.Fprintln(outputFile, line)
			if err != nil {
				log.Fatalf("Failed to write to %s: %v", path, err)
			}
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
	fmt.Printf("Java version %s uninstalled\n", version)
	removed := common.RemoveVersion(*installed)
	if !removed {
		log.Fatalf("Failed to remove version %s from config", version)
	}
}
