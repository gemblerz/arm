package robot

import (
	"fmt"
	"sync"
	"time"

	"github.com/gemblerz/arm/pkg/kinematics"
	"github.com/gemblerz/arm/pkg/stepper"
)

// DisplayInterface defines the interface for robot display
type DisplayInterface interface {
	UpdateStatus(status DisplayStatus)
	ShowStartupScreen()
	ShowError(errorMsg string)
	ShowCalibration(jointID string, progress float64)
	Start()
	Stop()
}

// DisplayStatus represents the robot status for display
type DisplayStatus struct {
	IsHomed    bool
	IsEnabled  bool
	IsMoving   bool
	Joints     []DisplayJointStatus
	CurrentOp  string
	Timestamp  time.Time
}

// DisplayJointStatus represents joint status for display
type DisplayJointStatus struct {
	ID       string
	Position int
	Target   int
	MinPos   int
	MaxPos   int
	Enabled  bool
}

// Joint represents a single joint in the robot arm
type Joint struct {
	ID     string
	Motor  stepper.StepperMotor
	MinPos int // Minimum position in steps
	MaxPos int // Maximum position in steps
	currentTarget int // Current target position
}

// Arm represents a 6-DOF robot arm
type Arm struct {
	joints           []Joint
	mutex            sync.RWMutex
	homed            bool
	display          DisplayInterface
	isMoving         bool
	currentOperation string
	kinematicsSolver *kinematics.KinematicsSolver
}

// NewArm creates a new robot arm with the specified joints
func NewArm(joints []Joint, solver *kinematics.KinematicsSolver) *Arm {
	return &Arm{
		joints:           joints,
		homed:            false,
		isMoving:         false,
		currentOperation: "Initializing",
		kinematicsSolver: solver,
	}
}

// SetDisplay sets the display interface for the robot arm
func (a *Arm) SetDisplay(display DisplayInterface) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.display = display
	if display != nil {
		display.Start()
		a.updateDisplay()
	}
}

// updateDisplay sends current status to display
func (a *Arm) updateDisplay() {
	if a.display == nil {
		return
	}

	status := DisplayStatus{
		IsHomed:   a.homed,
		IsEnabled: a.areAllEnabled(),
		IsMoving:  a.isMoving,
		CurrentOp: a.currentOperation,
		Joints:    make([]DisplayJointStatus, len(a.joints)),
		Timestamp: time.Now(),
	}

	for i, joint := range a.joints {
		status.Joints[i] = DisplayJointStatus{
			ID:       joint.ID,
			Position: joint.Motor.GetPosition(),
			Target:   joint.currentTarget,
			MinPos:   joint.MinPos,
			MaxPos:   joint.MaxPos,
			Enabled:  joint.Motor.IsEnabled(),
		}
	}

	a.display.UpdateStatus(status)
}

// areAllEnabled checks if all joints are enabled
func (a *Arm) areAllEnabled() bool {
	for _, joint := range a.joints {
		if !joint.Motor.IsEnabled() {
			return false
		}
	}
	return len(a.joints) > 0
}

// GetJoint returns a joint by its ID
func (a *Arm) GetJoint(id string) (*Joint, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	for i := range a.joints {
		if a.joints[i].ID == id {
			return &a.joints[i], nil
		}
	}
	return nil, fmt.Errorf("joint with ID '%s' not found", id)
}

// EnableAll enables all joints in the arm
func (a *Arm) EnableAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	a.currentOperation = "Enabling motors"
	a.updateDisplay()
	
	for i := range a.joints {
		if err := a.joints[i].Motor.Enable(); err != nil {
			return fmt.Errorf("failed to enable joint %s: %w", a.joints[i].ID, err)
		}
	}
	
	a.currentOperation = "Motors enabled"
	a.updateDisplay()
	return nil
}

// DisableAll disables all joints in the arm
func (a *Arm) DisableAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	for i := range a.joints {
		if err := a.joints[i].Motor.Disable(); err != nil {
			return fmt.Errorf("failed to disable joint %s: %w", a.joints[i].ID, err)
		}
	}
	return nil
}

// Home performs a homing sequence for all joints
func (a *Arm) Home() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	a.currentOperation = "Homing robot"
	a.updateDisplay()
	
	fmt.Println("Starting homing sequence...")
	
	// Simple homing: move all joints to position 0
	for i := range a.joints {
		fmt.Printf("Homing joint %s...\n", a.joints[i].ID)
		
		// Show calibration progress on display
		if a.display != nil {
			a.display.ShowCalibration(a.joints[i].ID, 0.0)
		}
		
		// Simulate homing progress
		for progress := 0.0; progress <= 1.0; progress += 0.1 {
			if a.display != nil {
				a.display.ShowCalibration(a.joints[i].ID, progress)
			}
			time.Sleep(50 * time.Millisecond) // Simulate homing time
		}
		
		a.joints[i].Motor.SetPosition(0)
		a.joints[i].currentTarget = 0
	}
	
	a.homed = true
	a.currentOperation = "Robot homed"
	a.updateDisplay()
	
	fmt.Println("Homing sequence completed")
	return nil
}

