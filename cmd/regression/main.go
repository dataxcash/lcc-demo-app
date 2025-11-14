package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// DemoStats mirrors the JSON schema exposed by the demo status server.
type DemoStats struct {
	AdvancedCalls  int     `json:"advanced_calls"`
	PDFExports     int     `json:"pdf_exports"`
	Projects       int     `json:"projects"`
	LastTPS        float64 `json:"last_tps"`
	ConcurrentJobs int     `json:"concurrent_jobs"`
}

func main() {
	addr := ":19080"
	statusURL := "http://127.0.0.1" + addr + "/status/json"

	if _, err := os.Stat("./bin/demo-app"); err != nil {
		fmt.Println("demo binary ./bin/demo-app not found; please run 'make build' first")
		os.Exit(1)
	}

	cmd := exec.Command("./bin/demo-app")
	cmd.Env = append(os.Environ(), "LCC_DEMO_STATUS_ADDR="+addr)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("failed to get stdin pipe: %v\n", err)
		os.Exit(1)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting demo-app for regression run...")
	if err := cmd.Start(); err != nil {
		fmt.Printf("failed to start demo-app: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	if err := waitForStatus(statusURL, 10*time.Second); err != nil {
		fmt.Printf("status server not reachable: %v\n", err)
		fmt.Println("make sure LCC server is running and demo-app can start successfully")
		os.Exit(1)
	}

	sendChoice := func(choice int) {
		_, _ = fmt.Fprintf(stdin, "%d\n", choice)
		// Give the demo some time to execute the scenario
		time.Sleep(400 * time.Millisecond)
	}

	// Drive a few key scenarios via the menu:
	// 2: advanced analytics (consumption)
	// 3: PDF export (consumption)
	// 7: create project (capacity)
	// 8: call demo API (TPS)
	// 9: simulate concurrent jobs (concurrency)
	sendChoice(2)
	sendChoice(3)
	sendChoice(7)
	sendChoice(8)
	sendChoice(9)

	if err := waitForMetrics(statusURL, 8*time.Second); err != nil {
		fmt.Printf("Regression FAILED: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Regression PASSED: all four limitation signals observed via status JSON")
}

func waitForStatus(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("status endpoint did not become ready in time")
}

func waitForMetrics(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	seenConcurrent := false

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		var s DemoStats
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			resp.Body.Close()
			time.Sleep(200 * time.Millisecond)
			continue
		}
		resp.Body.Close()

		if s.ConcurrentJobs > 0 {
			seenConcurrent = true
		}

		if s.AdvancedCalls >= 1 && s.PDFExports >= 1 && s.Projects >= 1 && s.LastTPS > 0 && seenConcurrent {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("expected metrics not observed before timeout")
}
