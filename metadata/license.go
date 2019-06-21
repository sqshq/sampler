package metadata

import (
	"gopkg.in/yaml.v3"
	"log"
)

type License struct {
	Key      *string `yaml:"k"`
	Username *string `yaml:"u"`
	Company  *string `yaml:"c"`
	Valid    bool    `yaml:"v"`
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
			log.Fatalf("Failed to read license file: %v", err)
		}

		return license
	}
}

func SaveLicense(license License) {

	file, err := yaml.Marshal(license)
	if err != nil {
		log.Fatalf("Failed to marshal license file: %v", err)
	}

	saveStorageFile(file, licenseFileName)
}
