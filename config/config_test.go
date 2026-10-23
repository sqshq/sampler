package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConfigPath_ExplicitExistingFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "sampler-test-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName)
	tempFile.Close()

	resolved, err := ResolveConfigPathWithDirs(&fileName, []string{"."})
	if err != nil {
		t.Fatalf("unexpected error resolving existing file: %v", err)
	}
	if resolved != filepath.Clean(tempFile.Name()) {
		t.Errorf("expected %s, got %s", tempFile.Name(), resolved)
	}
}

func TestResolveConfigPath_RelativeFileInSearchDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "sampler-search-dir-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configName := "custom-config.yml"
	targetPath := filepath.Join(tempDir, configName)
	if err := ioutil.WriteFile(targetPath, []byte("theme: dark\n"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	relativeName := configName
	resolved, err := ResolveConfigPathWithDirs(&relativeName, []string{tempDir})
	if err != nil {
		t.Fatalf("unexpected error resolving relative file in search dir: %v", err)
	}
	if resolved != targetPath {
		t.Errorf("expected %s, got %s", targetPath, resolved)
	}
}

func TestResolveConfigPath_DefaultCandidates(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "sampler-default-dir-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetPath := filepath.Join(tempDir, "example.yml")
	if err := ioutil.WriteFile(targetPath, []byte("theme: dark\n"), 0644); err != nil {
		t.Fatalf("failed to write example.yml: %v", err)
	}

	resolved, err := ResolveConfigPathWithDirs(nil, []string{tempDir})
	if err != nil {
		t.Fatalf("unexpected error resolving default candidate: %v", err)
	}
	if resolved != targetPath {
		t.Errorf("expected %s, got %s", targetPath, resolved)
	}
}

func TestResolveConfigPath_NotFound(t *testing.T) {
	nonExistent := "non_existent_file.yml"
	_, err := ResolveConfigPathWithDirs(&nonExistent, []string{"."})
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "config file not found") {
		t.Errorf("expected 'config file not found' in error, got: %v", err)
	}
}

func TestResolveConfigPath_NoConfigAndNoDefaultFound(t *testing.T) {
	emptyDir, err := ioutil.TempDir("", "sampler-empty-*")
	if err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	_, err = ResolveConfigPathWithDirs(nil, []string{emptyDir})
	if err == nil {
		t.Fatal("expected error when no config and no defaults exist, got nil")
	}
	if !strings.Contains(err.Error(), "Please specify config file using --config flag") {
		t.Errorf("expected flag guidance in error, got: %v", err)
	}
}

func TestLoadConfigFromPath_ExampleYml(t *testing.T) {
	examplePath := "../example.yml"
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("example.yml not found at expected relative path")
	}

	cfg, err := LoadConfigFromPath(examplePath)
	if err != nil {
		t.Fatalf("failed to load example.yml: %v", err)
	}

	if len(cfg.RunCharts) == 0 {
		t.Error("expected runcharts in example.yml, got 0")
	}
	if len(cfg.BarCharts) == 0 {
		t.Error("expected barcharts in example.yml, got 0")
	}
	if len(cfg.SparkLines) == 0 {
		t.Error("expected sparklines in example.yml, got 0")
	}
	if len(cfg.Gauges) == 0 {
		t.Error("expected gauges in example.yml, got 0")
	}
	if len(cfg.TextBoxes) == 0 {
		t.Error("expected textboxes in example.yml, got 0")
	}
	if len(cfg.AsciiBoxes) == 0 {
		t.Error("expected asciiboxes in example.yml, got 0")
	}
}

func TestLoadConfigFromPath_NonExistent(t *testing.T) {
	_, err := LoadConfigFromPath("does-not-exist.yml")
	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	yamlContent := `
runcharts:
  - title: Chart A
    rate-ms: 500
    items:
      - label: item1
        sample: echo 1
`
	cfg, err := ParseConfig([]byte(yamlContent))
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
	if len(cfg.RunCharts) != 1 {
		t.Fatalf("expected 1 runchart, got %d", len(cfg.RunCharts))
	}
	if *cfg.RunCharts[0].RateMs != 500 {
		t.Errorf("expected rate-ms 500, got %d", *cfg.RunCharts[0].RateMs)
	}
}

func TestValidate_DuplicateTitles(t *testing.T) {
	yamlContent := `
runcharts:
  - title: Same Title
    items:
      - label: item1
        sample: echo 1
  - title: Same Title
    items:
      - label: item2
        sample: echo 2
`
	_, err := ParseConfig([]byte(yamlContent))
	if err == nil {
		t.Fatal("expected error for duplicate titles, got nil")
	}
	if !strings.Contains(err.Error(), "component titles should be unique") {
		t.Errorf("expected unique title error, got: %v", err)
	}
}

func TestValidate_DuplicateLabels(t *testing.T) {
	yamlContent := `
runcharts:
  - title: Chart A
    items:
      - label: dup
        sample: echo 1
      - label: dup
        sample: echo 2
`
	_, err := ParseConfig([]byte(yamlContent))
	if err == nil {
		t.Fatal("expected error for duplicate labels, got nil")
	}
	if !strings.Contains(err.Error(), "item labels should be unique") {
		t.Errorf("expected unique label error, got: %v", err)
	}
}

func TestValidate_MissingSampleScript(t *testing.T) {
	yamlContent := `
runcharts:
  - title: Chart A
    items:
      - label: item1
`
	_, err := ParseConfig([]byte(yamlContent))
	if err == nil {
		t.Fatal("expected error for missing sample script, got nil")
	}
	if !strings.Contains(err.Error(), "sample script should be specified") {
		t.Errorf("expected missing sample script error, got: %v", err)
	}
}

func TestValidate_BothInitScripts(t *testing.T) {
	yamlContent := `
runcharts:
  - title: Chart A
    items:
      - label: item1
        sample: echo 1
        init: echo start
        multistep-init:
          - echo step 1
`
	_, err := ParseConfig([]byte(yamlContent))
	if err == nil {
		t.Fatal("expected error for both init scripts, got nil")
	}
	if !strings.Contains(err.Error(), "both init and multistep-init scripts are not allowed") {
		t.Errorf("expected both init scripts error, got: %v", err)
	}
}

func TestValidate_InvalidYaml(t *testing.T) {
	yamlContent := `runcharts: [invalid yaml`
	_, err := ParseConfig([]byte(yamlContent))
	if err == nil {
		t.Fatal("expected error for invalid yaml, got nil")
	}
}
