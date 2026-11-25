package main

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestParsePortRange(t *testing.T) {
	tests := []struct {
		value string
		start int
		end   int
		ok    bool
	}{
		{"36000-44000", 36000, 44000, true},
		{" 80 - 90 ", 80, 90, true},
		{"1-65535", 1, 65535, true},
		{"0-10", 0, 0, false},
		{"10-1", 0, 0, false},
		{"70000-70001", 0, 0, false},
		{"80", 0, 0, false},
		{"a-b", 0, 0, false},
	}
	for _, test := range tests {
		start, end, err := parsePortRange(test.value)
		if test.ok {
			if err != nil || start != test.start || end != test.end {
				t.Errorf("parsePortRange(%q) = %d,%d,%v", test.value, start, end, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("parsePortRange(%q) should fail", test.value)
		}
	}
}

func TestParseExcludedPorts(t *testing.T) {
	excluded := parseExcludedPorts("80, 443 8080，9000; abc 70000")
	for _, port := range []int{80, 443, 8080, 9000} {
		if !excluded[port] {
			t.Errorf("port %d should be excluded", port)
		}
	}
	if excluded[70000] {
		t.Error("out-of-range port should be ignored")
	}
	if len(excluded) != 4 {
		t.Fatalf("excluded = %v", excluded)
	}
}

func TestScanRangeFindsOpenPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	open, err := scanRange(ctx, "127.0.0.1", port, port, 4, 200*time.Millisecond, nil)
	if err != nil {
		t.Fatalf("scanRange: %v", err)
	}
	if len(open) != 1 || open[0] != port {
		t.Fatalf("open = %v, want [%d]", open, port)
	}
}

func TestScanRangeHonorsExcludedPorts(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	open, err := scanRange(ctx, "127.0.0.1", port, port, 2, 200*time.Millisecond, map[int]bool{port: true})
	if err != nil {
		t.Fatalf("scanRange: %v", err)
	}
	if len(open) != 0 {
		t.Fatalf("open = %v, want none", open)
	}
}

func TestRunScanStartupChecksSavedPortFirst(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	result, err := runScan(context.Background(), ScanRequest{
		IP:        "127.0.0.1",
		Port:      strconv.Itoa(port),
		PortRange: "not-a-range",
		Threads:   2,
		Timeout:   0.2,
		Startup:   true,
	})
	if err != nil {
		t.Fatalf("runScan: %v", err)
	}
	if !result.Single || len(result.Ports) != 1 || result.Ports[0] != port {
		t.Fatalf("result = %+v", result)
	}
}

func TestRunScanRejectsInvalidIP(t *testing.T) {
	if _, err := runScan(context.Background(), ScanRequest{IP: "not-an-ip", PortRange: "1-2"}); err == nil {
		t.Fatal("invalid IP should fail")
	}
}
