package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStatusServerJSON(t *testing.T) {
	// Use a test-specific port to avoid conflicts
	os.Setenv("LCC_DEMO_STATUS_ADDR", ":18080")

	// Seed some stats
	statsMu.Lock()
	stats = DemoStats{
		AdvancedCalls:  3,
		PDFExports:     2,
		Projects:       5,
		LastTPS:        7.5,
		ConcurrentJobs: 1,
	}
	statsMu.Unlock()

	// Start server
	go startStatusServer()

	// Wait for server to be ready (simple retry loop)
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		resp, err = http.Get("http://127.0.0.1:18080/status/json")
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("failed to reach status server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got DemoStats
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if got.AdvancedCalls != 3 || got.PDFExports != 2 || got.Projects != 5 || got.LastTPS != 7.5 || got.ConcurrentJobs != 1 {
		t.Fatalf("unexpected stats: %+v", got)
	}
}

func TestStatusServerHTML(t *testing.T) {
	// Use another test-specific port
	os.Setenv("LCC_DEMO_STATUS_ADDR", ":18081")

	statsMu.Lock()
	stats = DemoStats{
		AdvancedCalls:  1,
		PDFExports:     0,
		Projects:       2,
		LastTPS:        3.5,
		ConcurrentJobs: 0,
	}
	statsMu.Unlock()

	go startStatusServer()

	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		resp, err = http.Get("http://127.0.0.1:18081/status")
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("failed to reach status server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read HTML body: %v", err)
	}
	body := string(bodyBytes)

	// Simple content checks
	if !strings.Contains(body, "LCC Demo Status") {
		t.Fatalf("HTML does not contain title: %s", body)
	}
	if !strings.Contains(body, "Advanced analytics calls:") {
		t.Fatalf("HTML does not contain advanced calls label: %s", body)
	}
	if !strings.Contains(body, "1") {
		t.Fatalf("HTML does not contain expected value for AdvancedCalls: %s", body)
	}
}