// IsHomed returns true if the arm has been homed
func (a *Arm) IsHomed() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.homed
}

// MoveJoint moves a specific joint to a target position
func (a *Arm) MoveJoint(jointID string, targetPosition int) error {
	joint, err := a.GetJoint(jointID)
	if err != nil {
		return err
	}
	
	// Check position limits
	if targetPosition < joint.MinPos || targetPosition > joint.MaxPos {
		return fmt.Errorf("target position %d is out of range [%d, %d] for joint %s",
			targetPosition, joint.MinPos, joint.MaxPos, jointID)
	}
	
	fmt.Printf("Moving joint %s to position %d\n", jointID, targetPosition)
	return joint.Motor.MoveTo(targetPosition)
}

// MoveJoints moves multiple joints simultaneously to target positions
func (a *Arm) MoveJoints(targets map[string]int) error {
	if !a.IsHomed() {
		return fmt.Errorf("arm must be homed before moving")
	}
	
	// Validate all targets first
	for jointID, targetPos := range targets {
		joint, err := a.GetJoint(jointID)
		if err != nil {
			return err
		}
		
		if targetPos < joint.MinPos || targetPos > joint.MaxPos {
			return fmt.Errorf("target position %d is out of range [%d, %d] for joint %s",
				targetPos, joint.MinPos, joint.MaxPos, jointID)
		}
	}
	
	// Move all joints concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, len(targets))
	
	for jointID, targetPos := range targets {
		wg.Add(1)
		go func(id string, pos int) {
			defer wg.Done()
			if err := a.MoveJoint(id, pos); err != nil {
				errChan <- fmt.Errorf("joint %s: %w", id, err)
			}
		}(jointID, targetPos)
	}
	
	wg.Wait()
	close(errChan)
	
	// Check for errors
	for err := range errChan {
		return err
	}
	
	return nil
}

// GetPositions returns the current positions of all joints
func (a *Arm) GetPositions() map[string]int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	positions := make(map[string]int)
	for _, joint := range a.joints {
		positions[joint.ID] = joint.Motor.GetPosition()
	}
	return positions
}

// ExecuteSequence executes a sequence of movements with delays
func (a *Arm) ExecuteSequence(sequence []map[string]int, delays []time.Duration) error {
	if len(sequence) != len(delays) {
		return fmt.Errorf("sequence length (%d) must match delays length (%d)", 
			len(sequence), len(delays))
	}
	
	for i, targets := range sequence {
		fmt.Printf("Executing movement %d/%d\n", i+1, len(sequence))
		
		if err := a.MoveJoints(targets); err != nil {
			return fmt.Errorf("movement %d failed: %w", i+1, err)
		}
		
		if i < len(delays) {
			fmt.Printf("Waiting %v before next movement...\n", delays[i])
			time.Sleep(delays[i])
		}
	}
	
	return nil
}

// MoveToXYZ moves the arm's end-effector to a target Cartesian coordinate.
// This is a placeholder for the full inverse kinematics implementation.
func (a *Arm) MoveToXYZ(x, y, z float64) error {
	if a.kinematicsSolver == nil {
		return fmt.Errorf("kinematics solver is not initialized")
	}

	fmt.Printf("Attempting to move to Cartesian coordinate: (%.2f, %.2f, %.2f)\n", x, y, z)

	// TODO:
	// 1. Convert current motor positions (steps) to joint angles (radians).
	// 2. Use inverse kinematics to calculate the target joint angles for the given x, y, z.
	// 3. Convert target joint angles back to motor steps.
	// 4. Command the motors to move to the new step positions using MoveJoints.

	fmt.Println("Placeholder: Inverse kinematics calculation not yet implemented.")
	return nil
}

// GetCurrentCartesianPosition calculates and returns the current end-effector position.
// This is a placeholder for the full forward kinematics implementation.
func (a *Arm) GetCurrentCartesianPosition() (kinematics.CartesianPoint, error) {
	if a.kinematicsSolver == nil {
		return kinematics.CartesianPoint{}, fmt.Errorf("kinematics solver is not initialized")
	}

	// TODO:
	// 1. Get current motor positions (steps) for all joints.
	// 2. Convert these step positions to joint angles (radians).
	//    This will require knowing the steps-per-degree for each motor.
	// 3. Pass the joint angles to the forward kinematics solver.

	// Using placeholder joint angles for now.
	placeholderAngles := []float64{0, 0, 0, 0, 0, 0}
	currentPos := a.kinematicsSolver.ForwardKinematics(placeholderAngles)

	return currentPos, nil
}

