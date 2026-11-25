//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

func prepareCommand(*exec.Cmd) {}

// startScrcpyProcess 在独立会话中启动 scrcpy，启动器退出后它继续运行。
func startScrcpyProcess(path string, args []string) (*exec.Cmd, error) {
	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// detachConsole 在非 Windows 平台无需处理。
func detachConsole() {}
