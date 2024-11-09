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
		prefix = "win"
	case "darwin":
		prefix = "mac"
	case "linux":
		prefix = "linux"
	default:
		return ""
	}

	switch arch {
	case "amd64":
		prefix += "64"
	default:
		return ""
	}

	return prefix
}

func ParseAssetName(assetName string) string {
	// Разделяем имя по символу "-"
	parts := strings.Split(assetName, "-")

	// Проверяем, что у нас достаточно частей после разделения
	if len(parts) < 2 {
		return ""
	}

	// Удаляем суффикс ".zip" из второй части
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
	// Выполняем HTTP GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем успешность ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Создаем временный файл в системной временной директории
	tmpFile, err := os.CreateTemp("", "javaenv-jdk-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// Копируем содержимое ответа в временный файл
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return tmpFile.Name(), nil
}

func Unzip(src, dest string) error {
	// Открываем ZIP-файл
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Перебираем файлы в архиве
	for _, f := range r.File {
		// Определяем путь для распаковки файла
		fpath := filepath.Join(dest, f.Name)

		// Создаем директории для вложенных файлов
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Создаем необходимые директории для файла
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories: %v", err)
		}

		// Открываем файл из архива
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %v", err)
		}
		defer rc.Close()

		// Создаем файл на диске
		outFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file on disk: %v", err)
		}
		defer outFile.Close()

		// Копируем содержимое из ZIP-файла в новый файл
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to write file to disk: %v", err)
		}
	}
	return nil
}
