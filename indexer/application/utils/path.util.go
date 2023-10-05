package utils

import (
	"io/fs"
	"os"
	"path/filepath"

	customlogger "pkg/custom-logger"
)

var loggerPathUtil = customlogger.NewLogger()

func PathExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func GetFilesPaths(dirPath string, extension string) ([]string, error) {
	filesPaths := []string{}

	err := filepath.WalkDir(
		dirPath,
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				loggerPathUtil.Println("error getting files paths: ", err)

				return err
			}

			if extension == "" {
				if !info.IsDir() {
					filesPaths = append(filesPaths, path)
				}
			} else {
				if !info.IsDir() && filepath.Ext(info.Name()) == extension {
					filesPaths = append(filesPaths, path)
				}
			}

			return nil
		},
	)

	if err != nil {
		loggerPathUtil.Println("error getting files paths: ", err)

		return nil, err
	}

	return filesPaths, nil
}
