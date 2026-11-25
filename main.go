package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivaui/viva/app"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

func main() {
	detachConsole()
	cfgPath := resolveConfigPath()
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	m := newModel(cfg, cfgPath)
	if scrcpyRunning() {
		m.alert = "scrcpy 已运行，不要重复执行。"
		m.showAlert = true
		m.quitOnAck = true
	}

	theme := themeForAppearance(app.SystemAppearance())
	invalidator := new(app.Invalidator)
	var stopAppearance func()
	err = app.RunApplication(func(application *app.Application) error {
		m.application = application
		window, err := ui.OpenWindow(application, app.Config{
			Title:      appVersion,
			Width:      windowWidth,
			Height:     windowHeight,
			MinWidth:   windowWidth,
			MinHeight:  windowHeight,
			OnFileDrop: m.handleFileDrop,
			OnClosed: func() {
				if stopAppearance != nil {
					stopAppearance()
				}
			},
		}, ui.PageConfig{View: m.page, Theme: &theme, Invalidator: invalidator})
		if err != nil {
			return err
		}
		m.window = window
		_ = window.SetTitle(appVersion)
		stopAppearance = followSystemAppearance(&theme, invalidator)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// themeForAppearance 按系统外观选择浅色或深色调色板，
// 并保留本窗口更紧凑的尺寸度量。
func themeForAppearance(appearance app.Appearance) gfx.Theme {
	theme := gfx.LightTheme()
	if appearance == app.AppearanceDark {
		theme = gfx.DarkTheme()
	}
	return theme.WithControlHeight(30).WithFontSize(13)
}

// followSystemAppearance 在系统于浅色与深色外观之间切换时重写 theme 并唤醒窗口。
// 调用方通过 themeForAppearance 设置初始主题；返回值用于停止监听。
func followSystemAppearance(theme *gfx.Theme, invalidator *app.Invalidator) func() {
	current := app.SystemAppearance()
	return app.WatchSystemAppearance(func() {
		next := app.SystemAppearance()
		if next == current {
			return
		}
		current = next
		*theme = themeForAppearance(next)
		invalidator.Invalidate()
	})
}
