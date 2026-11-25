package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func writeTool(t *testing.T, dir, name string) {
	t.Helper()
	content := "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		content = ""
	}
	path := filepath.Join(dir, executableName(name))
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func testLaunchInput(t *testing.T) launchInput {
	t.Helper()
	dir := t.TempDir()
	writeTool(t, dir, "scrcpy")
	writeTool(t, dir, "adb")
	return launchInput{
		Path:      dir,
		IP:        "127.0.0.1",
		Port:      "40015",
		RecordDir: filepath.Join(dir, "rec"),
		Stamp:     "2026-01-02_03-04-05",
	}
}

func hasArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

func TestBuildLaunchPlanOptions(t *testing.T) {
	in := testLaunchInput(t)
	in.Scan = true
	in.PortRange = "36000-44000"
	in.AlwaysOnTop = true
	in.StayAwake = true
	in.PowerOff = true
	in.HEVC = true
	in.NoAudio = true
	in.NoMouseHover = true
	in.Key = "lctrl"
	in.Size = "1200"
	in.Bitrate = "4"
	in.Buffer = "100"
	in.Pair = true
	in.PairPort = "37123"
	in.PairCode = "123456"

	plan, err := buildLaunchPlan(in)
	if err != nil {
		t.Fatalf("buildLaunchPlan: %v", err)
	}
	wantArgs := []string{
		"-s127.0.0.1:40015",
		"--max-fps=60",
		"--always-on-top",
		"-Sw",
		"--power-off-on-close",
		"--video-codec=h265",
		"--no-audio",
		"--no-mouse-hover",
		"--shortcut-mod=lctrl",
		"-m1200",
		"-b4M",
		"--display-buffer=100",
	}
	if !reflect.DeepEqual(plan.Args, wantArgs) {
		t.Fatalf("args = %v\nwant %v", plan.Args, wantArgs)
	}
	if plan.Device != "127.0.0.1:40015" {
		t.Fatalf("device = %q", plan.Device)
	}
	if plan.PairAddr != "127.0.0.1:37123" || plan.PairCode != "123456" {
		t.Fatalf("pair = %q %q", plan.PairAddr, plan.PairCode)
	}
	if plan.OnlyAudio {
		t.Fatal("OnlyAudio should be false")
	}
}

func TestBuildLaunchPlanRecordVideo(t *testing.T) {
	in := testLaunchInput(t)
	in.Record = true

	plan, err := buildLaunchPlan(in)
	if err != nil {
		t.Fatalf("buildLaunchPlan: %v", err)
	}
	want := "--record=" + filepath.ToSlash(filepath.Join(in.RecordDir, "2026-01-02_03-04-05.mp4"))
	if !hasArg(plan.Args, want) {
		t.Fatalf("args = %v, want %s", plan.Args, want)
	}
	if info, err := os.Stat(in.RecordDir); err != nil || !info.IsDir() {
		t.Fatalf("record directory was not created: %v", err)
	}
}

func TestBuildLaunchPlanRecordAudio(t *testing.T) {
	in := testLaunchInput(t)
	in.Record = true
	in.RecordAudio = true
	in.StayAwake = true
	in.PowerOff = true

	plan, err := buildLaunchPlan(in)
	if err != nil {
		t.Fatalf("buildLaunchPlan: %v", err)
	}
	if !plan.OnlyAudio {
		t.Fatal("OnlyAudio should be true")
	}
	for _, removed := range []string{"-Sw", "--power-off-on-close"} {
		if hasArg(plan.Args, removed) {
			t.Fatalf("%s should not be passed in audio-only mode: %v", removed, plan.Args)
		}
	}
	for _, required := range []string{
		"--no-video",
		"--audio-buffer=100",
		"--record=" + filepath.ToSlash(filepath.Join(in.RecordDir, "2026-01-02_03-04-05.opus")),
	} {
		if !hasArg(plan.Args, required) {
			t.Fatalf("args = %v, want %s", plan.Args, required)
		}
	}
	for _, arg := range plan.Args {
		if strings.HasSuffix(arg, ".mp4") {
			t.Fatalf("audio-only mode should not record video: %v", plan.Args)
		}
	}
}

func TestBuildLaunchPlanValidation(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*launchInput)
		wantMsg string
	}{
		{"missing ip", func(in *launchInput) { in.IP = "" }, "IP"},
		{"invalid ip", func(in *launchInput) { in.IP = "300.1.1.1" }, "IP 地址无效"},
		{"invalid port", func(in *launchInput) { in.Port = "0" }, "端口无效"},
		{"invalid scan range", func(in *launchInput) { in.Scan = true; in.PortRange = "abc" }, "端口范围"},
		{"missing pair code", func(in *launchInput) { in.Pair = true; in.PairPort = "37123" }, "配对"},
		{"audio only with no audio", func(in *launchInput) { in.RecordAudio = true; in.NoAudio = true }, "仅录音"},
		{"invalid shortcut", func(in *launchInput) { in.Key = "meta" }, "热键无效"},
		{"invalid size", func(in *launchInput) { in.Size = "0" }, "高度无效"},
		{"invalid bitrate", func(in *launchInput) { in.Bitrate = "-1" }, "码率无效"},
		{"invalid buffer", func(in *launchInput) { in.Buffer = "0" }, "缓冲无效"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := testLaunchInput(t)
			test.mutate(&in)
			_, err := buildLaunchPlan(in)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), test.wantMsg) {
				t.Fatalf("error = %q, want it to contain %q", err, test.wantMsg)
			}
		})
	}
}
