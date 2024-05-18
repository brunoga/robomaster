package support

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// FindFile tries to locate the given file in the following order:
//
// 1. Only the file name in the install directory.
// 2. Only the file name In the current directory.
// 3. The full path in the current directory.
func FindFile(expectedPath string) string {
	fileName := filepath.Base(expectedPath)

	if path, err := findInInstallDir(fileName); err == nil {
		return path
	}

	if path, err := findInCurrentDir(fileName); err == nil {
		return path
	}

	if path, err := findInCurrentDir(expectedPath); err == nil {
		return path
	}

	// Return filename in case the OS supports searching for libraries in
	// specific paths.
	return fileName
}

func findInInstallDir(fileName string) (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	installDirAndFilename := filepath.Join(homeDir, ".unitybridge", fileName)

	if _, err := os.Stat(installDirAndFilename); err != nil {
		if runtime.GOOS == "windows" {
			// HACK! We might be under Wine so replace C:\users with Z:\home
			// as we we know we are under Linux.
			//
			// TODO(bga): There must be a better way to do this.
			installDirAndFilename = filepath.Join(
				"Z:\\home", installDirAndFilename[8:])
			if _, err := os.Stat(installDirAndFilename); err != nil {
				return "", err
			}

			return installDirAndFilename, nil
		} else {
			return "", err
		}
	}

	return installDirAndFilename, nil
}

func findInCurrentDir(fileName string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	currentDir = filepath.Join(currentDir, fileName)
	if _, err := os.Stat(currentDir); err == nil {
		return currentDir, nil
	}

	return "", fmt.Errorf(
		"could not find file %s in current dir %s",
		fileName, currentDir)
}
