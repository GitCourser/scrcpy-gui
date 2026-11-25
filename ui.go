package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vivaui/viva/app"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

const (
	appVersion   = "scrcpy GUI v0.8"
	windowWidth  = 640
	windowHeight = 580
)

// scanIcon 是 Lucide 的 search 图标：圆形镜片加斜向手柄。
var scanIcon = gfx.Icon{
	ViewBox: 24,
	Stroke:  2,
	Cap:     gfx.CapRound,
	Join:    gfx.JoinRound,
	Commands: []gfx.PathCommand{
		gfx.MoveTo(11, 3),
		gfx.CubicTo(15.42, 3, 19, 6.58, 19, 11),
		gfx.CubicTo(19, 15.42, 15.42, 19, 11, 19),
		gfx.CubicTo(6.58, 19, 3, 15.42, 3, 11),
		gfx.CubicTo(3, 6.58, 6.58, 3, 11, 3),
		gfx.ClosePath(),
		gfx.MoveTo(16.65, 16.65),
		gfx.LineTo(21, 21),
	},
}

type startRequest struct {
	Seq  int
	Plan string
}

type model struct {
	application *app.Application
	window      *app.Window
	cfgPath     string

	path        string
	ip          string
	port        string
	scanEnabled bool
	portRange   string
	exclude     string
	threads     string
	timeout     string
	pairEnabled bool
	pairPort    string
	pairCode    string

	alwaysOnTop  bool
	stayAwake    bool
	powerOff     bool
	record       bool
	hevc         bool
	noAudio      bool
	noMouseHover bool
	recordAudio  bool
	key          string
	size         string
	bitrate      string
	buffer       string

	status    string
	busy      bool
	scanning  bool
	alert     string
	showAlert bool
	quitOnAck bool

	nextScan   int
	activeScan int
	scanReq    ScanRequest
	scanEffect *ui.EffectState[ScanResult, ScanRequest]

	nextStart   int
	activeStart int
	startReq    startRequest
	startEffect *ui.EffectState[bool, startRequest]
}

func newModel(cfg Config, cfgPath string) *model {
	m := &model{
		cfgPath:      cfgPath,
		path:         cfg.Path,
		ip:           cfg.IP,
		port:         cfg.Port,
		scanEnabled:  cfg.Scan,
		portRange:    cfg.PortRange,
		exclude:      cfg.Exclude,
		threads:      strconv.Itoa(cfg.Threads),
		timeout:      strconv.FormatFloat(cfg.Timeout, 'f', -1, 64),
		pairEnabled:  false,
		pairPort:     "",
		pairCode:     "",
		alwaysOnTop:  cfg.AlwaysOnTop,
		stayAwake:    cfg.StayAwake,
		powerOff:     cfg.PowerOff,
		record:       cfg.Record,
		hevc:         cfg.HEVC,
		noAudio:      cfg.NoAudio,
		noMouseHover: cfg.NoMouseHover,
		recordAudio:  cfg.RecordAudio,
		key:          cfg.Key,
		size:         cfg.Size,
		bitrate:      cfg.Bitrate,
		buffer:       cfg.Buffer,
	}
	m.scanEffect = ui.NewEffectKeyed[ScanResult, ScanRequest](runScan)
	m.startEffect = ui.NewEffectKeyed[bool, startRequest](func(ctx context.Context, request startRequest) (bool, error) {
		var plan LaunchPlan
		if err := json.Unmarshal([]byte(request.Plan), &plan); err != nil {
			return false, err
		}
		if err := executeLaunch(ctx, plan); err != nil {
			return false, err
		}
		return true, nil
	})
	if cfg.Scan && strings.TrimSpace(m.ip) != "" {
		m.startScan(true)
	}
	return m
}

func (m *model) page(c ui.C) {
	if m.activeScan > 0 {
		ui.Effect(c, m.scanEffect, m.scanReq)
	}
	if m.activeStart > 0 {
		ui.Effect(c, m.startEffect, m.startReq)
	}
	m.takeEffects()
	if ui.Shortcut(c, 0, "enter") {
		m.start()
	}
	theme := c.Theme()
	ui.Column(c).Fill().Pad(14).Gap(8).Children(func() {
		m.pathRow(c)
		ui.Divider(c)
		ui.Row(c).FillWidth().Gap(10).Items(ui.Stretch).Children(func() {
			m.connectionCard(c)
			m.scanCard(c, theme)
			m.pairCard(c)
		})
		m.optionsCard(c, theme)
		m.statusRow(c)
		ui.Spacer(c)
		m.footer(c)
		if m.alert != "" {
			ui.ModalWith(c, &m.showAlert, "提示", ui.ModalSpec{Width: 420, NoClose: m.quitOnAck}, func() {
				ui.Wrap(c, m.alert).FillWidth()
			}, func() {
				if ui.Button(c, "知道了").Clicked() {
					m.showAlert = false
					if m.quitOnAck {
						m.quitOnAck = false
						go m.quit()
					}
				}
			})
		}
	})
}

