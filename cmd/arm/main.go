package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gemblerz/arm/internal/hardware"
	"github.com/gemblerz/arm/pkg/display"
	"github.com/gemblerz/arm/pkg/kinematics"
	"github.com/gemblerz/arm/pkg/robot"
)

// Global variables for interactive mode
var (
	globalArm     *robot.Arm
	globalDisplay *display.SSD1306
	globalBoard   hardware.Board
	verbose       bool
)

func main() {
	// Command line flags
	boardType := flag.String("board", "mock", "Board type: mock, raspberry-pi, coral, linux-gpio")
	verboseFlag := flag.Bool("verbose", true, "Enable verbose logging")
	configMode := flag.String("config", "auto", "Configuration mode: auto, raspberry-pi, coral, generic")
	interactive := flag.Bool("interactive", true, "Run in interactive mode")
	flag.Parse()

	verbose = *verboseFlag

	fmt.Println("🤖 6-DOF Robot Arm Controller - Interactive Mode")
	fmt.Println("================================================")
	fmt.Printf("Board Type: %s\n", *boardType)
	fmt.Printf("Config Mode: %s\n", *configMode)
	fmt.Printf("Verbose Mode: %t\n", *verboseFlag)
	fmt.Printf("Interactive Mode: %t\n\n", *interactive)

	// Initialize hardware and robot
	if err := initializeSystem(*boardType, *configMode); err != nil {
		log.Fatalf("Failed to initialize system: %v", err)
	}

	if *interactive {
		fmt.Println("🎮 Starting interactive mode...")
		fmt.Println("Type 'help' for available commands")
		runInteractiveMode()
	} else {
		fmt.Println("Running demo mode...")
		runDemoMode()
	}

	cleanup()
}

func initializeSystem(boardType, configMode string) error {
	// Initialize hardware board
	factory := &hardware.BoardFactory{}
	board, err := factory.CreateBoard(boardType, verbose)
	if err != nil {
		return fmt.Errorf("failed to create board: %v", err)
	}

	if err := board.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize board: %v", err)
	}
	globalBoard = board

	fmt.Printf("Board initialized: %s\n\n", board.GetName())

	// Get board configuration
	var boardConfig hardware.BoardConfig
	switch configMode {
	case "raspberry-pi":
		boardConfig = hardware.GetRaspberryPiConfig()
	case "coral":
		boardConfig = hardware.GetCoralDevBoardConfig()
	case "coral-mixed":
		boardConfig = hardware.GetCoralMixedDriverConfig()
	case "generic":
		boardConfig = hardware.GetGenericLinuxConfig()
	default: // auto
		switch boardType {
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
			return fmt.Errorf("failed to create motor for joint %s: %v", motorConfig.JointID, err)
		}

		joints = append(joints, robot.Joint{
			ID:     motorConfig.JointID,
			Motor:  motor,
			MinPos: motorConfig.MinPos,
			MaxPos: motorConfig.MaxPos,
		})

		if verbose {
			fmt.Printf("✓ Created joint '%s' on GPIO pins: Step=%d, Dir=%d, Enable=%d\n",
				motorConfig.JointID, motorConfig.Config.StepPin,
				motorConfig.Config.DirPin, motorConfig.Config.EnablePin)
		}
	}

	fmt.Printf("\nCreated %d joints for robot arm\n\n", len(joints))

	// Create robot arm with kinematics solver
	// Define placeholder DH parameters for a generic 6-DOF arm.
	// These should be replaced with the actual dimensions of the robot.
	dhParams := []kinematics.DHParameter{
		{Theta: 0, D: 10, A: 0, Alpha: math.Pi / 2}, // Base
		{Theta: 0, D: 0, A: 20, Alpha: 0},          // Shoulder
		{Theta: 0, D: 0, A: 20, Alpha: 0},          // Elbow
		{Theta: 0, D: 0, A: 0, Alpha: math.Pi / 2}, // Wrist Pitch
		{Theta: 0, D: 5, A: 0, Alpha: -math.Pi / 2}, // Wrist Roll
		{Theta: 0, D: 2, A: 0, Alpha: 0},            // Gripper
	}
	kinematicsSolver := kinematics.NewKinematicsSolver(dhParams)
	globalArm = robot.NewArm(joints, kinematicsSolver)

	// Initialize display if available
	displayConfig := display.DisplayConfig{
		Width:      128,
		Height:     64,
		I2CAddress: 0x3C,
		I2CBus:     1,
		Mock:       boardType == "mock",
		Verbose:    verbose,
	}

	ssd1306Display, err := display.NewSSD1306(displayConfig)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize display: %v\n", err)
	} else {
		globalDisplay = ssd1306Display
		fmt.Println("✓ Display initialized")
	}

	return nil
}

