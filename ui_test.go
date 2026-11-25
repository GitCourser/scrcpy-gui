package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vivaui/viva/app"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

func newTestModel(t *testing.T) *model {
	t.Helper()
	return newModel(DefaultConfig(), filepath.Join(t.TempDir(), configFileName))
}

func newTestDriver(t *testing.T, m *model) *app.Driver {
	t.Helper()
	theme := themeForAppearance(app.AppearanceLight)
	page := ui.NewPage(ui.PageConfig{View: m.page, Theme: &theme})
	driver, err := app.NewDriver(ui.Configure(page, app.Config{
		Title:  appVersion,
		Width:  windowWidth,
		Height: windowHeight,
	}), nil)
	if err != nil {
		t.Fatalf("NewDriver: %v", err)
	}
	t.Cleanup(driver.Close)
	return driver
}

func clickLabel(t *testing.T, driver *app.Driver, label string) {
	t.Helper()
	for _, hit := range driver.Snapshot().Hits {
		if hit.Label == label && !hit.IsLabel {
			bounds := hit.Bounds
			if err := driver.Click(bounds.X+bounds.Width/2, bounds.Y+bounds.Height/2); err != nil {
				t.Fatalf("click %q: %v", label, err)
			}
			return
		}
	}
	t.Fatalf("no control named %q", label)
}

func hasLabel(snapshot app.Snapshot, text string) bool {
	for _, label := range snapshot.Labels {
		if label.Value == text {
			return true
		}
	}
	return false
}

func hitBounds(snapshot app.Snapshot, label string) (gfx.Rect, bool) {
	for _, hit := range snapshot.Hits {
		if hit.Label == label && !hit.IsLabel {
			return hit.Bounds, true
		}
	}
	return gfx.Rect{}, false
}

func findLabel(snapshot app.Snapshot, value string) (app.Label, bool) {
	for _, label := range snapshot.Labels {
		if label.Value == value {
			return label, true
		}
	}
	return app.Label{}, false
}

func labelBounds(snapshot app.Snapshot, value string) (gfx.Rect, bool) {
	label, ok := findLabel(snapshot, value)
	return label.Bounds, ok
}

func waitFor(t *testing.T, driver *app.Driver, done func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !done() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for the effect")
		}
		if err := driver.Paint(); err != nil {
			t.Fatalf("paint: %v", err)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestPageRendersControls(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	snapshot := driver.Snapshot()
	if snapshot.Width != windowWidth || snapshot.Height != windowHeight {
		t.Fatalf("window = %dx%d", snapshot.Width, snapshot.Height)
	}
	for _, want := range []string{
		"连接", "扫描", "配对", "选项",
		"选择目录", "启 动",
		"置顶", "熄屏且禁止休眠", "退出后锁屏", "录屏",
		"H.265", "禁用音频转发", "禁用鼠标悬停", "仅录音",
		"热键", "高度", "码率(Mbps)", "缓冲(毫秒)",
	} {
		if !hasLabel(snapshot, want) {
			t.Errorf("missing label %q", want)
		}
	}
	if _, ok := hitBounds(snapshot, "扫描端口"); !ok {
		t.Error("扫描端口 icon button is missing")
	}
}

func TestPageFitsInsideTheWindow(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	for _, hit := range driver.Snapshot().Hits {
		bounds := hit.Bounds
		if bounds.Width <= 0 || bounds.Height <= 0 {
			continue
		}
		if bounds.X < -1 || bounds.Y < -1 ||
			bounds.X+bounds.Width > windowWidth+1 ||
			bounds.Y+bounds.Height > windowHeight+1 {
			t.Errorf("control %q is outside the window: %+v", hit.Label, bounds)
		}
	}
}

func TestScanCheckboxTogglesScanFields(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	if m.scanEnabled {
		t.Fatal("scan should start disabled")
	}
	clickLabel(t, driver, "扫描")
	if !m.scanEnabled {
		t.Fatal("clicking 扫描 should enable scanning")
	}
}

func TestStartWithoutAddressShowsAlert(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	clickLabel(t, driver, "启 动")
	if !m.showAlert {
		t.Fatal("missing IP and port should show an alert")
	}
	if !strings.Contains(m.alert, "IP") {
		t.Fatalf("alert = %q", m.alert)
	}
	if !hasLabel(driver.Snapshot(), "知道了") {
		t.Fatal("alert should offer 知道了")
	}
}

func TestAlertCanBeDismissed(t *testing.T) {
	m := newTestModel(t)
	m.alert = "测试提示"
	m.showAlert = true
	driver := newTestDriver(t, m)
	if !hasLabel(driver.Snapshot(), "测试提示") {
		t.Fatal("alert text is missing")
	}
	clickLabel(t, driver, "知道了")
	if m.showAlert {
		t.Fatal("知道了 should close the alert")
	}
}

func TestScanEffectUpdatesOpenPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	m := newTestModel(t)
	m.ip = "127.0.0.1"
	m.port = "1"
	m.portRange = fmt.Sprintf("%d-%d", port, port)
	m.threads = "4"
	m.timeout = "0.3"
	driver := newTestDriver(t, m)

	m.startScan(false)
	if !m.scanning {
		t.Fatal("scan should be running")
	}
	waitFor(t, driver, func() bool { return !m.scanning })
	if m.showAlert {
		t.Fatalf("scan alert: %s", m.alert)
	}
	if m.port != strconv.Itoa(port) {
		t.Fatalf("port = %q, want %d", m.port, port)
	}
	if !strings.Contains(m.status, "扫描成功") {
		t.Fatalf("status = %q", m.status)
	}
}

