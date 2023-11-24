package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	var userHomeDir string
	var err error

	if userHomeDir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}

	unityBridgeDir := filepath.Join(userHomeDir, ".unitybridge")

	if err = os.MkdirAll(unityBridgeDir, 0755); err != nil {
		panic(err)
	}

	if err = installUnityBridgeLibrary(unityBridgeDir); err != nil {
		panic(err)
	}

	if err = buildAndInstallDLLHost(unityBridgeDir); err != nil {
		panic(err)
	}
}

func buildAndInstallDLLHost(unityBridgeDir string) error {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil
	}

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	os.Setenv("GOOS", "windows")
	defer os.Unsetenv("GOOS")

	os.Setenv("CGO_ENABLED", "1")
	defer os.Unsetenv("CGO_ENABLED")

	os.Setenv("CC", "x86_64-w64-mingw32-gcc")
	defer os.Unsetenv("CC")

	cmd := exec.Command("go", "build", "-o", filepath.Join(unityBridgeDir,
		"dllhost.exe"), filepath.Join(currDir, "wrapper", "internal", "implementations",
		"wine", "dllhost"))

	return cmd.Run()
}

func installUnityBridgeLibrary(unityBridgeDir string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	goos := runtime.GOOS
	if goos == "linux" {
		// Linux will use the Windows DLL through Wine.
		goos = "windows"
	}

	srcDir := filepath.Join(currDir, "wrapper", "lib", goos,
		runtime.GOARCH)
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS,
			runtime.GOARCH)
	}

	err = copyDir(srcDir, unityBridgeDir)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(srcDir, destDir string) error {
	// Open source directory.
	src, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination directory.
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// Read source directory contents.
	fileInfos, err := src.Readdir(-1)
	if err != nil {
		return err
	}

	// Copy each file and directory
	for _, fileInfo := range fileInfos {
		srcPath := filepath.Join(srcDir, fileInfo.Name())
		destPath := filepath.Join(destDir, fileInfo.Name())

		if fileInfo.IsDir() {
			// Recursively copy subdirectories
			err = copyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err = copyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(srcPath, destPath string) error {
	// Open source file
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// Copy file contents
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}
