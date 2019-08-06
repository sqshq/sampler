package metadata

import (
	"os"
	"testing"
)

func Test_getEmptyLicense(t *testing.T) {

	cleanupPlatformStorage()
	license := GetLicense()

	if license != nil {
		t.Errorf("expected to be nil")
	}
}

func Test_saveAndGetExistingLicense(t *testing.T) {

	cleanupPlatformStorage()

	original := License{
		Valid: true,
	}

	SaveLicense(original)

	retrieved := *GetLicense()

	if original != retrieved {
		t.Errorf("read file != saved file")
	}
}

func cleanupPlatformStorage() {
	_ = os.RemoveAll(getPlatformStoragePath(""))
	_ = os.Remove(getPlatformStoragePath(""))
}