func (m *model) pathRow(c ui.C) {
	ui.Row(c).FillWidth().Gap(8).Items(ui.Center).Children(func() {
		ui.Text(c, "scrcpy 路径:").FontSize(13)
		ui.TextField(c, &m.path, "scrcpy 路径", "可拖入目录或文件，与 scrcpy 同目录可留空").Grow(1).OnKey(m.enterStarts)
		if ui.Button(c, "选择目录").Disabled(m.busy).Clicked() {
			directory, err := app.OpenDirectory(app.OpenDirectoryOptions{Title: "选择 scrcpy 目录"})
			if err == nil && directory != "" {
				m.path = directory
			}
		}
	})
}

func (m *model) connectionCard(c ui.C) {
	ui.Column(c).Grow(1).Basis(0).Pad(12).Gap(6).Background(c.Theme().Card).Border(1, c.Theme().Border).Radius(10).Children(func() {
		ui.Text(c, "连接").FontSize(14).Bold()
		ui.Spacer(c)
		m.stackedField(c, "IP", "", &m.ip, m.busy)
		ui.Spacer(c)
		m.stackedField(c, "端口", "", &m.port, m.busy)
		ui.Spacer(c)
	})
}

func (m *model) scanCard(c ui.C, theme *gfx.Theme) {
	disabled := !m.scanEnabled || m.scanning || m.busy
	ui.Column(c).Grow(1).Basis(0).Pad(12).Gap(6).Background(theme.Card).Border(1, theme.Border).Radius(10).Children(func() {
		ui.Row(c).FillWidth().Gap(8).Items(ui.Center).Children(func() {
			ui.Checkbox(c, &m.scanEnabled, "扫描").Bold().Disabled(m.busy || m.scanning)
			ui.Spacer(c)
			if ui.Button(c, "扫描端口").IconOnly().Icon(scanIcon).Disabled(disabled).Clicked() {
				m.startScan(false)
			}
		})
		m.stackedField(c, "端口范围", "36000-44000", &m.portRange, disabled)
		m.stackedField(c, "屏蔽端口", "多个用逗号或空格分隔", &m.exclude, disabled)
		ui.Row(c).FillWidth().Gap(8).Children(func() {
			m.growingField(c, "线程", "400", &m.threads, disabled)
			m.growingField(c, "超时", "1.5", &m.timeout, disabled)
		})
	})
}

func (m *model) pairCard(c ui.C) {
	disabled := !m.pairEnabled || m.busy
	ui.Column(c).Grow(1).Basis(0).Pad(12).Gap(6).Background(c.Theme().Card).Border(1, c.Theme().Border).Radius(10).Children(func() {
		ui.Checkbox(c, &m.pairEnabled, "配对").Bold().Disabled(m.busy)
		ui.Spacer(c)
		m.stackedField(c, "配对端口", "例如 37123", &m.pairPort, disabled)
		ui.Spacer(c)
		m.stackedField(c, "配对码", "手机显示的 6 位配对码", &m.pairCode, disabled)
		ui.Spacer(c)
	})
}

func (m *model) optionsCard(c ui.C, theme *gfx.Theme) {
	ui.Column(c).FillWidth().Pad(12).Gap(8).Background(theme.Card).Border(1, theme.Border).Radius(10).Children(func() {
		ui.Text(c, "选项").FontSize(14).Bold()
		ui.Row(c).FillWidth().Gap(12).Children(func() {
			ui.Checkbox(c, &m.alwaysOnTop, "置顶").Grow(1).Basis(0)
			ui.Checkbox(c, &m.stayAwake, "熄屏且禁止休眠").Grow(1).Basis(0)
			ui.Checkbox(c, &m.powerOff, "退出后锁屏").Grow(1).Basis(0)
			record := ui.Checkbox(c, &m.record, "录屏").Grow(1).Basis(0)
			if m.record {
				record.Color(theme.Danger)
			}
		})
		ui.Row(c).FillWidth().Gap(12).Children(func() {
			ui.Checkbox(c, &m.hevc, "H.265").Grow(1).Basis(0)
			ui.Checkbox(c, &m.noAudio, "禁用音频转发").Grow(1).Basis(0)
			ui.Checkbox(c, &m.noMouseHover, "禁用鼠标悬停").Grow(1).Basis(0)
			ui.Checkbox(c, &m.recordAudio, "仅录音").Grow(1).Basis(0)
		})
		ui.Row(c).FillWidth().Gap(12).Children(func() {
			m.growingField(c, "热键", "lalt", &m.key, m.busy)
			m.growingField(c, "高度", "", &m.size, m.busy)
			m.growingField(c, "码率(Mbps)", "8", &m.bitrate, m.busy)
			m.growingField(c, "缓冲(毫秒)", "", &m.buffer, m.busy)
		})
	})
}