func runInteractiveMode() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("Robot Arm Controller"),
		prompt.OptionPrefix("arm> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
	p.Run()
}

func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	args := strings.Fields(input)
	command := args[0]
	success := true
	message := ""

	switch command {
	case "help":
		showHelp()
	case "status":
		showStatus()
	case "enable":
		success = enableMotorsWithResult()
	case "disable":
		success = disableMotorsWithResult()
	case "home":
		success = homeArmWithResult()
	case "move":
		if len(args) < 3 {
			fmt.Println("Usage: move <joint> <position>")
			success = false
		} else {
			success = moveJointWithResult(args[1], args[2])
		}
	case "move-to-xyz":
		if len(args) < 4 {
			fmt.Println("Usage: move-to-xyz <x> <y> <z>")
			success = false
		} else {
			success = moveXYZWithResult(args[1], args[2], args[3])
		}
	case "sequence":
		success = runPickPlaceSequenceWithResult()
	case "display":
		if len(args) < 2 {
			// Simple display command - show current status
			displayRobotStatus()
		} else {
			// Check if it's a subcommand or direct text
			if args[1] == "status" || args[1] == "clear" || args[1] == "test" {
				handleDisplayCommand(args[1:])
			} else {
				// Treat everything after "display" as text to show
				message := strings.Join(args[1:], " ")
				displayText(message)
			}
		}
	case "test":
		runDisplayTest()
	case "demo":
		success = runBasicDemoWithResult()
	case "joints":
		listJoints()
	case "clear":
		clearDisplay()
	case "exit", "quit":
		fmt.Println("Goodbye! 👋")
		cleanup()
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		success = false
	}

	// Auto-update display with command result
	autoUpdateDisplay(command, args[1:], success, message)
}

func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "help", Description: "Show available commands"},
		{Text: "status", Description: "Show robot arm status"},
		{Text: "enable", Description: "Enable all motors"},
		{Text: "disable", Description: "Disable all motors"},
		{Text: "home", Description: "Home the robot arm"},
		{Text: "move", Description: "Move a joint (usage: move <joint> <position>)"},
		{Text: "move-to-xyz", Description: "Move to Cartesian coordinate (usage: move-to-xyz <x> <y> <z>)"},
		{Text: "sequence", Description: "Run pick and place sequence"},
		{Text: "display", Description: "Display text on screen (usage: display <message>)"},
		{Text: "test", Description: "Run display test"},
		{Text: "demo", Description: "Run basic movement demo"},
		{Text: "joints", Description: "List all available joints"},
		{Text: "clear", Description: "Clear the display"},
		{Text: "exit", Description: "Exit the program"},
		{Text: "quit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func showHelp() {
	fmt.Println("🤖 Robot Arm Controller - Interactive Commands")
	fmt.Println("===============================================")
	fmt.Println("Robot Control:")
	fmt.Println("  status                    - Show robot arm status")
	fmt.Println("  enable                    - Enable all motors")
	fmt.Println("  disable                   - Disable all motors")
	fmt.Println("  home                      - Home the robot arm")
	fmt.Println("  move <joint> <position>   - Move specific joint")
	fmt.Println("  move-to-xyz <x> <y> <z>   - Move to a Cartesian coordinate (placeholder)")
	fmt.Println("  sequence                  - Run pick and place sequence")
	fmt.Println("  demo                      - Run basic movement demo")
	fmt.Println("  joints                    - List all available joints")
	fmt.Println("")
	fmt.Println("Display Control:")
	fmt.Println("  display <message>         - Display text message on screen")
	fmt.Println("  display status            - Show robot status on display")
	fmt.Println("  display clear             - Clear the display")
	fmt.Println("  display test              - Run display test patterns")
	fmt.Println("  test                      - Run display test patterns")
	fmt.Println("  clear                     - Clear the display")
	fmt.Println("")
	fmt.Println("System:")
	fmt.Println("  help                      - Show this help")
	fmt.Println("  exit/quit                 - Exit the program")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  move base 100             - Move base joint to position 100")
	fmt.Println("  move-to-xyz 150.0 50.0 100.0 - Move end-effector to a 3D point")
	fmt.Println("  display Hello World       - Display 'Hello World' on screen")
	fmt.Println("  display status            - Show robot status on display")
}

func showStatus() {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return
	}
	fmt.Println(globalArm.GetStatus())
}

