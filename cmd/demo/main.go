package main

import (
	"fmt"
	"log"
	"time"

	"demo-app/internal/analytics"
	"demo-app/internal/export"
	"demo-app/internal/reporting"

	"github.com/yourorg/lcc-sdk/pkg/client"
	"github.com/yourorg/lcc-sdk/pkg/config"
)

var lccClient *client.Client

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
	
	// Check if feature is enabled
	status, err := lccClient.CheckFeature("advanced_analytics")
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
	analytics.RunAdvanced()
	
	// Report usage
	if err := lccClient.ReportUsage("advanced_analytics", 1); err != nil {
		log.Printf("Warning: Failed to report usage: %v", err)
	}

	fmt.Println("✓ Advanced analytics completed")
}

func exportToPDF() {
	fmt.Println("\n[PDF Export]")
	
	// Check if feature is enabled
	status, err := lccClient.CheckFeature("pdf_export")
	if err != nil {
		fmt.Printf("✗ Failed to check feature: %v\n", err)
		return
	}

	if !status.Enabled {
		fmt.Printf("✗ Feature not available: %s\n", status.Reason)
		if status.Reason == "insufficient_tier" {
			fmt.Println("  Please upgrade to Professional tier")
		} else if status.Reason == "quota_exceeded" {
			fmt.Println("  Daily quota exceeded, please try tomorrow")
		}
		return
	}

	// Feature is enabled, use it
	export.GeneratePDF("report.pdf")
	
	// Report usage
	if err := lccClient.ReportUsage("pdf_export", 1); err != nil {
		log.Printf("Warning: Failed to report usage: %v", err)
	}

	fmt.Println("✓ PDF exported successfully")
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
		fmt.Println()
	}
	fmt.Println("-----------------------------------")
}