func TestStartEffectLaunchesScrcpy(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses shell scripts")
	}
	dir := t.TempDir()
	adb := "#!/bin/sh\ncase \"$1\" in\ndevices) printf '127.0.0.1:40015\\tdevice\\n' ;;\nesac\nexit 0\n"
	if err := os.WriteFile(filepath.Join(dir, "adb"), []byte(adb), 0o755); err != nil {
		t.Fatal(err)
	}
	scrcpy := "#!/bin/sh\nsleep 1\nexit 0\n"
	if err := os.WriteFile(filepath.Join(dir, "scrcpy"), []byte(scrcpy), 0o755); err != nil {
		t.Fatal(err)
	}

	m := newTestModel(t)
	m.path = dir
	m.ip = "127.0.0.1"
	m.port = "40015"
	driver := newTestDriver(t, m)

	m.start()
	if !m.busy {
		t.Fatal("start should set busy")
	}
	waitFor(t, driver, func() bool { return !m.busy })
	if m.showAlert {
		t.Fatalf("start alert: %s", m.alert)
	}
	if !strings.Contains(m.status, "已启动") {
		t.Fatalf("status = %q", m.status)
	}
}

func TestOptionsSpreadEvenlyAcrossTheRow(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	snapshot := driver.Snapshot()
	for _, row := range [][]string{
		{"置顶", "熄屏且禁止休眠", "退出后锁屏", "录屏"},
		{"H.265", "禁用音频转发", "禁用鼠标悬停", "仅录音"},
	} {
		positions := make([]float32, 0, len(row))
		for _, want := range row {
			found := false
			for _, label := range snapshot.Labels {
				if label.Value == want {
					positions = append(positions, label.Bounds.X)
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("missing label %q", want)
			}
		}
		step := positions[1] - positions[0]
		for i := 2; i < len(positions); i++ {
			if gap := positions[i] - positions[i-1]; math.Abs(float64(gap-step)) > 1.5 {
				t.Fatalf("%v are not evenly spread: %v", row, positions)
			}
		}
		if positions[len(positions)-1] < windowWidth/2 {
			t.Fatalf("%v stay on the left: %v", row, positions)
		}
	}
}

func TestEnterStartsFromAnywhere(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	if err := driver.Key(gfx.KeyEnter, 0); err != nil {
		t.Fatal(err)
	}
	if !m.showAlert || !strings.Contains(m.alert, "IP") {
		t.Fatalf("回车未触发启动：alert = %q", m.alert)
	}
}

func TestEnterInAFieldStarts(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	clickLabel(t, driver, "端口")
	if err := driver.Type("40015"); err != nil {
		t.Fatal(err)
	}
	if err := driver.Key(gfx.KeyEnter, 0); err != nil {
		t.Fatal(err)
	}
	if m.port != "40015" {
		t.Fatalf("field lost the typed text: %q", m.port)
	}
	if !m.showAlert || !strings.Contains(m.alert, "IP") {
		t.Fatalf("输入框内回车未触发启动：alert = %q", m.alert)
	}
}

func TestPathFromDropUsesDirectory(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "scrcpy")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := pathFromDrop(file); got != dir {
		t.Fatalf("pathFromDrop(file) = %q, want %q", got, dir)
	}
	if got := pathFromDrop(dir); got != dir {
		t.Fatalf("pathFromDrop(dir) = %q, want %q", got, dir)
	}
}

func TestThemeFollowsSystemAppearance(t *testing.T) {
	dark := themeForAppearance(app.AppearanceDark)
	light := themeForAppearance(app.AppearanceLight)
	if dark.Background != gfx.DarkTheme().Background {
		t.Fatal("dark appearance should use the dark palette")
	}
	if light.Background != gfx.LightTheme().Background {
		t.Fatal("light appearance should use the light palette")
	}
	if dark.Background == light.Background {
		t.Fatal("light and dark themes should differ")
	}
	if dark.Metrics.ControlHeight != 30 {
		t.Fatalf("control height = %v", dark.Metrics.ControlHeight)
	}
	unknown := themeForAppearance(app.AppearanceUnknown)
	if unknown.Background != light.Background {
		t.Fatal("an unknown appearance should fall back to the light palette")
	}
}

