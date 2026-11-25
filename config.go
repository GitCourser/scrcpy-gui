package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configFileName = "scrcpy.cfg"
	defaultThreads = 400
	defaultTimeout = 1.5
)

// Config 对应 scrcpy.cfg 中的配置项。
type Config struct {
	Path         string  `json:"strPath"`
	IP           string  `json:"strIP"`
	Port         string  `json:"strPort"`
	Scan         bool    `json:"isScan"`
	PortRange    string  `json:"strPortRange"`
	Exclude      string  `json:"strUnPort"`
	Threads      int     `json:"numThread"`
	Timeout      float64 `json:"numTimeout"`
	AlwaysOnTop  bool    `json:"isTop"`
	StayAwake    bool    `json:"isSw"`
	PowerOff     bool    `json:"isPoweroff"`
	HEVC         bool    `json:"isHEVC"`
	NoAudio      bool    `json:"isNoAudio"`
	NoMouseHover bool    `json:"isNoMousehover"`
	Record       bool    `json:"isRecord"`
	RecordAudio  bool    `json:"isRecordAudio"`
	Key          string  `json:"strKey"`
	Size         string  `json:"strSize"`
	Bitrate      string  `json:"strBitrate"`
	Buffer       string  `json:"strBuffer"`
}

func DefaultConfig() Config {
	return Config{Threads: defaultThreads, Timeout: defaultTimeout}
}

// LoadConfig 读取 path，为文件中缺失的键补上默认值。
// 文件不存在不算错误：这表示首次运行。
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("读取配置 %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("解析配置 %s: %w", path, err)
	}
	if cfg.Threads <= 0 {
		cfg.Threads = defaultThreads
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	return cfg, nil
}

// SaveConfig 只写入与默认值不同的项，配置文件保持精简。
func SaveConfig(path string, cfg Config) error {
	values := make(map[string]any)
	if cfg.Path != "" {
		values["strPath"] = cfg.Path
	}
	if cfg.IP != "" {
		values["strIP"] = cfg.IP
	}
	if cfg.Port != "" {
		values["strPort"] = cfg.Port
	}
	if cfg.Scan {
		values["isScan"] = true
	}
	if cfg.PortRange != "" {
		values["strPortRange"] = cfg.PortRange
	}
	if cfg.Exclude != "" {
		values["strUnPort"] = cfg.Exclude
	}
	if cfg.Threads > 0 && cfg.Threads != defaultThreads {
		values["numThread"] = cfg.Threads
	}
	if cfg.Timeout > 0 && cfg.Timeout != defaultTimeout {
		values["numTimeout"] = cfg.Timeout
	}
	for key, on := range map[string]bool{
		"isTop":          cfg.AlwaysOnTop,
		"isSw":           cfg.StayAwake,
		"isPoweroff":     cfg.PowerOff,
		"isHEVC":         cfg.HEVC,
		"isNoAudio":      cfg.NoAudio,
		"isNoMousehover": cfg.NoMouseHover,
		"isRecord":       cfg.Record,
		"isRecordAudio":  cfg.RecordAudio,
	} {
		if on {
			values[key] = true
		}
	}
	if cfg.Key != "" {
		values["strKey"] = cfg.Key
	}
	if cfg.Size != "" {
		values["strSize"] = cfg.Size
	}
	if cfg.Bitrate != "" {
		values["strBitrate"] = cfg.Bitrate
	}
	if cfg.Buffer != "" {
		values["strBuffer"] = cfg.Buffer
	}

	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return fmt.Errorf("生成配置: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("写入配置 %s: %w", path, err)
	}
	return nil
}

// resolveConfigPath 优先使用工作目录中的 scrcpy.cfg，其次是可执行文件所在目录，
// 都不存在时回到工作目录创建新文件。
func resolveConfigPath() string {
	if abs, err := filepath.Abs(configFileName); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}
	if executable, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(executable), configFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	abs, err := filepath.Abs(configFileName)
	if err != nil {
		return configFileName
	}
	return abs
}