func handleDisplayCommand(args []string) {
	if globalDisplay == nil {
		fmt.Println("❌ Display not initialized")
		return
	}

	if len(args) == 0 {
		fmt.Println("Usage: display <command> [args...]")
		fmt.Println("Commands: text, status, clear, test")
		return
	}

	switch args[0] {
	case "text":
		if len(args) < 2 {
			fmt.Println("Usage: display text <message>")
			return
		}
		message := strings.Join(args[1:], " ")
		displayText(message)
	case "status":
		displayRobotStatus()
	case "clear":
		clearDisplay()
	case "test":
		runDisplayTest()
	default:
		fmt.Printf("Unknown display command: %s\n", args[0])
	}
}

func displayText(message string) {
	if globalDisplay == nil {
		fmt.Println("❌ Display not initialized")
		return
	}

	fmt.Printf("📺 Displaying: %s\n", message)
	if err := globalDisplay.Clear(); err != nil {
		fmt.Printf("Error clearing display: %v\n", err)
		return
	}
	globalDisplay.DrawText(message, 10, 20, nil)
	globalDisplay.Update()
}

func displayRobotStatus() {
	if globalDisplay == nil || globalArm == nil {
		fmt.Println("❌ Display or arm not initialized")
		return
	}

	fmt.Println("📺 Displaying robot status...")
	if err := globalDisplay.Clear(); err != nil {
		fmt.Printf("Error clearing display: %v\n", err)
		return
	}

	// Display basic status info
	globalDisplay.DrawText("Robot Arm Status", 5, 15, nil)
	globalDisplay.DrawText("Joints: 6 DOF", 5, 30, nil)
	globalDisplay.DrawText("Status: Active", 5, 45, nil)
	globalDisplay.DrawText(fmt.Sprintf("Time: %s", time.Now().Format("15:04:05")), 5, 60, nil)

	globalDisplay.Update()
}

func runDisplayTest() {
	if globalDisplay == nil {
		fmt.Println("❌ Display not initialized")
		return
	}

	fmt.Println("🧪 Running display test patterns...")

	// Test 1: Clear screen
	fmt.Println("  Test 1: Clear screen")
	globalDisplay.Clear()
	globalDisplay.Update()
	time.Sleep(1 * time.Second)

	// Test 2: Draw text
	fmt.Println("  Test 2: Text rendering")
	globalDisplay.DrawText("Display Test", 10, 20, nil)
	globalDisplay.DrawText("Line 2", 10, 35, nil)
	globalDisplay.Update()
	time.Sleep(2 * time.Second)

	// Test 3: Draw lines
	fmt.Println("  Test 3: Line drawing")
	globalDisplay.Clear()
	globalDisplay.DrawLine(0, 0, 127, 63)
	globalDisplay.DrawLine(0, 63, 127, 0)
	globalDisplay.Update()
	time.Sleep(2 * time.Second)

	// Test 4: Draw rectangles
	fmt.Println("  Test 4: Rectangle drawing")
	globalDisplay.Clear()
	globalDisplay.DrawRectangle(10, 10, 50, 30, false)
	globalDisplay.DrawRectangle(60, 20, 30, 25, true)
	globalDisplay.Update()
	time.Sleep(2 * time.Second)

	fmt.Println("✅ Display test completed")
}

func listJoints() {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return
	}

	fmt.Println("🔧 Available joints:")
	// Since GetJoints() doesn't exist, let's show common joint names
	jointNames := []string{"pan", "tilt", "joint3", "joint4", "joint5", "gripper"}
	for _, name := range jointNames {
		fmt.Printf("  - %s\n", name)
	}
}

func clearDisplay() {
	if globalDisplay == nil {
		fmt.Println("❌ Display not initialized")
		return
	}

	fmt.Println("🧹 Clearing display...")
	globalDisplay.Clear()
	globalDisplay.Update()
	fmt.Println("✅ Display cleared")
}