func (m *model) statusRow(c ui.C) {
	ui.Row(c).FillWidth().Height(24).Gap(8).Items(ui.Center).Children(func() {
		if m.scanning || m.busy {
			ui.Spinner(c)
		}
		text := m.status
		if text == "" {
			text = " "
		}
		ui.Text(c, text).FontSize(12).Muted().Ellipsis().Grow(1)
	})
}

func (m *model) footer(c ui.C) {
	ui.Row(c).FillWidth().Gap(10).Items(ui.Center).Children(func() {
		if ui.Button(c, "scrcpy v3.2").Link().Clicked() {
			_ = app.OpenURL(scrcpyURL)
		}
		ui.Spacer(c)
		if ui.Button(c, "启 动").Width(180).Disabled(m.busy || m.scanning).Clicked() {
			m.start()
		}
		ui.Spacer(c)
		ui.Text(c, "by Courser").FontSize(12).Muted()
	})
}

func (m *model) stackedField(c ui.C, label, placeholder string, value *string, disabled bool) {
	ui.Column(c).FillWidth().Gap(3).Children(func() {
		ui.Text(c, label).FontSize(12).Muted()
		ui.TextField(c, value, label, placeholder).FillWidth().OnKey(m.enterStarts).Disabled(disabled)
	})
}

func (m *model) growingField(c ui.C, label, placeholder string, value *string, disabled bool) {
	ui.Column(c).Grow(1).Basis(0).Gap(3).Children(func() {
		ui.Text(c, label).FontSize(12).Muted()
		ui.TextField(c, value, label, placeholder).FillWidth().OnKey(m.enterStarts).Disabled(disabled)
	})
}

// enterStarts 把输入框里的回车接到“启动”上：页面级快捷键不接管正在输入的字段，
// 所以字段自己处理这一次按键。
func (m *model) enterStarts(event gfx.KeyEvent) gfx.EventResult {
	if event.Key == gfx.KeyEnter && event.Modifiers == 0 {
		m.start()
		return gfx.EventResult{Handled: true}
	}
	return gfx.EventResult{}
}

func (m *model) takeEffects() {
	if m.scanEffect != nil {
		if result, err, ok := m.scanEffect.Take(); ok {
			m.scanning = false
			m.activeScan = 0
			if err != nil {
				m.alert = "扫描失败: " + err.Error()
				m.showAlert = true
				m.status = "扫描失败"
				m.setTitle("扫描失败")
			} else {
				m.applyScan(result)
			}
		}
	}
	if m.startEffect != nil {
		if _, err, ok := m.startEffect.Take(); ok {
			m.busy = false
			m.activeStart = 0
			if err != nil {
				m.alert = "启动失败: " + err.Error()
				m.showAlert = true
				m.status = "启动失败"
				m.setTitle("启动失败")
			} else {
				m.status = "scrcpy 已启动"
				m.setTitle("已启动")
				if m.application != nil {
					go m.application.Quit()
				}
			}
		}
	}
}

func (m *model) applyScan(result ScanResult) {
	elapsed := result.Elapsed.Seconds()
	if result.Single {
		if len(result.Ports) == 1 {
			m.port = strconv.Itoa(result.Ports[0])
		}
		m.status = fmt.Sprintf("端口正常，用时 %.3fs", elapsed)
		m.setTitle("端口正常")
		return
	}
	switch len(result.Ports) {
	case 0:
		m.status = fmt.Sprintf("未发现开放端口，用时 %.3fs", elapsed)
		m.setTitle("扫描结果 0 个端口")
	case 1:
		m.port = strconv.Itoa(result.Ports[0])
		m.status = fmt.Sprintf("扫描成功，端口 %d，用时 %.3fs", result.Ports[0], elapsed)
		m.setTitle("扫描成功")
	default:
		list := make([]string, len(result.Ports))
		for index, port := range result.Ports {
			list[index] = strconv.Itoa(port)
		}
		csv := strings.Join(list, ",")
		if err := app.SetClipboardText(csv); err == nil {
			m.status = fmt.Sprintf("扫描到 %d 个端口，已复制到剪贴板: %s", len(result.Ports), csv)
		} else {
			m.status = fmt.Sprintf("扫描到 %d 个端口: %s", len(result.Ports), csv)
		}
		m.setTitle(fmt.Sprintf("扫描结果 %d 个端口", len(result.Ports)))
	}
}

