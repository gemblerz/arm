package main

// Example configuration for Google Coral Dev Board robot arm setup
// This file demonstrates how to customize the robot arm for Coral hardware

import (
	"log"
	"time"

	"github.com/gemblerz/arm/internal/hardware"
	"github.com/gemblerz/arm/pkg/robot"
)

// CoralArmConfig represents a complete Coral Dev Board robot arm setup
type CoralArmConfig struct {
	BoardType    string
	MotorDrivers string // e.g., "A4988", "DRV8825", "TMC2209"
	PowerVoltage string // e.g., "12V", "24V"
	Microstepping int   // 1, 2, 4, 8, 16, 32
}

// CreateCoralArm creates a robot arm configured specifically for Coral Dev Board
func CreateCoralArm(config CoralArmConfig) (*robot.Arm, error) {
	// Initialize Coral board
	factory := &hardware.BoardFactory{}
	board, err := factory.CreateBoard("coral", true)
	if err != nil {
		return nil, err
	}

	if err := board.Initialize(); err != nil {
		return nil, err
	}

	// Get Coral-optimized configuration
	boardConfig := hardware.GetCoralDevBoardConfig()
	
	// Adjust for microstepping
	for i := range boardConfig.Motors {
		boardConfig.Motors[i].Config.StepsPerRev *= config.Microstepping
		// Adjust speed for microstepping
		boardConfig.Motors[i].Config.MaxSpeed *= config.Microstepping
		boardConfig.Motors[i].Config.DefaultSpeed *= config.Microstepping
	}

	// Create joints
	var joints []robot.Joint
	for _, motorConfig := range boardConfig.Motors {
		motor, err := board.CreateStepperMotor(motorConfig.Config)
		if err != nil {
			return nil, err
		}

		joints = append(joints, robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motor,
			MinPos: motorConfig.MinPos * config.Microstepping,
			MaxPos: motorConfig.MaxPos * config.Microstepping,
		})
	}

	return robot.NewArm(joints), nil
}

// CoralPerformanceTest runs a performance test specifically for Coral hardware
func CoralPerformanceTest(arm *robot.Arm) error {
	log.Println("🧪 Running Coral Dev Board performance test...")

	// Test high-speed movements (Coral can handle higher frequencies)
	testSequence := []map[string]int{
		{"base": 100, "shoulder": 50},           // Fast positioning
		{"elbow": -100, "wrist_pitch": 75},      // Coordinated movement
		{"wrist_roll": 50, "gripper": 25},       // Fine control
		{"base": 0, "shoulder": 0, "elbow": 0},  // Return to home
	}

	// Shorter delays for Coral's performance
	delays := []time.Duration{
		500 * time.Millisecond,
		500 * time.Millisecond, 
		500 * time.Millisecond,
		500 * time.Millisecond,
	}

	return arm.ExecuteSequence(testSequence, delays)
}

// ExampleCoralUsage demonstrates typical usage patterns for Coral Dev Board
func ExampleCoralUsage() {
	config := CoralArmConfig{
		BoardType:     "coral",
		MotorDrivers:  "TMC2209", // Quiet, high-performance drivers
		PowerVoltage:  "24V",     // Higher voltage for better performance
		Microstepping: 16,        // Smooth, precise movement
	}

	arm, err := CreateCoralArm(config)
	if err != nil {
		log.Fatalf("Failed to create Coral arm: %v", err)
	}

	// Enable all motors
	if err := arm.EnableAll(); err != nil {
		log.Fatalf("Failed to enable motors: %v", err)
	}

	// Home the arm
	if err := arm.Home(); err != nil {
		log.Fatalf("Failed to home arm: %v", err)
	}

	// Run Coral-specific performance test
	if err := CoralPerformanceTest(arm); err != nil {
		log.Printf("Performance test failed: %v", err)
	}

	// Advanced Coral features could be added here:
	// - Edge TPU integration for path planning
	// - Camera integration for computer vision
	// - Real-time feedback control

	log.Println("✅ Coral Dev Board robot arm demo completed")
}
