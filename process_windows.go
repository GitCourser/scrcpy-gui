//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

const createNoWindow = 0x08000000

// prepareCommand 让子进程在没有可见控制台的情况下运行，
// 效果等同于原来的 CREATE_NO_WINDOW 标志。
// 不设置 HideWindow：它把 SW_HIDE 作为子进程首个窗口的显示状态，
// 会把 scrcpy 的画面窗口一并隐藏。
func prepareCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: createNoWindow,
	}
}

// startScrcpyProcess 用 CREATE_NO_WINDOW 启动 scrcpy：
// 它有自己的控制台，启动器退出后继续运行，画面窗口正常显示。
func startScrcpyProcess(path string, args []string) (*exec.Cmd, error) {
	cmd := exec.Command(path, args...)
	prepareCommand(cmd)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// detachConsole 断开启动时继承的控制台，
// 即使程序由终端启动也不会留下一个终端窗口。
func detachConsole() {
	freeConsole := syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole")
	_, _, _ = freeConsole.Call()
}
