package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gemblerz/arm/internal/hardware"
	"github.com/gemblerz/arm/pkg/display"
	"github.com/gemblerz/arm/pkg/robot"
	"github.com/gemblerz/arm/pkg/stepper"
)

// Interactive Robot Display Demo with Menu Navigation and Notifications
func main() {
	fmt.Println("🎮 Interactive 6-DOF Robot Arm Display Demo")
	fmt.Println("============================================")

	// Create SSD1306 display
	displayConfig := display.DisplayConfig{
		Width:      128,
		Height:     64,
		I2CAddress: 0x3C,
		I2CBus:     1,
		Mock:       true,  // Set to false for real hardware
		Verbose:    true,
	}

	ssd1306, err := display.NewSSD1306(displayConfig)
	if err != nil {
		log.Fatalf("Failed to create SSD1306 display: %v", err)
	}
	defer ssd1306.Close()

	// Create interactive display
	interactiveDisplay := display.NewInteractiveDisplay(ssd1306)
	defer interactiveDisplay.Close()

	// Create notification manager
	notificationManager := display.NewNotificationManager(ssd1306)

	// Create performance monitor
	performanceMonitor := display.NewPerformanceMonitor(ssd1306)

	// Create robot arm
	config := hardware.GetCoralMixedDriverConfig()
	joints := make([]robot.Joint, len(config.Motors))
	
	for i, motorConfig := range config.Motors {
		motor := stepper.NewMockMotor(motorConfig.Config, false)
		
		joints[i] = robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motor,
			MinPos: motorConfig.MinPos,
			MaxPos: motorConfig.MaxPos,
		}
	}
	
	arm := robot.NewArm(joints)

	// Demo the interactive display system
	fmt.Println("\n=== Interactive Display Demo ===")
	
	// Show welcome notification
	notificationManager.ShowInfo("System Starting...")
	notificationManager.Update()
	time.Sleep(2 * time.Second)

	// Show startup screen with interactive menu
	fmt.Println("Showing interactive menu system...")
	
	// Simulate user navigation through the menu
	demoInteractiveNavigation(interactiveDisplay, notificationManager, performanceMonitor, arm)
}

