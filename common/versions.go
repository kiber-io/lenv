package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Version struct {
	Version string `json:"version"`
	Path    string `json:"path"`
	Vendor  string `json:"vendor"`
}

func (v Version) Name() string {
	return fmt.Sprintf("%s-%s", v.Version, v.Vendor)
}

type ServerVersion struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
}

func FetchVersions(platform string, arch string) ([]Version, error) {
	platformPrefix := GetPlatformPrefix(platform, arch)
	if platformPrefix == "" {
		return nil, fmt.Errorf("unknown operating system and architecture: %s/%s", platform, arch)
	}

	url := "https://api.github.com/repos/kiber-io/jdks/releases"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var versions []ServerVersion
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	filteredVersions := []Version{}
	for _, serverVersion := range versions {
		filteredAssets := []Asset{}
		for _, asset := range serverVersion.Assets {
			if strings.HasPrefix(asset.Name, platformPrefix) {
				filteredAssets = append(filteredAssets, asset)
			}
		}

		if len(filteredAssets) > 0 {
			for _, asset := range filteredAssets {
				version := Version{
					Version: serverVersion.TagName,
					Path:    "",
					Vendor:  ParseAssetName(asset.Name),
				}
				filteredVersions = append(filteredVersions, version)
			}
		}
	}

	return filteredVersions, nil
}