// GetStatus returns a formatted status of the entire arm
func (a *Arm) GetStatus() string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	status := fmt.Sprintf("Robot Arm Status (Homed: %t):\n", a.homed)
	positions := a.GetPositions()
	
	for _, joint := range a.joints {
		enabled := "DISABLED"
		if joint.Motor.IsEnabled() {
			enabled = "ENABLED"
		}
		status += fmt.Sprintf("  Joint %s: %s, Position: %d, Range: [%d, %d]\n",
			joint.ID, enabled, positions[joint.ID], joint.MinPos, joint.MaxPos)
	}
	
	// Add Cartesian position to the status if solver is available
	if a.kinematicsSolver != nil {
		pos, err := a.GetCurrentCartesianPosition()
		if err != nil {
			status += fmt.Sprintf("  End-Effector (X,Y,Z): Error - %v\n", err)
		} else {
			status += fmt.Sprintf("  End-Effector (X,Y,Z): (%.2f, %.2f, %.2f) (placeholder)\n", pos.X, pos.Y, pos.Z)
		}
	}

	return status
}

// UpdateDisplayStatus updates the robot status on the display
func (a *Arm) UpdateDisplayStatus() {
	if a.display == nil {
		return
	}
	
	status := DisplayStatus{
		IsHomed:   a.homed,
		IsEnabled: a.isMoving,
		IsMoving:  a.isMoving,
		CurrentOp: a.currentOperation,
		Timestamp: time.Now(),
		Joints:    make([]DisplayJointStatus, len(a.joints)),
	}
	
	for i, joint := range a.joints {
		status.Joints[i] = DisplayJointStatus{
			ID:       joint.ID,
			Position: joint.Motor.GetPosition(),
			Target:   joint.currentTarget,
			MinPos:   joint.MinPos,
			MaxPos:   joint.MaxPos,
			Enabled:  joint.Motor.IsEnabled(),
		}
	}
	
	a.display.UpdateStatus(status)
}

// StartDisplay starts the robot display
func (a *Arm) StartDisplay() {
	if a.display == nil {
		return
	}
	a.display.Start()
	a.display.ShowStartupScreen()
	a.UpdateDisplayStatus()
}

// StopDisplay stops the robot display
func (a *Arm) StopDisplay() {
	if a.display == nil {
		return
	}
	a.display.Stop()
}

// MoveJointWithDisplay moves a joint and updates the display
func (a *Arm) MoveJointWithDisplay(jointID string, targetPosition int) error {
	if err := a.MoveJoint(jointID, targetPosition); err != nil {
		return err
	}
	
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Update the current target position
	for i := range a.joints {
		if a.joints[i].ID == jointID {
			a.joints[i].currentTarget = targetPosition
			break
		}
	}
	
	// Update the display with the new status
	a.updateDisplay()
	return nil
}

// MoveJointsWithDisplay moves multiple joints and updates the display
func (a *Arm) MoveJointsWithDisplay(targets map[string]int) error {
	if err := a.MoveJoints(targets); err != nil {
		return err
	}
	
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Update the current target positions
	for jointID, targetPosition := range targets {
		for i := range a.joints {
			if a.joints[i].ID == jointID {
				a.joints[i].currentTarget = targetPosition
				break
			}
		}
	}
	
	// Update the display with the new status
	a.UpdateDisplayStatus()
	return nil
}

// HomeWithDisplay performs homing and updates the display
func (a *Arm) HomeWithDisplay() error {
	if err := a.Home(); err != nil {
		return err
	}
	
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Update the current target positions to 0
	for i := range a.joints {
		a.joints[i].currentTarget = 0
	}
	
	// Update the display with the new status
	a.UpdateDisplayStatus()
	return nil
}

// ExecuteSequenceWithDisplay executes a sequence and updates the display
func (a *Arm) ExecuteSequenceWithDisplay(sequence []map[string]int, delays []time.Duration) error {
	if err := a.ExecuteSequence(sequence, delays); err != nil {
		return err
	}
	
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Update the current target positions based on the last sequence
	lastTargets := sequence[len(sequence)-1]
	for jointID, targetPosition := range lastTargets {
		for i := range a.joints {
			if a.joints[i].ID == jointID {
				a.joints[i].currentTarget = targetPosition
				break
			}
		}
	}
	
	// Update the display with the new status
	a.UpdateDisplayStatus()
	return nil
}
