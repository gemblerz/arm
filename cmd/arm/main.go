package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gemblerz/arm/internal/hardware"
	"github.com/gemblerz/arm/pkg/robot"
)

func main() {
	// Command line flags
	boardType := flag.String("board", "mock", "Board type: mock, raspberry-pi, coral, linux-gpio")
	verbose := flag.Bool("verbose", true, "Enable verbose logging")
	configMode := flag.String("config", "auto", "Configuration mode: auto, raspberry-pi, coral, generic")
	flag.Parse()

	fmt.Println("🤖 6-DOF Robot Arm Controller")
	fmt.Println("==============================")
	fmt.Printf("Board Type: %s\n", *boardType)
	fmt.Printf("Config Mode: %s\n", *configMode)
	fmt.Printf("Verbose Mode: %t\n\n", *verbose)

	// Initialize hardware board
	factory := &hardware.BoardFactory{}
	board, err := factory.CreateBoard(*boardType, *verbose)
	if err != nil {
		log.Fatalf("Failed to create board: %v", err)
	}

	if err := board.Initialize(); err != nil {
		log.Fatalf("Failed to initialize board: %v", err)
	}
	defer board.Cleanup()

	fmt.Printf("Board initialized: %s\n\n", board.GetName())

	// Get board configuration
	var boardConfig hardware.BoardConfig
	switch *configMode {
	case "raspberry-pi":
		boardConfig = hardware.GetRaspberryPiConfig()
	case "coral":
		boardConfig = hardware.GetCoralDevBoardConfig()
	case "coral-mixed":
		boardConfig = hardware.GetCoralMixedDriverConfig()
	case "generic":
		boardConfig = hardware.GetGenericLinuxConfig()
	default: // auto
		switch *boardType {
		case "raspberry-pi", "rpi":
			boardConfig = hardware.GetRaspberryPiConfig()
		case "coral":
			boardConfig = hardware.GetCoralDevBoardConfig()
		case "coral-extended":
			boardConfig = hardware.GetCoralMixedDriverConfig()
		default:
			boardConfig = hardware.GetGenericLinuxConfig()
		}
	}

	fmt.Printf("Using configuration: %s\n", boardConfig.Name)
	fmt.Printf("GPIO Settings: %s numbering, %s logic\n\n", 
		boardConfig.GPIO.Numbering, boardConfig.GPIO.Voltage)

	// Create joints from configuration
	var joints []robot.Joint
	for _, motorConfig := range boardConfig.Motors {
		motor, err := board.CreateStepperMotor(motorConfig.Config)
		if err != nil {
			log.Fatalf("Failed to create motor for joint %s: %v", motorConfig.JointID, err)
		}

		joints = append(joints, robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motor,
			MinPos: motorConfig.MinPos,
			MaxPos: motorConfig.MaxPos,
		})
		
		if *verbose {
			fmt.Printf("✓ Created joint '%s' on GPIO pins: Step=%d, Dir=%d, Enable=%d\n",
				motorConfig.JointID, motorConfig.Config.StepPin, 
				motorConfig.Config.DirPin, motorConfig.Config.EnablePin)
		}
	}

	fmt.Printf("\nCreated %d joints for robot arm\n\n", len(joints))

	// Create robot arm
	arm := robot.NewArm(joints)

	// Enable all motors
	fmt.Println("Enabling all motors...")
	if err := arm.EnableAll(); err != nil {
		log.Fatalf("Failed to enable motors: %v", err)
	}

	// Home the arm
	if err := arm.Home(); err != nil {
		log.Fatalf("Failed to home arm: %v", err)
	}

	// Display status
	fmt.Println("\n" + arm.GetStatus())

	// Demo sequence
	fmt.Println("\n🎬 Starting demo sequence...")
	
	// Sequence 1: Basic movements
	fmt.Println("\n📍 Sequence 1: Basic joint movements")
	movements1 := []map[string]int{
		{"pan": 100, "tilt": 50},
		{"joint3": -100, "joint4": 75},
		{"joint5": 50, "gripper": 25},
	}
	delays1 := []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second}

	if err := arm.ExecuteSequence(movements1, delays1); err != nil {
		log.Printf("Error in sequence 1: %v", err)
	}

	// Sequence 2: Pick and place simulation
	fmt.Println("\n📦 Sequence 2: Pick and place simulation")
	movements2 := []map[string]int{
		{"pan": 200, "tilt": 100, "joint3": -150},     // Move to pick position
		{"gripper": 80},                               // Close gripper
		{"tilt": -50, "joint3": 100},                  // Lift object
		{"pan": -200},                                 // Rotate to place position
		{"tilt": 100, "joint3": -150},                 // Lower to place position
		{"gripper": 0},                                // Open gripper
		{"pan": 0, "tilt": 0, "joint3": 0},           // Return to home
	}
	delays2 := []time.Duration{
		3 * time.Second, 1 * time.Second, 2 * time.Second,
		2 * time.Second, 3 * time.Second, 1 * time.Second, 3 * time.Second,
	}

	if err := arm.ExecuteSequence(movements2, delays2); err != nil {
		log.Printf("Error in sequence 2: %v", err)
	}

	// Final status
	fmt.Println("\n📊 Final Status:")
	fmt.Println(arm.GetStatus())

	// Demonstrate individual joint control
	fmt.Println("\n🎯 Individual joint control demo:")
	
	// Move pan joint (base rotation)
	fmt.Println("Moving pan joint...")
	if err := arm.MoveJoint("pan", 300); err != nil {
		log.Printf("Error moving pan: %v", err)
	}
	
	time.Sleep(1 * time.Second)
	
	// Move tilt joint (elevation)
	fmt.Println("Moving tilt joint...")
	if err := arm.MoveJoint("tilt", -200); err != nil {
		log.Printf("Error moving tilt: %v", err)
	}

	// Return to home position
	fmt.Println("\n🏠 Returning to home position...")
	homePositions := map[string]int{
		"pan": 0, "tilt": 0, "joint3": 0,
		"joint4": 0, "joint5": 0, "gripper": 0,
	}
	
	if err := arm.MoveJoints(homePositions); err != nil {
		log.Printf("Error returning home: %v", err)
	}

	// Disable all motors
	fmt.Println("\n🔌 Disabling all motors...")
	if err := arm.DisableAll(); err != nil {
		log.Printf("Error disabling motors: %v", err)
	}

	fmt.Println("\n✅ Demo completed successfully!")
	fmt.Println("📝 Check the logs above to see the mock hardware interactions.")
}
