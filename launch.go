package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const scrcpyURL = "https://github.com/Genymobile/scrcpy"

// LaunchPlan 包含配对、连接并启动 scrcpy 所需的全部信息。
// 它以 JSON 形式编码进启动 Effect 的键，使后台任务读取不可变快照，
// 而不是实时的 model。
type LaunchPlan struct {
	AdbPath    string   `json:"adbPath"`
	ScrcpyPath string   `json:"scrcpyPath"`
	Device     string   `json:"device"`
	Pair       bool     `json:"pair"`
	PairAddr   string   `json:"pairAddr"`
	PairCode   string   `json:"pairCode"`
	Args       []string `json:"args"`
	OnlyAudio  bool     `json:"onlyAudio"`
}

type launchInput struct {
	Path         string
	IP           string
	Port         string
	Scan         bool
	PortRange    string
	Pair         bool
	PairPort     string
	PairCode     string
	AlwaysOnTop  bool
	StayAwake    bool
	PowerOff     bool
	HEVC         bool
	NoAudio      bool
	NoMouseHover bool
	Record       bool
	RecordAudio  bool
	Key          string
	Size         string
	Bitrate      string
	Buffer       string
	RecordDir    string
	Stamp        string
}

var allowedShortcutMods = map[string]bool{
	"lctrl": true, "rctrl": true,
	"lalt": true, "ralt": true,
	"lsuper": true, "rsuper": true,
}

func executableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func isToolFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}

// findTool 在指定了 scrcpy 目录时于该目录查找；
// 目录为空时则依次查找可执行文件所在目录与 PATH。
func findTool(dir, name string) (string, error) {
	fileName := executableName(name)
	if strings.TrimSpace(dir) != "" {
		candidate := filepath.Join(strings.TrimSpace(dir), fileName)
		if isToolFile(candidate) {
			return candidate, nil
		}
		return "", fmt.Errorf("未在 %s 中找到 %s", dir, fileName)
	}

	var candidates []string
	if executable, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(executable), fileName))
	}
	if path, err := exec.LookPath(fileName); err == nil {
		candidates = append(candidates, path)
	}
	if runtime.GOOS != "windows" {
		if path, err := exec.LookPath(name); err == nil {
			candidates = append(candidates, path)
		}
	}
	for _, candidate := range candidates {
		if isToolFile(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("未找到 %s，请设置 scrcpy 目录或将它加入 PATH", fileName)
}

func buildLaunchPlan(in launchInput) (*LaunchPlan, error) {
	ip := strings.TrimSpace(in.IP)
	portText := strings.TrimSpace(in.Port)
	if ip == "" || portText == "" {
		return nil, errors.New("请先设置 IP 和端口")
	}
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("IP 地址无效: %s", ip)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("端口无效: %s", portText)
	}
	if in.Scan {
		if _, _, err := parsePortRange(in.PortRange); err != nil {
			return nil, fmt.Errorf("扫描已开启: %w", err)
		}
	}
	pairPort := strings.TrimSpace(in.PairPort)
	pairCode := strings.TrimSpace(in.PairCode)
	pairPortNumber := 0
	if in.Pair {
		if pairPort == "" || pairCode == "" {
			return nil, errors.New("配对已开启，请填写配对端口和配对码")
		}
		number, err := strconv.Atoi(pairPort)
		if err != nil || number < 1 || number > 65535 {
			return nil, fmt.Errorf("配对端口无效: %s", pairPort)
		}
		pairPortNumber = number
	}
	if in.RecordAudio && in.NoAudio {
		return nil, errors.New("“仅录音”与“禁用音频转发”不能同时开启")
	}

	scrcpyPath, err := findTool(in.Path, "scrcpy")
	if err != nil {
		return nil, fmt.Errorf("未找到 scrcpy，请设置正确的 scrcpy 目录: %w", err)
	}
	adbPath, err := findTool(in.Path, "adb")
	if err != nil {
		return nil, fmt.Errorf("未找到 adb，请将它和 scrcpy 放在同一目录或加入 PATH: %w", err)
	}

	device := net.JoinHostPort(ip, strconv.Itoa(port))
	args := []string{"-s" + device, "--max-fps=60"}
	if in.AlwaysOnTop {
		args = append(args, "--always-on-top")
	}
	if in.StayAwake && !in.RecordAudio {
		args = append(args, "-Sw")
	}
	if in.PowerOff && !in.RecordAudio {
		args = append(args, "--power-off-on-close")
	}
	if in.HEVC {
		args = append(args, "--video-codec=h265")
	}
	if in.NoAudio {
		args = append(args, "--no-audio")
	}
	if in.NoMouseHover {
		args = append(args, "--no-mouse-hover")
	}

	recordDir := strings.TrimSpace(in.RecordDir)
	if recordDir == "" {
		recordDir = "rec"
	}
	stamp := in.Stamp
	if stamp == "" {
		stamp = time.Now().Format("2006-01-02_15-04-05")
	}
	onlyAudio := in.RecordAudio
	if in.RecordAudio || in.Record {
		if err := os.MkdirAll(recordDir, 0o755); err != nil {
			return nil, fmt.Errorf("创建录制目录 %s 失败: %w", recordDir, err)
		}
	}
	switch {
	case in.RecordAudio:
		args = append(args,
			"--no-video",
			"--audio-buffer=100",
			"--record="+filepath.ToSlash(filepath.Join(recordDir, stamp+".opus")),
		)
	case in.Record:
		args = append(args, "--record="+filepath.ToSlash(filepath.Join(recordDir, stamp+".mp4")))
	}

	key := strings.TrimSpace(in.Key)
	if key != "" {
		if !allowedShortcutMods[key] {
			return nil, fmt.Errorf("热键无效: %s（可选 lctrl、rctrl、lalt、ralt、lsuper、rsuper）", key)
		}
		args = append(args, "--shortcut-mod="+key)
	}
	size := strings.TrimSpace(in.Size)
	if size != "" {
		number, err := strconv.Atoi(size)
		if err != nil || number < 1 {
			return nil, fmt.Errorf("高度无效: %s", size)
		}
		args = append(args, "-m"+size)
	}
	bitrate := strings.TrimSpace(in.Bitrate)
	if bitrate != "" {
		number, err := strconv.ParseFloat(bitrate, 64)
		if err != nil || number <= 0 {
			return nil, fmt.Errorf("码率无效: %s", bitrate)
		}
		args = append(args, "-b"+bitrate+"M")
	}
	buffer := strings.TrimSpace(in.Buffer)
	if buffer != "" {
		number, err := strconv.Atoi(buffer)
		if err != nil || number < 1 {
			return nil, fmt.Errorf("缓冲无效: %s", buffer)
		}
		args = append(args, "--display-buffer="+buffer)
	}

	plan := &LaunchPlan{
		AdbPath:    adbPath,
		ScrcpyPath: scrcpyPath,
		Device:     device,
		Pair:       in.Pair,
		Args:       args,
		OnlyAudio:  onlyAudio,
	}
	if in.Pair {
		plan.PairAddr = net.JoinHostPort(ip, strconv.Itoa(pairPortNumber))
		plan.PairCode = pairCode
	}
	return plan, nil
}

