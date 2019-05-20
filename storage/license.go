package storage

import (
	"gopkg.in/yaml.v2"
	"log"
)

type License struct {
	Purchased bool
	Valid     bool
	Key       *string
	Username  *string
	Company   *string
}

const licenseFileName = "license.yml"

func GetLicense() *License {
	if !fileExists(licenseFileName) {
		return nil
	} else {
		file := readStorageFile(getPlatformStoragePath(licenseFileName))

		license := new(License)
		err := yaml.Unmarshal(file, license)

		if err != nil {
			log.Fatalf("Can't read license file: %v", err)
		}

		return license
	}
}

func InitLicense() {

	license := License{
		Purchased: false,
		Valid:     false,
	}

	file, err := yaml.Marshal(license)
	if err != nil {
		log.Fatalf("Can't marshal config file: %v", err)
	}

	initStorage()
	saveStorageFile(file, getPlatformStoragePath(licenseFileName))
}

func SaveLicense() {
	// TODO
}
