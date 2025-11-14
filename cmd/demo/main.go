package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"demo-app/internal/analytics"
	"demo-app/internal/export"
	"demo-app/internal/reporting"

	"github.com/yourorg/lcc-sdk/pkg/client"
	"github.com/yourorg/lcc-sdk/pkg/config"
)

var lccClient *client.Client

// Demo state for capacity/TPS/concurrency examples
var (
	projectCount int

	requestHistoryMu sync.Mutex
	requestHistory   []time.Time
)

func main() {
	fmt.Println("=== LCC Demo Application ===\n")

	// Initialize LCC SDK
	if err := initLCC(); err != nil {
		log.Fatalf("Failed to initialize LCC SDK: %v", err)
	}
	defer lccClient.Close()

	fmt.Printf("Instance ID: %s\n\n", lccClient.GetInstanceID())

	// Demo menu
	for {
		showMenu()
		var choice int
		fmt.Print("Select option: ")
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			runBasicAnalytics()
		case 2:
			runAdvancedAnalytics()
		case 3:
			exportToPDF()
		case 4:
			exportToExcel()
		case 5:
			scheduleReport()
		case 6:
			showLicenseInfo()
		case 7:
			createProjectDemo()
		case 8:
			callDemoAPIDemo()
		case 9:
			simulateConcurrentJobsDemo()
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option")
		}
		fmt.Println()
	}
}

func initLCC() error {
	cfg := &config.SDKConfig{
		LCCURL:         "https://localhost:8088",
		ProductID:      "demo-app",
		ProductVersion: "1.0.0",
		Timeout:        30 * time.Second,
		CacheTTL:       10 * time.Second,
	}

	var err error
	lccClient, err = client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Register with LCC
	if err := lccClient.Register(); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	fmt.Println("✓ Successfully registered with LCC")
	return nil
}

func showMenu() {
	fmt.Println("-----------------------------------")
	fmt.Println("1. Run Basic Analytics (Free)")
	fmt.Println("2. Run Advanced Analytics (Professional)")
	fmt.Println("3. Export to PDF (Professional, Quota: 100/day)")
	fmt.Println("4. Export to Excel (Enterprise)")
	fmt.Println("5. Schedule Report (Professional)")
	fmt.Println("6. Show License Info")
	fmt.Println("7. Create Project (State Capacity Demo)")
	fmt.Println("8. Call Demo API (TPS Demo)")
	fmt.Println("9. Simulate Concurrent Jobs (Concurrency Demo)")
	fmt.Println("0. Exit")
	fmt.Println("-----------------------------------")
}

func runBasicAnalytics() {
	fmt.Println("\n[Basic Analytics]")
	
	// This feature is always available
	analytics.RunBasic()
	fmt.Println("✓ Basic analytics completed")
}

func runAdvancedAnalytics() {
	fmt.Println("\n[Advanced Analytics]")
	
	// 消耗型控制示例：每次运行高级分析消耗 1 次额度
	allowed, remaining, reason, err := lccClient.Consume("advanced_analytics", 1, nil)
	if err != nil {
		fmt.Printf("✗ Failed to check feature: %v\n", err)
		return
	}
	if !allowed {
		fmt.Printf("✗ Feature not available: %s\n", reason)
		fmt.Println("  Please upgrade license or check quota")
		return
	}

	// Feature is enabled, use it
	analytics.RunAdvanced()

	fmt.Printf("✓ Advanced analytics completed (remaining approx: %d)\n", remaining)
}

func exportToPDF() {
	fmt.Println("\n[PDF Export]")
	
	// 消耗型控制示例：PDF 导出每天最多 N 次（由 LIC/QuotaLimits 决定）
	allowed, remaining, reason, err := lccClient.Consume("pdf_export", 1, nil)
	if err != nil {
		fmt.Printf("✗ Failed to check feature: %v\n", err)
		return
	}
	if !allowed {
		fmt.Printf("✗ Feature not available: %s\n", reason)
		if reason == "insufficient_tier" {
			fmt.Println("  Please upgrade to Professional tier")
		} else if reason == "quota_exceeded" {
			fmt.Println("  Daily quota exceeded, please try tomorrow")
		}
		return
	}

	// Feature is enabled, use it
	export.GeneratePDF("report.pdf")

	fmt.Printf("✓ PDF exported successfully (remaining approx: %d)\n", remaining)
}

func exportToExcel() {
	fmt.Println("\n[Excel Export]")
	
	// Check if feature is enabled (Enterprise only)
	status, err := lccClient.CheckFeature("excel_export")
	if err != nil {
		fmt.Printf("✗ Failed to check feature: %v\n", err)
		return
	}

	if !status.Enabled {
		fmt.Printf("✗ Feature not available: %s\n", status.Reason)
		fmt.Println("  Please upgrade to Enterprise tier")
		return
	}

	// Feature is enabled, use it
	export.GenerateExcel("report.xlsx")
	
	// Report usage
	if err := lccClient.ReportUsage("excel_export", 1); err != nil {
		log.Printf("Warning: Failed to report usage: %v", err)
	}

	fmt.Println("✓ Excel exported successfully")
}