// demoInteractiveNavigation demonstrates the interactive display navigation
func demoInteractiveNavigation(interactive *display.InteractiveDisplay, notifications *display.NotificationManager, performance *display.PerformanceMonitor, arm *robot.Arm) {
	
	// Demo 1: Navigate through main menu
	fmt.Println("Demo 1: Main Menu Navigation")
	time.Sleep(2 * time.Second)
	
	// Navigate down through menu items
	for i := 0; i < 3; i++ {
		interactive.NavigateDown()
		time.Sleep(1 * time.Second)
	}
	
	// Select "Performance" submenu
	fmt.Println("Selecting Performance submenu...")
	interactive.Select()
	time.Sleep(2 * time.Second)
	
	// Navigate in submenu
	interactive.NavigateDown()
	time.Sleep(1 * time.Second)
	
	// Select "CPU Usage"
	fmt.Println("Showing CPU Usage chart...")
	interactive.Select()
	time.Sleep(3 * time.Second)
	
	// Go back to submenu
	interactive.NavigateDown() // Navigate to "Back"
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.Select() // Select "Back"
	time.Sleep(1 * time.Second)
	
	// Demo 2: Show notification system
	fmt.Println("\nDemo 2: Notification System")
	
	// Show different types of notifications
	notifications.ShowInfo("Robot initialized successfully")
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	notifications.ShowWarning("Motor temperature high")
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	notifications.ShowSuccess("Calibration complete")
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	// Show progress notification
	for progress := 0.0; progress <= 1.0; progress += 0.1 {
		notifications.ShowProgress("Homing robot arm...", progress)
		notifications.Update()
		time.Sleep(500 * time.Millisecond)
	}
	
	notifications.ShowSuccess("Homing completed")
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	// Demo 3: Performance monitoring
	fmt.Println("\nDemo 3: Performance Monitoring")
	
	// Update performance metrics and show charts
	for i := 0; i < 10; i++ {
		metrics := performance.SimulateMetrics()
		performance.UpdateMetrics(metrics)
		
		// Cycle through different performance views
		switch i % 4 {
		case 0:
			performance.DisplaySystemOverview()
		case 1:
			performance.DisplayCPUChart()
		case 2:
			performance.DisplayMemoryChart()
		case 3:
			performance.DisplayMotorLoadChart()
		}
		
		time.Sleep(1 * time.Second)
	}
	
	// Demo 4: Interactive workspace visualization
	fmt.Println("\nDemo 4: Workspace Visualization")
	
	// Navigate to workspace map
	interactive.NavigateUp() // Go to "Workspace Map"
	interactive.NavigateUp()
	interactive.Select()
	time.Sleep(5 * time.Second)
	
	// Demo 5: System health monitoring
	fmt.Println("\nDemo 5: System Health")
	
	// Navigate to system health
	interactive.NavigateDown() // Navigate to "System Health"
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.Select()
	time.Sleep(5 * time.Second)
	
	// Demo 6: Emergency scenarios
	fmt.Println("\nDemo 6: Emergency Notifications")
	
	// Show critical error
	notifications.ShowCriticalError("EMERGENCY STOP ACTIVATED")
	for i := 0; i < 6; i++ {
		notifications.Update()
		time.Sleep(500 * time.Millisecond)
	}
	
	// Clear critical notification
	notifications.DismissAll()
	
	// Show recovery
	notifications.ShowSuccess("System recovered")
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	// Demo 7: Settings and configuration
	fmt.Println("\nDemo 7: Settings Menu")
	
	// Navigate to settings
	interactive.NavigateDown() // Navigate to "Settings"
	interactive.NavigateDown()
	interactive.Select()
	time.Sleep(2 * time.Second)
	
	// Toggle some settings
	interactive.Select() // Toggle "Auto Screen"
	time.Sleep(1 * time.Second)
	
	interactive.NavigateDown()
	interactive.Select() // Show "Update Rate"
	time.Sleep(2 * time.Second)
	
	// Go back to main menu
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.NavigateDown()
	interactive.Select() // Select "Back"
	time.Sleep(1 * time.Second)
	
	// Demo 8: Robot control integration
	fmt.Println("\nDemo 8: Robot Control Integration")
	
	// Enable robot
	if err := arm.EnableAll(); err != nil {
		notifications.ShowError(fmt.Sprintf("Failed to enable motors: %v", err))
	} else {
		notifications.ShowSuccess("All motors enabled")
	}
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	// Simulate robot operations with progress feedback
	notifications.ShowProgress("Initializing robot...", 0.0)
	notifications.Update()
	time.Sleep(1 * time.Second)
	
	notifications.ShowProgress("Calibrating sensors...", 0.3)
	notifications.Update()
	time.Sleep(1 * time.Second)
	
	notifications.ShowProgress("Checking limits...", 0.6)
	notifications.Update()
	time.Sleep(1 * time.Second)
	
	notifications.ShowProgress("Homing complete", 1.0)
	notifications.Update()
	time.Sleep(2 * time.Second)
	
	// Show final status
	notifications.ShowSuccess("Robot ready for operation")
	notifications.Update()
	time.Sleep(3 * time.Second)
	
	// Final performance overview
	fmt.Println("\nDemo Complete: Showing final system status")
	metrics := performance.SimulateMetrics()
	performance.UpdateMetrics(metrics)
	performance.DisplaySystemOverview()
	time.Sleep(5 * time.Second)
	
	fmt.Println("\n🎉 Interactive Display Demo Complete!")
	fmt.Println("========================================")
	fmt.Println("Features demonstrated:")
	fmt.Println("✓ Interactive menu navigation")
	fmt.Println("✓ Multi-level menu system")
	fmt.Println("✓ Real-time notifications")
	fmt.Println("✓ Progress indicators")
	fmt.Println("✓ Performance monitoring")
	fmt.Println("✓ System health displays")
	fmt.Println("✓ Workspace visualization")
	fmt.Println("✓ Emergency alerts")
	fmt.Println("✓ Settings configuration")
	fmt.Println("✓ Robot control integration")
}
