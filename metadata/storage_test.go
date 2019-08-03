package metadata

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"testing"
)

type f struct {
	a int
}

func Test_fileExists(t *testing.T) {

	initStorage()

	_, err := os.Create(getPlatformStoragePath("exists"))
	if err != nil {
		panic(err)
	}

	defer os.Remove(getPlatformStoragePath("exists"))

	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"should verify that file does not exist", args{"does-not-exist"}, false},
		{"should verify that file exists", args{"exists"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.filename); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveStorageFile(t *testing.T) {

	initStorage()

	file, _ := yaml.Marshal(f{a: 1})
	name := "test"

	saveStorageFile(file, name)

	read, _ := ioutil.ReadFile(getPlatformStoragePath(name))

	if !bytes.Equal(file, read) {
		t.Errorf("read file != saved file")
	}
}

func Test_readStorageFile(t *testing.T) {

	initStorage()

	file, _ := yaml.Marshal(f{a: 1})
	name := "test"

	err := ioutil.WriteFile(getPlatformStoragePath(name), file, os.ModePerm)
	if err != nil {
		panic(err)
	}
	read := readStorageFile(getPlatformStoragePath(name))

	if !bytes.Equal(file, read) {
		t.Errorf("read file != saved file")
	}
}
