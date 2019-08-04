package metadata

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

const (
	unixDir    = "/.config/Sampler"
	windowsDir = "Sampler"
)

func fileExists(filename string) bool {
	_, err := os.Stat(getPlatformStoragePath(filename))
	return !os.IsNotExist(err)
}

func getPlatformStoragePath(filename string) string {
	switch runtime.GOOS {
	case "windows":
		cache, _ := os.UserCacheDir()
		return filepath.Join(cache, windowsDir, filename)
	default:
		home, _ := homedir.Dir()
		return filepath.Join(home, unixDir, filename)
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
		log.Fatalf("Failed to save the storage file: %s %v", fileName, err)
	}
}
