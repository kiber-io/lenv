package common

import (
	"fmt"
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

func FindVersion(versions []Version, targetVersion string, targetVendor string) *Version {
	for _, v := range versions {
		if v.Version == targetVersion && v.Vendor == targetVendor {
			return &v
		}
	}
	return nil
}

func ParseAssetName(assetName string) string {
	parts := strings.Split(assetName, "-")
	if len(parts) < 2 {
		return ""
	}
	nameWithoutSuffix := strings.TrimSuffix(parts[1], ".zip")
	return nameWithoutSuffix
}
