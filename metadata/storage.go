package metadata

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	macOSDir   = "/Library/Application Support/Sampler"
	linuxDir   = "/.config/Sampler"
	windowsDir = "%LOCALAPPDATA%\\Sampler"
)

func fileExists(filename string) bool {
	_, err := os.Stat(getPlatformStoragePath(filename))
	return !os.IsNotExist(err)
}

func getPlatformStoragePath(filename string) string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, macOSDir, filename)
	case "windows":
		return filepath.Join(home, windowsDir, filename)
	default:
		return filepath.Join(linuxDir, filename)
	}
}

func initStorage() {
	err := os.MkdirAll(getPlatformStoragePath(""), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}
}

func readStorageFile(path string) []byte {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to the read storage file: %s", path)
	}

	return file
}

func saveStorageFile(file []byte, fileName string) {
	err := ioutil.WriteFile(getPlatformStoragePath(fileName), file, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to save the storage file: %v", err)
	}
}
