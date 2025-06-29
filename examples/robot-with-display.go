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

// Example: 6-DOF Robot with SSD1306 Display Integration
func main() {
	fmt.Println("🤖 6-DOF Robot Arm with SSD1306 Display")
	fmt.Println("=========================================")

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

	// Create display adapter
	displayAdapter := robot.NewDisplayAdapter(ssd1306, 100*time.Millisecond)
	defer displayAdapter.Stop()

	// Show startup screen
	displayAdapter.ShowStartupScreen()

	// Get the Coral mixed driver configuration
	config := hardware.GetCoralMixedDriverConfig()
	
	fmt.Printf("Robot Configuration: %s\n", config.Name)
	fmt.Printf("Board Type: %s\n", config.Type)
	fmt.Printf("GPIO Setup: %s at %s\n", config.GPIO.Numbering, config.GPIO.Voltage)
	
	// Create motors and joints
	joints := make([]robot.Joint, len(config.Motors))
	
	for i, motorConfig := range config.Motors {
		// Create mock motor with logging enabled for demo
		motor := stepper.NewMockMotor(motorConfig.Config, false) // Disable verbose for cleaner output
		
		joints[i] = robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motor,
			MinPos: motorConfig.MinPos,
			MaxPos: motorConfig.MaxPos,
		}
		
		fmt.Printf("Motor %d (%s): Steps/Rev=%d, Speed=%d, Range=[%d,%d]\n",
			i+1, motorConfig.JointID, 
			motorConfig.Config.StepsPerRev,
			motorConfig.Config.DefaultSpeed,
			motorConfig.MinPos, motorConfig.MaxPos)
	}
	
	// Create the robot arm with display
	arm := robot.NewArm(joints)
	arm.SetDisplay(displayAdapter)
	
	// Enable all motors
	fmt.Println("\n=== Enabling all motors ===")
	if err := arm.EnableAll(); err != nil {
		displayAdapter.ShowError(fmt.Sprintf("Enable failed: %v", err))
		log.Fatalf("Failed to enable motors: %v", err)
	}
	
	time.Sleep(2 * time.Second) // Let display show the status
	
	// Home the robot
	fmt.Println("\n=== Homing robot ===")
	if err := arm.Home(); err != nil {
		displayAdapter.ShowError(fmt.Sprintf("Homing failed: %v", err))
		log.Fatalf("Failed to home robot: %v", err)
	}
	
	time.Sleep(2 * time.Second) // Let display show homing complete
	
	// Display driver assignments
	fmt.Println("\n=== Driver Assignments ===")
	assignments := hardware.GetCoralDriverAssignments()
	for _, assignment := range assignments {
		fmt.Printf("%s (Motor): %s Driver #%d\n", 
			assignment.JointID, assignment.DriverType, assignment.DriverID)
	}
	
	// Demonstrate movements with display feedback
	fmt.Println("\n=== Movement Demo with Display ===")
	
	// 1. Pan-Tilt demonstration
	fmt.Println("Pan-Tilt Movement...")
	panTiltTargets := map[string]int{
		"pan":  3200,  // 90° pan
		"tilt": 1600,  // 45° tilt
	}
	
	if err := arm.MoveJoints(panTiltTargets); err != nil {
		displayAdapter.ShowError(fmt.Sprintf("Movement failed: %v", err))
		log.Fatalf("Pan-tilt movement failed: %v", err)
	}
	time.Sleep(3 * time.Second)
	
	// 2. Complex sequence with real-time display updates
	fmt.Println("Full Arm Movement Sequence...")
	sequence := []map[string]int{
		// Home position
		{"pan": 0, "tilt": 0, "joint3": 0, "joint4": 0, "joint5": 0, "gripper": 0},
		// Extend and reach
		{"pan": 1600, "tilt": 800, "joint3": 1600, "joint4": -800, "joint5": 0, "gripper": 0},
		// Grab position
		{"pan": 1600, "tilt": 800, "joint3": 1600, "joint4": -800, "joint5": 0, "gripper": 800},
		// Retract
		{"pan": 0, "tilt": 0, "joint3": 800, "joint4": -400, "joint5": 0, "gripper": 800},
		// Release
		{"pan": 0, "tilt": 0, "joint3": 800, "joint4": -400, "joint5": 0, "gripper": 0},
	}
	
	delays := []time.Duration{
		3 * time.Second, // Move to position
		2 * time.Second, // Grab
		3 * time.Second, // Retract  
		2 * time.Second, // Release
	}
	
	if err := arm.ExecuteSequence(sequence, delays); err != nil {
		displayAdapter.ShowError(fmt.Sprintf("Sequence failed: %v", err))
		log.Fatalf("Sequence execution failed: %v", err)
	}
	
	// Display final status
	fmt.Println("\n=== Final Robot Status ===")
	fmt.Println(arm.GetStatus())
	
	// Show status on display for a few seconds
	time.Sleep(5 * time.Second)
	
	// Demonstrate display features
	fmt.Println("\n=== Display Feature Demo ===")
	
	// Show custom display content
	ssd1306.Clear()
	ssd1306.DrawText("6-DOF Robot Complete!", 5, 15, nil)
	ssd1306.DrawText("DM556T + SSD1306", 15, 30, nil)
	ssd1306.DrawText("Ready for operation", 10, 45, nil)
	
	// Draw some graphics
	ssd1306.DrawRectangle(5, 55, 118, 8, false)
	ssd1306.DrawProgressBar(6, 56, 116, 6, 1.0)
	
	ssd1306.Update()
	time.Sleep(3 * time.Second)
	
	// Disable all motors
	fmt.Println("\n=== Disabling all motors ===")
	if err := arm.DisableAll(); err != nil {
		displayAdapter.ShowError(fmt.Sprintf("Disable failed: %v", err))
		log.Printf("Error disabling motors: %v", err)
	}
	
	// Final display message
	ssd1306.Clear()
	ssd1306.DrawText("Demo Complete!", 25, 25, nil)
	ssd1306.DrawText("System Safe", 30, 40, nil)
	ssd1306.Update()
	time.Sleep(2 * time.Second)
	
	fmt.Println("\n✅ Robot and Display Demo completed successfully!")
}
