package java

import (
	"encoding/json"
	"fmt"
	"io"
	"kiber-io/lenv/common"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var showAll bool

func Init(javaCmd *cobra.Command) {
	var installCmd = &cobra.Command{
		Use:     "install",
		Short:   "Install specific Java version",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			install(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	var uninstallCmd = &cobra.Command{
		Use:     "uninstall [version]",
		Short:   "Uninstall specific Java version",
		Aliases: []string{"u"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uninstall(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed or available Java versions",
		Run: func(cmd *cobra.Command, args []string) {
			if showAll {
				listAvailable()
			} else {
				listInstalled()
			}
		},
	}
	var globalCmd = &cobra.Command{
		Use:     "global",
		Aliases: []string{"g"},
		Short:   "Set global Java version",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			setGlobal(args[0])
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all available versions")

	javaCmd.AddCommand(installCmd)
	javaCmd.AddCommand(uninstallCmd)
	javaCmd.AddCommand(listCmd)
	javaCmd.AddCommand(globalCmd)

	javaCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		common.LoadConfig("java")
	}
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
	url := fmt.Sprintf("https://github.com/kiber-io/lenv-java-versions/releases/download/%s/%s-%s.zip", parts[0], platformPrefix, parts[1])
	filePath, err := common.DownloadFile(url)
	if err != nil {
		fmt.Println("Failed to download file: ", err)
		return
	}
	jdkDir := filepath.Join(common.Config.VersionsDir, version)
	fmt.Println("Extracting...")
	common.Unzip(filePath, jdkDir)
	os.Remove(filePath)
	cmd := exec.Command("chmod", "-R", "755", jdkDir)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to set permissions: %v", err)
	}
	fmt.Printf("Java version %s installed\n", version)
}

func listInstalled() {
	if len(common.Config.InstalledVersions) == 0 {
		fmt.Println("No versions installed")
		return
	}
	fmt.Println("Installed Versions:")
	for _, version := range common.Config.InstalledVersions {
		prefix := "    "
		if version.Name() == common.Config.GlobalVersion.Name() {
			prefix = " -> "
		}
		fmt.Printf("%s%s-%s", prefix, version.Version, version.Vendor)
		fmt.Println()
	}
}

func listAvailable() {
	versions, err := FetchVersions(runtime.GOOS, runtime.GOARCH)
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
			if installed.Name() == common.Config.GlobalVersion.Name() {
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
	common.SetGlobalVersion(*installed)
	switch runtime.GOOS {
	case "windows":
		setGlobalWindows(*installed)
	case "linux", "android":
		setGlobalLinux(*installed)
	default:
		log.Fatalf("Unknown operating system: %s", runtime.GOOS)
	}
	fmt.Printf("Java version %s set as global\n", version)
}

func setGlobalWindows(version common.Version) {
	os.Remove(common.Config.CurrentVersionDir)
	cmd := exec.Command("cmd", "/c", "mklink", "/J", common.Config.CurrentVersionDir, version.Path)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to set global version: %v", err)
	}
}

func setGlobalLinux(version common.Version) {
	os.Remove(common.Config.CurrentVersionDir)
	err := os.Symlink(version.Path, common.Config.CurrentVersionDir)
	if err != nil {
		log.Fatalf("Failed to set global version: %v", err)
	}
}

func uninstall(version string) {
	parts := strings.Split(version, "-")
	installed := common.FindVersion(common.Config.InstalledVersions, parts[0], parts[1])
	if installed == nil {
		fmt.Printf("Java version %s is not installed\n", version)
		return
	}
	if installed.Name() == common.Config.GlobalVersion.Name() {
		fmt.Printf("Java version %s is set as global, are you sure you want to uninstall it? [y/N]: ", version)
		var response string
		fmt.Scanln(&response)
		if strings.TrimSpace(strings.ToLower(response)) != "y" {
			return
		}
	}
	err := os.RemoveAll(installed.Path)
	if err != nil {
		fmt.Printf("Failed to uninstall Java version %s: %v\n", version, err)
		return
	}
	fmt.Printf("Java version %s uninstalled\n", version)
}

func FetchVersions(platform string, arch string) ([]common.Version, error) {
	platformPrefix := common.GetPlatformPrefix(platform, arch)
	if platformPrefix == "" {
		return nil, fmt.Errorf("unknown operating system and architecture: %s/%s", platform, arch)
	}

	url := "https://api.github.com/repos/kiber-io/lenv-java-versions/releases"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var versions []common.ServerVersion
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	filteredVersions := []common.Version{}
	for _, serverVersion := range versions {
		filteredAssets := []common.Asset{}
		for _, asset := range serverVersion.Assets {
			if strings.HasPrefix(asset.Name, platformPrefix) {
				filteredAssets = append(filteredAssets, asset)
			}
		}

		if len(filteredAssets) > 0 {
			for _, asset := range filteredAssets {
				version := common.Version{
					Version: serverVersion.TagName,
					Path:    "",
					Vendor:  common.ParseAssetName(asset.Name),
				}
				filteredVersions = append(filteredVersions, version)
			}
		}
	}

	return filteredVersions, nil
}
