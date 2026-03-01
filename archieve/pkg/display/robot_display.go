package display

import (
	"fmt"
	"time"

	"golang.org/x/image/font/basicfont"
)

// RobotDisplayManager manages the robot arm display interface
type RobotDisplayManager struct {
	display     DisplayInterface
	updateRate  time.Duration
	running     bool
	statusChan  chan RobotStatus
	stopChan    chan bool
	lastStatus  RobotStatus
}

// RobotStatus represents the current state of the robot arm
type RobotStatus struct {
	IsHomed    bool
	IsEnabled  bool
	IsMoving   bool
	Joints     []JointStatus
	CurrentOp  string
	Timestamp  time.Time
}

// JointStatus represents the status of a single joint
type JointStatus struct {
	ID       string
	Position int
	Target   int
	MinPos   int
	MaxPos   int
	Enabled  bool
}

// NewRobotDisplayManager creates a new display manager
func NewRobotDisplayManager(display DisplayInterface, updateRate time.Duration) *RobotDisplayManager {
	return &RobotDisplayManager{
		display:    display,
		updateRate: updateRate,
		statusChan: make(chan RobotStatus, 10),
		stopChan:   make(chan bool, 1),
	}
}

// Start begins the display update loop
func (dm *RobotDisplayManager) Start() {
	dm.running = true
	go dm.updateLoop()
}

// Stop stops the display update loop
func (dm *RobotDisplayManager) Stop() {
	if dm.running {
		dm.running = false
		dm.stopChan <- true
	}
}

// UpdateStatus updates the robot status for display
func (dm *RobotDisplayManager) UpdateStatus(status RobotStatus) {
	status.Timestamp = time.Now()
	select {
	case dm.statusChan <- status:
		// Status updated
	default:
		// Channel full, skip this update
	}
}

// updateLoop runs the display update loop
func (dm *RobotDisplayManager) updateLoop() {
	ticker := time.NewTicker(dm.updateRate)
	defer ticker.Stop()

	for dm.running {
		select {
		case status := <-dm.statusChan:
			dm.lastStatus = status
			dm.renderDisplay()
		case <-ticker.C:
			// Periodic update even without status change
			dm.renderDisplay()
		case <-dm.stopChan:
			return
		}
	}
}

// renderDisplay renders the current robot status to the display
func (dm *RobotDisplayManager) renderDisplay() {
	dm.display.Clear()

	// Header with timestamp
	timeStr := time.Now().Format("15:04:05")
	dm.display.DrawText(fmt.Sprintf("Robot %s", timeStr), 2, 10, basicfont.Face7x13)

	// Status indicators
	y := 25
	statusText := dm.getStatusText()
	dm.display.DrawText(statusText, 2, y, basicfont.Face7x13)

	// Current operation
	if dm.lastStatus.CurrentOp != "" {
		y += 15
		dm.display.DrawText(fmt.Sprintf("Op: %s", dm.lastStatus.CurrentOp), 2, y, basicfont.Face7x13)
	}

	// Joint positions (show first 3 joints due to space)
	if len(dm.lastStatus.Joints) > 0 {
		y += 15
		dm.renderJointInfo(y)
	}

	dm.display.Update()
}

// getStatusText returns a formatted status string
func (dm *RobotDisplayManager) getStatusText() string {
	status := "STS:"
	if dm.lastStatus.IsHomed {
		status += "H"
	} else {
		status += "-"
	}
	if dm.lastStatus.IsEnabled {
		status += "E"
	} else {
		status += "-"
	}
	if dm.lastStatus.IsMoving {
		status += "M"
	} else {
		status += "-"
	}
	return status
}