func executeLaunch(ctx context.Context, plan LaunchPlan) error {
	if plan.Pair {
		pair := exec.CommandContext(ctx, plan.AdbPath, "pair", plan.PairAddr)
		pair.Stdin = strings.NewReader(plan.PairCode + "\n")
		prepareCommand(pair)
		_ = pair.Run()
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	devices := exec.CommandContext(ctx, plan.AdbPath, "devices")
	prepareCommand(devices)
	output, err := devices.CombinedOutput()
	if err != nil {
		return fmt.Errorf("运行 adb devices 失败: %s", strings.TrimSpace(string(output)))
	}

	connected := false
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == plan.Device && fields[1] == "device" {
			connected = true
			break
		}
	}
	if !connected {
		connect := exec.CommandContext(ctx, plan.AdbPath, "connect", plan.Device)
		prepareCommand(connect)
		output, err := connect.CombinedOutput()
		if err != nil {
			return fmt.Errorf("连接设备 %s 失败: %s", plan.Device, strings.TrimSpace(string(output)))
		}
	}

	scrcpy, err := startScrcpyProcess(plan.ScrcpyPath, plan.Args)
	if err != nil {
		return fmt.Errorf("启动 scrcpy 失败: %w", err)
	}
	exited := make(chan error, 1)
	go func() { exited <- scrcpy.Wait() }()
	select {
	case err := <-exited:
		if err != nil {
			return fmt.Errorf("scrcpy 启动后立即退出: %w", err)
		}
		return errors.New("scrcpy 启动后立即退出")
	case <-time.After(700 * time.Millisecond):
		go func() { <-exited }()
		return nil
	}
}

// scrcpyRunning 返回当前是否已有 scrcpy 进程在运行，
// 对应原程序的单实例检查。
func scrcpyRunning() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq scrcpy.exe", "/NH")
		prepareCommand(cmd)
		output, err := cmd.Output()
		return err == nil && strings.Contains(strings.ToLower(string(output)), "scrcpy.exe")
	}
	if _, err := exec.LookPath("pgrep"); err != nil {
		return false
	}
	output, err := exec.Command("pgrep", "-x", "scrcpy").Output()
	return err == nil && len(strings.TrimSpace(string(output))) > 0
}
