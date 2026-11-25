package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
)

const (
	maxScanWorkers = 6000
	minScanTimeout = 50 * time.Millisecond
	maxScanTimeout = 10 * time.Second
)

// ScanRequest 既是交给后台扫描器的快照，也是决定何时开始扫描的键。
// Seq 让每次请求互不相同。
type ScanRequest struct {
	Seq       int
	IP        string
	Port      string
	PortRange string
	Exclude   string
	Threads   int
	Timeout   float64
	Startup   bool
}

type ScanResult struct {
	Ports   []int
	Elapsed time.Duration
	Single  bool
}

func parsePortRange(value string) (int, int, error) {
	text := strings.TrimSpace(value)
	parts := strings.SplitN(text, "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("端口范围格式应为 起始-结束，例如 36000-44000")
	}
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("端口范围起点无效: %s", parts[0])
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("端口范围终点无效: %s", parts[1])
	}
	if start < 1 || end > 65535 || start > end {
		return 0, 0, fmt.Errorf("端口范围超出 1-65535: %s", value)
	}
	return start, end, nil
}

func parseExcludedPorts(value string) map[int]bool {
	excluded := make(map[int]bool)
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || unicode.IsSpace(r)
	})
	for _, field := range fields {
		port, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil || port < 1 || port > 65535 {
			continue
		}
		excluded[port] = true
	}
	return excluded
}

func scanPort(ctx context.Context, ip string, port int, timeout time.Duration) bool {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// scanRange 用 workers 个协程探测 [startPort, endPort] 内的每个端口，
// 按端口升序返回开放的端口。
func scanRange(ctx context.Context, ip string, startPort, endPort, workers int, timeout time.Duration, excluded map[int]bool) ([]int, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("端口范围超出 1-65535: %d-%d", startPort, endPort)
	}
	if workers < 1 {
		workers = 1
	}
	if workers > maxScanWorkers {
		workers = maxScanWorkers
	}

	ports := make([]int, 0, endPort-startPort+1)
	for port := startPort; port <= endPort; port++ {
		if !excluded[port] {
			ports = append(ports, port)
		}
	}
	if len(ports) == 0 {
		return nil, nil
	}
	if workers > len(ports) {
		workers = len(ports)
	}

	var next atomic.Int64
	var mu sync.Mutex
	open := make([]int, 0, 8)
	var group sync.WaitGroup
	for i := 0; i < workers; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				index := int(next.Add(1)) - 1
				if index >= len(ports) {
					return
				}
				if scanPort(ctx, ip, ports[index], timeout) {
					mu.Lock()
					open = append(open, ports[index])
					mu.Unlock()
				}
			}
		}()
	}
	group.Wait()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sort.Ints(open)
	return open, nil
}

func effectiveTimeout(seconds float64) time.Duration {
	timeout := time.Duration(seconds * float64(time.Second))
	if timeout < minScanTimeout {
		return minScanTimeout
	}
	if timeout > maxScanTimeout {
		return maxScanTimeout
	}
	return timeout
}

// runScan 是扫描 Effect 的后台函数。Startup 为真时先检查上次使用的端口，
// 只有该端口不通时才做范围扫描。
func runScan(ctx context.Context, request ScanRequest) (ScanResult, error) {
	started := time.Now()
	ip := strings.TrimSpace(request.IP)
	if ip == "" {
		return ScanResult{}, errors.New("请先设置 IP 地址")
	}
	if net.ParseIP(ip) == nil {
		return ScanResult{}, fmt.Errorf("IP 地址无效: %s", ip)
	}
	timeout := effectiveTimeout(request.Timeout)

	if request.Startup {
		if port, err := strconv.Atoi(strings.TrimSpace(request.Port)); err == nil && port >= 1 && port <= 65535 {
			if scanPort(ctx, ip, port, timeout) {
				return ScanResult{Ports: []int{port}, Elapsed: time.Since(started), Single: true}, nil
			}
		}
	}

	start, end, err := parsePortRange(request.PortRange)
	if err != nil {
		return ScanResult{}, err
	}
	open, err := scanRange(ctx, ip, start, end, request.Threads, timeout, parseExcludedPorts(request.Exclude))
	if err != nil {
		return ScanResult{}, err
	}
	return ScanResult{Ports: open, Elapsed: time.Since(started)}, nil
}