// renderJointInfo renders joint position information
func (dm *RobotDisplayManager) renderJointInfo(startY int) {
	maxJoints := 3 // Limit to 3 joints for display space
	for i, joint := range dm.lastStatus.Joints {
		if i >= maxJoints {
			break
		}

		y := startY + (i * 12)
		
		// Joint name and position
		text := fmt.Sprintf("%s:%d", joint.ID[:3], joint.Position)
		dm.display.DrawText(text, 2, y, basicfont.Face7x13)

		// Progress bar showing position relative to range
		progress := dm.calculateJointProgress(joint)
		barX := 45
		barY := y - 8
		dm.display.DrawProgressBar(barX, barY, 40, 8, progress)
		
		// Target indicator
		if joint.Target != joint.Position {
			targetProgress := dm.calculateJointTargetProgress(joint)
			targetX := barX + int(40*targetProgress)
			dm.display.DrawLine(targetX, barY-1, targetX, barY+9)
		}
	}
}

// calculateJointProgress calculates the progress of a joint within its range
func (dm *RobotDisplayManager) calculateJointProgress(joint JointStatus) float64 {
	if joint.MaxPos == joint.MinPos {
		return 0.5
	}
	
	progress := float64(joint.Position-joint.MinPos) / float64(joint.MaxPos-joint.MinPos)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	return progress
}

// calculateJointTargetProgress calculates the target progress of a joint
func (dm *RobotDisplayManager) calculateJointTargetProgress(joint JointStatus) float64 {
	if joint.MaxPos == joint.MinPos {
		return 0.5
	}
	
	progress := float64(joint.Target-joint.MinPos) / float64(joint.MaxPos-joint.MinPos)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	return progress
}

// ShowStartupScreen displays a startup/welcome screen
func (dm *RobotDisplayManager) ShowStartupScreen() {
	dm.display.Clear()
	
	// Title
	dm.display.DrawText("6-DOF Robot Arm", 10, 15, basicfont.Face7x13)
	dm.display.DrawText("DM556T Controller", 5, 30, basicfont.Face7x13)
	
	// Version info
	dm.display.DrawText("v1.0 - Coral Board", 5, 45, basicfont.Face7x13)
	
	// Progress animation
	for i := 0; i < 100; i += 10 {
		dm.display.DrawProgressBar(10, 55, 108, 6, float64(i)/100.0)
		dm.display.Update()
		time.Sleep(100 * time.Millisecond)
	}
	
	time.Sleep(1 * time.Second)
}

// ShowError displays an error message
func (dm *RobotDisplayManager) ShowError(errorMsg string) {
	dm.display.Clear()
	
	dm.display.DrawText("ERROR:", 2, 15, basicfont.Face7x13)
	
	// Word wrap the error message
	maxChars := 18 // Approximate characters per line
	words := splitString(errorMsg, maxChars)
	
	y := 30
	for i, line := range words {
		if i >= 3 { // Max 3 lines
			break
		}
		dm.display.DrawText(line, 2, y, basicfont.Face7x13)
		y += 12
	}
	
	dm.display.Update()
}

// ShowCalibration displays calibration/homing status
func (dm *RobotDisplayManager) ShowCalibration(jointID string, progress float64) {
	dm.display.Clear()
	
	dm.display.DrawText("HOMING...", 25, 15, basicfont.Face7x13)
	dm.display.DrawText(fmt.Sprintf("Joint: %s", jointID), 15, 30, basicfont.Face7x13)
	
	// Progress bar
	dm.display.DrawProgressBar(10, 40, 108, 8, progress)
	
	// Percentage
	dm.display.DrawText(fmt.Sprintf("%d%%", int(progress*100)), 55, 55, basicfont.Face7x13)
	
	dm.display.Update()
}

// Helper function to split strings for word wrapping
func splitString(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	
	var lines []string
	for len(text) > maxLen {
		// Find the last space before maxLen
		breakPoint := maxLen
		for i := maxLen - 1; i >= 0; i-- {
			if text[i] == ' ' {
				breakPoint = i
				break
			}
		}
		
		lines = append(lines, text[:breakPoint])
		text = text[breakPoint+1:]
	}
	
	if len(text) > 0 {
		lines = append(lines, text)
	}
	
	return lines
}
