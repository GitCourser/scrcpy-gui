package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	want := Config{
		Path:         "/opt/scrcpy",
		IP:           "192.168.0.11",
		Port:         "40015",
		Scan:         true,
		PortRange:    "36000-44000",
		Exclude:      "80,443",
		Threads:      1200,
		Timeout:      2.5,
		AlwaysOnTop:  true,
		StayAwake:    true,
		PowerOff:     true,
		HEVC:         true,
		NoAudio:      true,
		NoMouseHover: true,
		Record:       true,
		RecordAudio:  true,
		Key:          "lctrl",
		Size:         "1200",
		Bitrate:      "4",
		Buffer:       "100",
	}
	if err := SaveConfig(path, want); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip mismatch:\n got %+v\nwant %+v", got, want)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var keys map[string]any
	if err := json.Unmarshal(data, &keys); err != nil {
		t.Fatalf("saved file is not JSON: %v", err)
	}
	if keys["strPath"] != "/opt/scrcpy" {
		t.Fatalf("strPath = %v", keys["strPath"])
	}
	if keys["numThread"] != float64(1200) {
		t.Fatalf("numThread = %v", keys["numThread"])
	}
}

func TestConfigOmitsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	if err := SaveConfig(path, DefaultConfig()); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var keys map[string]any
	if err := json.Unmarshal(data, &keys); err != nil {
		t.Fatalf("saved file is not JSON: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("defaults should produce an empty object, got %v", keys)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatal("config file should end with a newline")
	}
}

func TestLoadMissingConfigUsesDefaults(t *testing.T) {
	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.cfg"))
	if err != nil {
		t.Fatalf("missing file should not fail: %v", err)
	}
	if cfg.Threads != defaultThreads || cfg.Timeout != defaultTimeout {
		t.Fatalf("defaults = %+v", cfg)
	}
}

func TestLoadInvalidConfigReturnsDefaultsAndError(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err == nil {
		t.Fatal("invalid JSON should return an error")
	}
	if cfg.Threads != defaultThreads || cfg.Timeout != defaultTimeout {
		t.Fatalf("defaults = %+v", cfg)
	}
}