func TestFollowSystemAppearanceSwitchesTheme(t *testing.T) {
	driver, err := app.NewDriver(app.Config{Title: "appearance", Width: 100, Height: 100}, nil)
	if err != nil {
		t.Fatalf("NewDriver: %v", err)
	}
	defer driver.Close()

	driver.SetSystemAppearance(app.AppearanceLight)
	theme := themeForAppearance(app.SystemAppearance())
	stop := followSystemAppearance(&theme, new(app.Invalidator))
	defer stop()

	if theme.Background != gfx.LightTheme().Background {
		t.Fatal("the window should start light")
	}
	driver.SetSystemAppearance(app.AppearanceDark)
	if theme.Background != gfx.DarkTheme().Background {
		t.Fatal("switching the system to dark should switch the window")
	}
	driver.SetSystemAppearance(app.AppearanceLight)
	if theme.Background != gfx.LightTheme().Background {
		t.Fatal("switching the system back to light should switch the window")
	}
	stop()
	driver.SetSystemAppearance(app.AppearanceDark)
	if theme.Background != gfx.LightTheme().Background {
		t.Fatal("a stopped watch should leave the theme alone")
	}
}

func textPainted(snapshot app.Snapshot, value string) bool {
	for _, text := range snapshot.Texts {
		if text.Value == value {
			return true
		}
	}
	return false
}

func TestScanButtonIsAnIconInTheHeader(t *testing.T) {
	m := newTestModel(t)
	driver := newTestDriver(t, m)
	snapshot := driver.Snapshot()

	bounds, ok := hitBounds(snapshot, "扫描端口")
	if !ok {
		t.Fatal("扫描图标按钮缺少可访问名称")
	}
	if hasLabel(snapshot, "扫描端口") {
		t.Fatal("扫描按钮不应再绘制文字")
	}
	if math.Abs(float64(bounds.Width-bounds.Height)) > 1 {
		t.Fatalf("扫描按钮应为方形图标按钮: %+v", bounds)
	}
	if bounds.X < windowWidth/2 {
		t.Fatalf("扫描按钮应位于右上角: %+v", bounds)
	}
	checkbox, ok := hitBounds(snapshot, "扫描")
	if !ok {
		t.Fatal("扫描复选框缺少可访问名称")
	}
	if bounds.Y+bounds.Height < checkbox.Y || checkbox.Y+checkbox.Height < bounds.Y {
		t.Fatalf("扫描按钮不在标题行: button %+v, checkbox %+v", bounds, checkbox)
	}
}

func TestRecordLabelTurnsRed(t *testing.T) {
	theme := themeForAppearance(app.AppearanceLight)

	off := newTestModel(t)
	offLabel, ok := findLabel(newTestDriver(t, off).Snapshot(), "录屏")
	if !ok {
		t.Fatal("missing 录屏 label")
	}
	if offLabel.Color == theme.Danger {
		t.Fatal("未勾选时录屏标签不应为红色")
	}

	m := newTestModel(t)
	m.record = true
	label, ok := findLabel(newTestDriver(t, m).Snapshot(), "录屏")
	if !ok {
		t.Fatal("missing 录屏 label")
	}
	if label.Color != theme.Danger {
		t.Fatalf("录屏标签颜色 = %+v, 期望 %+v", label.Color, theme.Danger)
	}
}

func TestAddressFieldsHaveNoPlaceholder(t *testing.T) {
	m := newTestModel(t)
	snapshot := newTestDriver(t, m).Snapshot()
	for _, placeholder := range []string{"192.168.0.11", "40015", "1100"} {
		if textPainted(snapshot, placeholder) {
			t.Errorf("输入框仍显示默认提示 %q", placeholder)
		}
	}
	if !textPainted(snapshot, "36000-44000") {
		t.Error("扫描框的端口范围提示应保留")
	}
}

func TestConnectionAndPairFieldsSpreadVertically(t *testing.T) {
	m := newTestModel(t)
	snapshot := newTestDriver(t, m).Snapshot()

	scanTop, ok := labelBounds(snapshot, "端口范围")
	if !ok {
		t.Fatal("missing 端口范围 label")
	}
	scanBottom, ok := labelBounds(snapshot, "屏蔽端口")
	if !ok {
		t.Fatal("missing 屏蔽端口 label")
	}
	scanPitch := scanBottom.Y - scanTop.Y

	for _, fields := range []struct{ first, second string }{
		{"IP", "端口"},
		{"配对端口", "配对码"},
	} {
		first, ok := labelBounds(snapshot, fields.first)
		if !ok {
			t.Fatalf("missing %q label", fields.first)
		}
		second, ok := labelBounds(snapshot, fields.second)
		if !ok {
			t.Fatalf("missing %q label", fields.second)
		}
		spread := second.Y - first.Y
		if spread < scanPitch+8 {
			t.Errorf("%q 与 %q 的纵向间距 %v 未拉大（扫描框基准 %v）", fields.first, fields.second, spread, scanPitch)
		} else {
			t.Logf("%q 与 %q 的纵向间距 %.1f（扫描框基准 %.1f）", fields.first, fields.second, spread, scanPitch)
		}
	}
}
