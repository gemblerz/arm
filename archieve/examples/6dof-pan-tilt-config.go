package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gemblerz/arm/internal/hardware"
	"github.com/gemblerz/arm/pkg/robot"
	"github.com/gemblerz/arm/pkg/stepper"
)

// Example: 6-DOF Pan-Tilt Robot Configuration
// This demonstrates how to set up a 6-DOF robot with:
// - Motors 1-2: Large motors (Pan/Tilt) using DM556T drivers
// - Motor 3: Medium motor (3rd joint) using DM320T driver  
// - Motors 4-6: Smaller motors using DRV8825 drivers

func main() {
	// Get the Coral mixed driver configuration
	config := hardware.GetCoralMixedDriverConfig()
	
	fmt.Printf("Robot Configuration: %s\n", config.Name)
	fmt.Printf("Board Type: %s\n", config.Type)
	fmt.Printf("GPIO Setup: %s at %s\n", config.GPIO.Numbering, config.GPIO.Voltage)
	
	// Create mock motors for demonstration
	motors := make([]stepper.StepperMotor, len(config.Motors))
	joints := make([]robot.Joint, len(config.Motors))
	
	for i, motorConfig := range config.Motors {
		// Create mock motor with logging enabled for demo
		motors[i] = stepper.NewMockMotor(motorConfig.Config, true)
		
		joints[i] = robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motors[i],
			MinPos: motorConfig.MinPos,
			MaxPos: motorConfig.MaxPos,
		}
		
		fmt.Printf("Motor %d (%s): Steps/Rev=%d, Speed=%d, Range=[%d,%d]\n",
			i+1, motorConfig.JointID, 
			motorConfig.Config.StepsPerRev,
			motorConfig.Config.DefaultSpeed,
			motorConfig.MinPos, motorConfig.MaxPos)
	}
	
	// Create the robot arm
	arm := robot.NewArm(joints)
	
	// Enable all motors
	fmt.Println("\n=== Enabling all motors ===")
	if err := arm.EnableAll(); err != nil {
		log.Fatalf("Failed to enable motors: %v", err)
	}
	
	// Home the robot
	fmt.Println("\n=== Homing robot ===")
	if err := arm.Home(); err != nil {
		log.Fatalf("Failed to home robot: %v", err)
	}
	
	// Display driver assignments
	fmt.Println("\n=== Driver Assignments ===")
	assignments := hardware.GetCoralDriverAssignments()
	for _, assignment := range assignments {
		fmt.Printf("%s (Motor): %s Driver #%d\n", 
			assignment.JointID, assignment.DriverType, assignment.DriverID)
		fmt.Printf("  Mode Pins: %d, %d, %d\n", 
			assignment.ModePin1, assignment.ModePin2, assignment.ModePin3)
		fmt.Printf("  Notes: %s\n\n", assignment.Notes)
	}
	
	// Demonstrate different movement patterns
	
	// 1. Pan-Tilt demonstration (large motors with DM556T)
	fmt.Println("=== Pan-Tilt Movement Demo ===")
	panTiltTargets := map[string]int{
		"pan":  3200,  // 90° pan
		"tilt": 1600,  // 45° tilt
	}
	
	if err := arm.MoveJoints(panTiltTargets); err != nil {
		log.Fatalf("Pan-tilt movement failed: %v", err)
	}
	time.Sleep(2 * time.Second)
	
	// 2. Full arm sequence
	fmt.Println("\n=== Full Arm Movement Sequence ===")
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
		2 * time.Second, // Move to position
		2 * time.Second, // Grab
		2 * time.Second, // Retract  
		1 * time.Second, // Release
	}
	
	if err := arm.ExecuteSequence(sequence, delays); err != nil {
		log.Fatalf("Sequence execution failed: %v", err)
	}
	
	// Display final status
	fmt.Println("\n=== Final Robot Status ===")
	fmt.Println(arm.GetStatus())
	
	// Demonstrate motor characteristics
	fmt.Println("\n=== Motor & Driver Characteristics ===")
	
	// Large motors (DM556T)
	fmt.Println("Large Motors (Pan/Tilt) - DM556T Drivers:")
	fmt.Println("  - High torque capability (up to 5.6A)")
	fmt.Println("  - 32x microstepping for smooth movement")
	fmt.Println("  - 0.5µs step pulse for high speed")
	fmt.Println("  - Digital current control via DIP switches")
	
	// Medium motor (DM320T)
	fmt.Println("\nMedium Motor (Joint3) - DM320T Driver:")
	fmt.Println("  - Intelligent features and auto-decay")
	fmt.Println("  - 32x microstepping precision")
	fmt.Println("  - 3A current capability")
	fmt.Println("  - 1µs step pulse timing")
	
	// Smaller motors (DRV8825)
	fmt.Println("\nSmaller Motors (Joint4-6) - DRV8825 Drivers:")
	fmt.Println("  - Cost-effective for lighter loads")
	fmt.Println("  - 16x microstepping (adequate precision)")
	fmt.Println("  - 2.2A maximum current")
	fmt.Println("  - 2µs step pulse timing")
	
	// Disable all motors
	fmt.Println("\n=== Disabling all motors ===")
	if err := arm.DisableAll(); err != nil {
		log.Fatalf("Failed to disable motors: %v", err)
	}
	
	fmt.Println("Demo completed successfully!")
}