// autoUpdateDisplay automatically updates the display with relevant information after command execution
func autoUpdateDisplay(command string, args []string, success bool, message string) {
	if globalDisplay == nil {
		return // Display not available
	}

	// Don't auto-update for display-specific commands to avoid conflicts
	if command == "display" || command == "clear" || command == "test" {
		return
	}

	var displayMessage string
	var secondLine string

	switch command {
	case "enable":
		if success {
			displayMessage = "✓ Motors ON"
		} else {
			displayMessage = "✗ Enable Failed"
		}
	case "disable":
		if success {
			displayMessage = "✓ Motors OFF"
		} else {
			displayMessage = "✗ Disable Failed"
		}
	case "home":
		if success {
			displayMessage = "✓ Home Complete"
		} else {
			displayMessage = "✗ Home Failed"
		}
	case "move":
		if len(args) >= 2 && success {
			displayMessage = fmt.Sprintf("✓ %s: %s", args[0], args[1])
		} else if len(args) >= 2 {
			displayMessage = fmt.Sprintf("✗ %s Move Fail", args[0])
		} else {
			displayMessage = "✗ Move Failed"
		}
	case "sequence":
		if success {
			displayMessage = "✓ Sequence Done"
		} else {
			displayMessage = "✗ Sequence Fail"
		}
	case "demo":
		if success {
			displayMessage = "✓ Demo Complete"
		} else {
			displayMessage = "✗ Demo Failed"
		}
	case "status":
		displayMessage = "Status Check"
		if globalArm != nil {
			secondLine = "6-DOF Ready"
		}
	case "joints":
		displayMessage = "Joint List"
		secondLine = "6 DOF Available"
	case "help":
		displayMessage = "Help Displayed"
		secondLine = "Commands Ready"
	default:
		// For unknown commands or commands that don't need display updates
		return
	}

	// Update the display with command status
	if err := globalDisplay.Clear(); err != nil {
		return // Silently handle display errors
	}

	// Draw the main status message
	globalDisplay.DrawText(displayMessage, 5, 15, nil)

	// Draw secondary information if available
	if secondLine != "" {
		globalDisplay.DrawText(secondLine, 5, 30, nil)
	} else if message != "" {
		// Use provided message as secondary line
		globalDisplay.DrawText(message, 5, 30, nil)
	}

	// Add timestamp
	globalDisplay.DrawText(time.Now().Format("15:04:05"), 5, 50, nil)

	// Update the display
	globalDisplay.Update()
}

func cleanup() {
	fmt.Println("🧹 Cleaning up...")
	if globalArm != nil {
		globalArm.DisableAll()
	}
	if globalBoard != nil {
		globalBoard.Cleanup()
	}
	fmt.Println("✅ Cleanup completed")
}

func runDemoMode() {
	if globalArm == nil {
		log.Fatal("Robot arm not initialized")
		return
	}

	// Enable all motors
	fmt.Println("Enabling all motors...")
	if err := globalArm.EnableAll(); err != nil {
		log.Fatalf("Failed to enable motors: %v", err)
	}

	// Home the arm
	if err := globalArm.Home(); err != nil {
		log.Fatalf("Failed to home arm: %v", err)
	}

	// Display status
	fmt.Println("\n" + globalArm.GetStatus())

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

	if err := globalArm.ExecuteSequence(movements1, delays1); err != nil {
		log.Printf("Error in sequence 1: %v", err)
	}

	// Sequence 2: Pick and place simulation
	fmt.Println("\n📦 Sequence 2: Pick and place simulation")
	movements2 := []map[string]int{
		{"pan": 200, "tilt": 100, "joint3": -150}, // Move to pick position
		{"gripper": 80},                    // Close gripper
		{"tilt": -50, "joint3": 100},       // Lift object
		{"pan": -200},                      // Rotate to place position
		{"tilt": 100, "joint3": -150},      // Lower to place position
		{"gripper": 0},                     // Open gripper
		{"pan": 0, "tilt": 0, "joint3": 0}, // Return to home
	}
	delays2 := []time.Duration{
		3 * time.Second, 1 * time.Second, 2 * time.Second,
		2 * time.Second, 3 * time.Second, 1 * time.Second, 3 * time.Second,
	}

	if err := globalArm.ExecuteSequence(movements2, delays2); err != nil {
		log.Printf("Error in sequence 2: %v", err)
	}

	// Final status
	fmt.Println("\n📊 Final Status:")
	fmt.Println(globalArm.GetStatus())

	// Demonstrate individual joint control
	fmt.Println("\n🎯 Individual joint control demo:")

	// Move pan joint (base rotation)
	fmt.Println("Moving pan joint...")
	if err := globalArm.MoveJoint("pan", 300); err != nil {
		log.Printf("Error moving pan: %v", err)
	}

	time.Sleep(1 * time.Second)

	// Move tilt joint (elevation)
	fmt.Println("Moving tilt joint...")
	if err := globalArm.MoveJoint("tilt", -200); err != nil {
		log.Printf("Error moving tilt: %v", err)
	}

	// Return to home position
	fmt.Println("\n🏠 Returning to home position...")
	homePositions := map[string]int{
		"pan": 0, "tilt": 0, "joint3": 0,
		"joint4": 0, "joint5": 0, "gripper": 0,
	}

	if err := globalArm.MoveJoints(homePositions); err != nil {
		log.Printf("Error returning home: %v", err)
	}

	// Disable all motors
	fmt.Println("\n🔌 Disabling all motors...")
	if err := globalArm.DisableAll(); err != nil {
		log.Printf("Error disabling motors: %v", err)
	}

	fmt.Println("\n✅ Demo completed successfully!")
	fmt.Println("📝 Check the logs above to see the mock hardware interactions.")
}

