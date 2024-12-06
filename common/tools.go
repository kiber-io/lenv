package common

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetPlatformPrefix(osName string, arch string) string {
	var prefix string
	switch osName {
	case "windows":
		prefix = "win_"
	case "linux":
		prefix = "linux_"
	case "android":
		prefix = "android_"
	default:
		return ""
	}

	switch arch {
	case "amd64":
		prefix += "x64"
	case "arm64":
		prefix += "arm64"
	default:
		return ""
	}

	return prefix
}

func ParseAssetName(assetName string) string {
	parts := strings.Split(assetName, "-")
	if len(parts) < 2 {
		return ""
	}
	nameWithoutSuffix := strings.TrimSuffix(parts[1], ".zip")
	return nameWithoutSuffix
}

func FindVersion(versions []Version, targetVersion string, targetVendor string) *Version {
	for _, v := range versions {
		if v.Version == targetVersion && v.Vendor == targetVendor {
			return &v
		}
	}
	return nil
}

func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "javaenv-jdk-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return tmpFile.Name(), nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories: %v", err)
		}
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %v", err)
		}
		defer rc.Close()
		outFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file on disk: %v", err)
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to write file to disk: %v", err)
		}
	}
	return nil
}
