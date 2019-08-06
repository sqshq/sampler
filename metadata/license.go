package metadata

import (
	"gopkg.in/yaml.v3"
	"log"
)

type License struct {
	Key      *string      `yaml:"k"`
	Username *string      `yaml:"u"`
	Company  *string      `yaml:"c"`
	Type     *LicenseType `yaml:"t"`
	Valid    bool         `yaml:"v"`
}

type LicenseType rune

const (
	TypePersonal   LicenseType = 0
	TypeCommercial LicenseType = 1
)

const licenseFileName = "license.yml"

func GetLicense() *License {

	if !fileExists(licenseFileName) {
		return nil
	}

	file := readStorageFile(getPlatformStoragePath(licenseFileName))

	license := new(License)
	err := yaml.Unmarshal(file, license)

	if err != nil {
		log.Fatalf("Failed to read license file: %v", err)
	}

	return license
}

func SaveLicense(license License) {

	initStorage()

	file, err := yaml.Marshal(license)
	if err != nil {
		log.Fatalf("Failed to marshal license file: %v", err)
	}

	saveStorageFile(file, licenseFileName)
}
