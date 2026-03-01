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

// Advanced Robot Display Demo with Multiple Screen Types
func main() {
	fmt.Println("🤖 Advanced 6-DOF Robot Arm Display Demo")
	fmt.Println("==========================================")

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

	// Create advanced display with multiple screens
	advancedDisplay := display.NewAdvancedRobotDisplay(ssd1306, 200*time.Millisecond)
	defer advancedDisplay.Stop()

	// Enable auto-switching between screens every 4 seconds
	advancedDisplay.EnableAutoSwitch(true, 4*time.Second)

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

	// Create display adapter that works with advanced display
	displayAdapter := &AdvancedDisplayAdapter{
		advancedDisplay: advancedDisplay,
	}
	arm.SetDisplay(displayAdapter)

	// Show startup screen
	advancedDisplay.ShowStartupScreen()
	time.Sleep(3 * time.Second)

	// Start the advanced display system
	advancedDisplay.Start()

	// Demo sequence with display feedback
	fmt.Println("\n=== Advanced Display Demo Sequence ===")

	// Enable motors
	if err := arm.EnableAll(); err != nil {
		log.Fatalf("Failed to enable motors: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Home the robot (shows calibration progress)
	if err := arm.Home(); err != nil {
		log.Fatalf("Failed to home robot: %v", err)
	}
	time.Sleep(3 * time.Second)

	// Demonstrate different display screens manually
	fmt.Println("Demonstrating different display screens...")

	// Screen 1: Status Overview
	fmt.Println("Screen 1: Status Overview")
	advancedDisplay.SetScreen(display.ScreenStatus)
	time.Sleep(4 * time.Second)

	// Screen 2: Joint Positions
	fmt.Println("Screen 2: Joint Positions with Progress Bars")
	advancedDisplay.SetScreen(display.ScreenJoints)
	
	// Move some joints to show progress bars
	arm.MoveJoint("pan", 1600)
	time.Sleep(1 * time.Second)
	arm.MoveJoint("tilt", 800)
	time.Sleep(3 * time.Second)

	// Screen 3: Angle Display
	fmt.Println("Screen 3: Joint Angles")
	advancedDisplay.SetScreen(display.ScreenAngles)
	time.Sleep(4 * time.Second)

	// Screen 4: Diagnostics
	fmt.Println("Screen 4: System Diagnostics")
	advancedDisplay.SetScreen(display.ScreenDiagnostics)
	time.Sleep(4 * time.Second)

	// Screen 5: Sequence Progress
	fmt.Println("Screen 5: Sequence Progress")
	advancedDisplay.SetScreen(display.ScreenSequence)
	
	// Execute a movement sequence
	sequence := []map[string]int{
		{"pan": 0, "tilt": 0, "joint3": 0, "joint4": 0, "joint5": 0, "gripper": 0},
		{"pan": 3200, "tilt": 1600, "joint3": 1600, "joint4": -800, "joint5": 0, "gripper": 0},
		{"pan": 0, "tilt": 0, "joint3": 0, "joint4": 0, "joint5": 0, "gripper": 800},
	}
	
	delays := []time.Duration{3 * time.Second, 3 * time.Second, 2 * time.Second}
	
	if err := arm.ExecuteSequence(sequence, delays); err != nil {
		log.Printf("Sequence error: %v", err)
	}

	// Screen 6: Driver Information
	fmt.Println("Screen 6: Driver Information")
	advancedDisplay.SetScreen(display.ScreenDriverInfo)
	time.Sleep(4 * time.Second)

	// Demonstrate specialized visualizations
	fmt.Println("\n=== Specialized Visualizations ===")

	// Arm visualization
	fmt.Println("Showing 2D arm visualization...")
	armViz := display.NewArmVisualization(ssd1306)
	
	// Create a sample robot status for visualization
	sampleStatus := display.RobotStatus{
		IsHomed:   true,
		IsEnabled: true,
		IsMoving:  false,
		CurrentOp: "Visualizing",
		Joints: []display.JointStatus{
			{ID: "pan", Position: 1600, Target: 1600, MinPos: -12800, MaxPos: 12800, Enabled: true},
			{ID: "tilt", Position: 800, Target: 800, MinPos: -3200, MaxPos: 3200, Enabled: true},
			{ID: "joint3", Position: 400, Target: 400, MinPos: -3200, MaxPos: 3200, Enabled: true},
		},
	}
	
	armViz.UpdateFromRobotStatus(sampleStatus)
	armViz.DrawArmVisualization()
	time.Sleep(5 * time.Second)

	// Workspace visualization
	fmt.Println("Showing workspace visualization...")
	workspaceViz := display.NewWorkspaceVisualization(ssd1306)
	workspaceViz.DrawWorkspace()
	time.Sleep(5 * time.Second)

	// Performance monitoring
	fmt.Println("Showing performance metrics...")
	perfMonitor := display.NewPerformanceMonitor(ssd1306)
	perfMonitor.UpdateMetrics(perfMonitor.SimulateMetrics())
	perfMonitor.DisplaySystemOverview()
	time.Sleep(5 * time.Second)

	// Resume automatic screen switching
	fmt.Println("Resuming automatic screen switching...")
	advancedDisplay.EnableAutoSwitch(true, 3*time.Second)
	time.Sleep(15 * time.Second)

	// Final demonstration: Error handling
	fmt.Println("Demonstrating error display...")
	advancedDisplay.ShowError("Joint limit exceeded: Pan motor at maximum position. Check trajectory planning.")
	time.Sleep(5 * time.Second)

	// Disable motors
	if err := arm.DisableAll(); err != nil {
		log.Printf("Error disabling motors: %v", err)
	}

	// Final status screen
	ssd1306.Clear()
	ssd1306.DrawText("Demo Complete!", 25, 25, nil)
	ssd1306.DrawText("All systems", 30, 40, nil)
	ssd1306.DrawText("operational", 35, 50, nil)
	ssd1306.Update()
	time.Sleep(3 * time.Second)

	fmt.Println("\n✅ Advanced Display Demo completed successfully!")
	fmt.Println("📊 Demonstrated features:")
	fmt.Println("  - Multi-screen display system")
	fmt.Println("  - Real-time joint monitoring")
	fmt.Println("  - Angle visualization")
	fmt.Println("  - System diagnostics")
	fmt.Println("  - 2D arm visualization")
	fmt.Println("  - Workspace mapping")
	fmt.Println("  - Performance metrics")
	fmt.Println("  - Error handling")
}

// AdvancedDisplayAdapter bridges robot and advanced display interfaces
type AdvancedDisplayAdapter struct {
	advancedDisplay *display.AdvancedRobotDisplay
}

func (ada *AdvancedDisplayAdapter) UpdateStatus(status robot.DisplayStatus) {
	// Convert robot.DisplayStatus to display.RobotStatus
	robotStatus := display.RobotStatus{
		IsHomed:   status.IsHomed,
		IsEnabled: status.IsEnabled,
		IsMoving:  status.IsMoving,
		CurrentOp: status.CurrentOp,
		Timestamp: status.Timestamp,
		Joints:    make([]display.JointStatus, len(status.Joints)),
	}

	for i, joint := range status.Joints {
		robotStatus.Joints[i] = display.JointStatus{
			ID:       joint.ID,
			Position: joint.Position,
			Target:   joint.Target,
			MinPos:   joint.MinPos,
			MaxPos:   joint.MaxPos,
			Enabled:  joint.Enabled,
		}
	}

	ada.advancedDisplay.UpdateStatus(robotStatus)
}

func (ada *AdvancedDisplayAdapter) ShowStartupScreen() {
	ada.advancedDisplay.ShowStartupScreen()
}

func (ada *AdvancedDisplayAdapter) ShowError(errorMsg string) {
	ada.advancedDisplay.ShowError(errorMsg)
}

func (ada *AdvancedDisplayAdapter) ShowCalibration(jointID string, progress float64) {
	ada.advancedDisplay.ShowCalibration(jointID, progress)
}

func (ada *AdvancedDisplayAdapter) Start() {
	ada.advancedDisplay.Start()
}

func (ada *AdvancedDisplayAdapter) Stop() {
	ada.advancedDisplay.Stop()
}