func (m *model) startScan(startup bool) {
	if m.scanning || m.busy {
		return
	}
	threads, err := strconv.Atoi(strings.TrimSpace(m.threads))
	if err != nil || threads <= 0 {
		threads = defaultThreads
	}
	timeout, err := strconv.ParseFloat(strings.TrimSpace(m.timeout), 64)
	if err != nil || timeout <= 0 {
		timeout = defaultTimeout
	}
	m.nextScan++
	m.scanReq = ScanRequest{
		Seq:       m.nextScan,
		IP:        strings.TrimSpace(m.ip),
		Port:      strings.TrimSpace(m.port),
		PortRange: m.portRange,
		Exclude:   m.exclude,
		Threads:   threads,
		Timeout:   timeout,
		Startup:   startup,
	}
	m.activeScan = m.nextScan
	m.scanning = true
	m.status = "正在扫描端口..."
	m.setTitle("扫描中")
}

func (m *model) start() {
	if m.busy || m.scanning {
		return
	}
	plan, err := buildLaunchPlan(m.launchInput())
	if err != nil {
		m.alert = err.Error()
		m.showAlert = true
		m.status = "请检查设置"
		return
	}
	if err := SaveConfig(m.cfgPath, m.config()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	data, err := json.Marshal(plan)
	if err != nil {
		m.alert = "生成启动参数失败: " + err.Error()
		m.showAlert = true
		return
	}
	m.nextStart++
	m.startReq = startRequest{Seq: m.nextStart, Plan: string(data)}
	m.activeStart = m.nextStart
	m.busy = true
	m.status = "正在连接设备并启动 scrcpy..."
	m.setTitle("启动中")
}

func (m *model) launchInput() launchInput {
	return launchInput{
		Path:         strings.TrimSpace(m.path),
		IP:           m.ip,
		Port:         m.port,
		Scan:         m.scanEnabled,
		PortRange:    m.portRange,
		Pair:         m.pairEnabled,
		PairPort:     m.pairPort,
		PairCode:     m.pairCode,
		AlwaysOnTop:  m.alwaysOnTop,
		StayAwake:    m.stayAwake,
		PowerOff:     m.powerOff,
		HEVC:         m.hevc,
		NoAudio:      m.noAudio,
		NoMouseHover: m.noMouseHover,
		Record:       m.record,
		RecordAudio:  m.recordAudio,
		Key:          m.key,
		Size:         m.size,
		Bitrate:      m.bitrate,
		Buffer:       m.buffer,
		RecordDir:    "rec",
	}
}

func (m *model) config() Config {
	threads, err := strconv.Atoi(strings.TrimSpace(m.threads))
	if err != nil || threads <= 0 {
		threads = defaultThreads
	}
	timeout, err := strconv.ParseFloat(strings.TrimSpace(m.timeout), 64)
	if err != nil || timeout <= 0 {
		timeout = defaultTimeout
	}
	return Config{
		Path:         strings.TrimSpace(m.path),
		IP:           strings.TrimSpace(m.ip),
		Port:         strings.TrimSpace(m.port),
		Scan:         m.scanEnabled,
		PortRange:    strings.TrimSpace(m.portRange),
		Exclude:      strings.TrimSpace(m.exclude),
		Threads:      threads,
		Timeout:      timeout,
		AlwaysOnTop:  m.alwaysOnTop,
		StayAwake:    m.stayAwake,
		PowerOff:     m.powerOff,
		HEVC:         m.hevc,
		NoAudio:      m.noAudio,
		NoMouseHover: m.noMouseHover,
		Record:       m.record,
		RecordAudio:  m.recordAudio,
		Key:          strings.TrimSpace(m.key),
		Size:         strings.TrimSpace(m.size),
		Bitrate:      strings.TrimSpace(m.bitrate),
		Buffer:       strings.TrimSpace(m.buffer),
	}
}

func (m *model) setTitle(suffix string) {
	title := appVersion
	if suffix != "" {
		title += " - " + suffix
	}
	if m.window != nil {
		_ = m.window.SetTitle(title)
	}
}

func (m *model) quit() {
	if m.application != nil {
		m.application.Quit()
		return
	}
	m.status = "已退出"
}

func (m *model) handleFileDrop(event gfx.FileDropEvent) gfx.EventResult {
	switch event.Phase {
	case gfx.FileDropOver:
		return gfx.EventResult{Handled: true}
	case gfx.FileDropDrop:
		if len(event.Paths) == 0 {
			return gfx.EventResult{}
		}
		m.path = pathFromDrop(event.Paths[0])
		m.status = "已载入路径: " + m.path
		return gfx.EventResult{Handled: true, Repaint: true}
	}
	return gfx.EventResult{}
}

func pathFromDrop(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return path
	}
	if info.IsDir() {
		return filepath.Clean(path)
	}
	return filepath.Dir(path)
}