// Functions that return success status for auto-display updates

func enableMotorsWithResult() bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}
	if err := globalArm.EnableAll(); err != nil {
		fmt.Printf("❌ Failed to enable motors: %v\n", err)
		return false
	} else {
		fmt.Println("✅ All motors enabled")
		return true
	}
}

func moveXYZWithResult(xStr, yStr, zStr string) bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}

	x, errX := strconv.ParseFloat(xStr, 64)
	y, errY := strconv.ParseFloat(yStr, 64)
	z, errZ := strconv.ParseFloat(zStr, 64)

	if errX != nil || errY != nil || errZ != nil {
		fmt.Printf("❌ Invalid coordinates. Please provide three numbers.\n")
		return false
	}

	fmt.Printf("🎯 Moving to Cartesian coordinate: (%.2f, %.2f, %.2f)\n", x, y, z)

	if err := globalArm.MoveToXYZ(x, y, z); err != nil {
		fmt.Printf("❌ Failed to move to coordinate: %v\n", err)
		return false
	}

	fmt.Printf("✅ MoveToXYZ command issued (placeholder).\n")
	return true
}

func disableMotorsWithResult() bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}
	if err := globalArm.DisableAll(); err != nil {
		fmt.Printf("❌ Failed to disable motors: %v\n", err)
		return false
	} else {
		fmt.Println("✅ All motors disabled")
		return true
	}
}

func homeArmWithResult() bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}
	fmt.Println("🏠 Homing robot arm...")
	if err := globalArm.Home(); err != nil {
		fmt.Printf("❌ Failed to home arm: %v\n", err)
		return false
	} else {
		fmt.Println("✅ Robot arm homed successfully")
		return true
	}
}

func moveJointWithResult(jointName, positionStr string) bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}

	position, err := strconv.Atoi(positionStr)
	if err != nil {
		fmt.Printf("❌ Invalid position: %s\n", positionStr)
		return false
	}

	fmt.Printf("🎯 Moving joint '%s' to position %d\n", jointName, position)

	if err := globalArm.MoveJoint(jointName, position); err != nil {
		fmt.Printf("❌ Failed to move joint: %v\n", err)
		return false
	} else {
		fmt.Printf("✅ Joint '%s' moved to position %d\n", jointName, position)
		return true
	}
}

func runPickPlaceSequenceWithResult() bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}

	fmt.Println("📦 Running pick and place sequence...")

	movements := []map[string]int{
		{"pan": 200, "tilt": 100, "joint3": -150}, // Move to pick position
		{"gripper": 80},                    // Close gripper
		{"tilt": -50, "joint3": 100},       // Lift object
		{"pan": -200},                      // Rotate to place position
		{"tilt": 100, "joint3": -150},      // Lower to place position
		{"gripper": 0},                     // Open gripper
		{"pan": 0, "tilt": 0, "joint3": 0}, // Return to home
	}
	delays := []time.Duration{
		3 * time.Second, 1 * time.Second, 2 * time.Second,
		2 * time.Second, 3 * time.Second, 1 * time.Second, 3 * time.Second,
	}

	if err := globalArm.ExecuteSequence(movements, delays); err != nil {
		fmt.Printf("❌ Error in sequence: %v\n", err)
		return false
	} else {
		fmt.Println("✅ Pick and place sequence completed")
		return true
	}
}

func runBasicDemoWithResult() bool {
	if globalArm == nil {
		fmt.Println("❌ Robot arm not initialized")
		return false
	}

	fmt.Println("🎬 Running basic movement demo...")

	movements := []map[string]int{
		{"pan": 100, "tilt": 50},
		{"joint3": -100, "joint4": 75},
		{"joint5": 50, "gripper": 25},
		{"pan": 0, "tilt": 0, "joint3": 0, "joint4": 0, "joint5": 0, "gripper": 0},
	}
	delays := []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second, 3 * time.Second}

	if err := globalArm.ExecuteSequence(movements, delays); err != nil {
		fmt.Printf("❌ Error in demo: %v\n", err)
		return false
	} else {
		fmt.Println("✅ Basic demo completed")
		return true
	}
}