func scheduleReport() {
	fmt.Println("\n[Schedule Report]")
	
	// Check if feature is enabled
	status, err := lccClient.CheckFeature("scheduled_reports")
	if err != nil {
		fmt.Printf("✗ Failed to check feature: %v\n", err)
		return
	}

	if !status.Enabled {
		fmt.Printf("✗ Feature not available: %s\n", status.Reason)
		fmt.Println("  Please upgrade to Professional tier")
		return
	}

	// Feature is enabled, use it
	reporting.Schedule("weekly", "Monday 9:00 AM")
	
	// Report usage
	if err := lccClient.ReportUsage("scheduled_reports", 1); err != nil {
		log.Printf("Warning: Failed to report usage: %v", err)
	}

	fmt.Println("✓ Report scheduled successfully")
}

func showLicenseInfo() {
	fmt.Println("\n[License Information]")
	
	features := []string{
		"advanced_analytics",
		"pdf_export",
		"excel_export",
		"scheduled_reports",
		"capacity.project.count",
		"api.v1.demo",
		"concurrent.user",
	}
	
	fmt.Println("Feature Status:")
	fmt.Println("-----------------------------------")
	
	for _, featureID := range features {
		status, err := lccClient.CheckFeature(featureID)
		if err != nil {
			fmt.Printf("  %s: ERROR (%v)\n", featureID, err)
			continue
		}
		
		statusSymbol := "✗"
		if status.Enabled {
			statusSymbol = "✓"
		}
		
		fmt.Printf("  %s %s", statusSymbol, featureID)
		if !status.Enabled {
			fmt.Printf(" (%s)", status.Reason)
		}
		// 展示部分控制上限信息
		if status.MaxCapacity > 0 {
			fmt.Printf(" [max_capacity=%d]", status.MaxCapacity)
		}
		if status.MaxTPS > 0 {
			fmt.Printf(" [max_tps=%.1f]", status.MaxTPS)
		}
		if status.MaxConcurrency > 0 {
			fmt.Printf(" [max_concurrency=%d]", status.MaxConcurrency)
		}
		fmt.Println()
	}
	fmt.Println("-----------------------------------")
}

// --- State capacity demo ---

func createProjectDemo() {
	fmt.Println("\n[State Capacity Demo: Project Count]")

	// 假定每次调用表示创建一个项目
	projectCount++

	allowed, max, reason, err := lccClient.CheckCapacity("capacity.project.count", projectCount)
	if err != nil {
		fmt.Printf("✗ Failed to check capacity: %v\n", err)
		projectCount--
		return
	}
	if !allowed {
		fmt.Printf("✗ Cannot create project %d: %s (max=%d)\n", projectCount, reason, max)
		projectCount--
		return
	}

	fmt.Printf("✓ Project %d created (max=%d)\n", projectCount, max)
}

// --- TPS demo ---

func callDemoAPIDemo() {
	fmt.Println("\n[TPS Demo: api.v1.demo]")

	currentTPS := recordRequestAndGetTPS()
	fmt.Printf("  Current TPS (approx): %.1f\n", currentTPS)

	allowed, maxTPS, reason, err := lccClient.CheckTPS("api.v1.demo", currentTPS)
	if err != nil {
		fmt.Printf("✗ Failed to check TPS: %v\n", err)
		return
	}
	if !allowed {
		fmt.Printf("✗ TPS limit exceeded: current=%.1f, max=%.1f (%s)\n", currentTPS, maxTPS, reason)
		return
	}

	fmt.Printf("✓ TPS within limit: current=%.1f, max=%.1f\n", currentTPS, maxTPS)
}

// 记录一次“请求”并计算最近 1 秒内的近似 TPS
func recordRequestAndGetTPS() float64 {
	now := time.Now()

	requestHistoryMu.Lock()
	defer requestHistoryMu.Unlock()

	requestHistory = append(requestHistory, now)

	// 保留最近 1 秒内的记录
	cutoff := now.Add(-1 * time.Second)
	idx := 0
	for i, ts := range requestHistory {
		if ts.After(cutoff) {
			idx = i
			break
		}
	}
	requestHistory = requestHistory[idx:]

	return float64(len(requestHistory)) / 1.0
}

// --- Concurrency demo ---

func simulateConcurrentJobsDemo() {
	fmt.Println("\n[Concurrency Demo: concurrent.user]")

	var wg sync.WaitGroup
	jobs := 15 // try to exceed MaxConcurrency=10 from demo license

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			release, allowed, reason, err := lccClient.AcquireSlot("concurrent.user", map[string]any{
				"job_id": id,
			})
			if err != nil {
				fmt.Printf("  Job %d: error acquiring slot: %v\n", id, err)
				return
			}
			if !allowed {
				fmt.Printf("  Job %d: denied (%s)\n", id, reason)
				return
			}
			defer release()

			fmt.Printf("  Job %d: running...\n", id)
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("  Job %d: done\n", id)
		}(i + 1)
	}

	wg.Wait()
	fmt.Println("  Concurrency demo finished")
}
