package display

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/font/basicfont"
)

// DisplayScreen represents different display modes
type DisplayScreen int

const (
	ScreenStatus DisplayScreen = iota
	ScreenJoints
	ScreenAngles
	ScreenDiagnostics
	ScreenSequence
	ScreenDriverInfo
)

// AdvancedRobotDisplay provides enhanced display functionality
type AdvancedRobotDisplay struct {
	manager       *RobotDisplayManager
	currentScreen DisplayScreen
	screenNames   []string
	lastSwitch    time.Time
	autoSwitch    bool
	switchInterval time.Duration
}

// NewAdvancedRobotDisplay creates an enhanced display with multiple screens
func NewAdvancedRobotDisplay(display DisplayInterface, updateRate time.Duration) *AdvancedRobotDisplay {
	manager := NewRobotDisplayManager(display, updateRate)
	
	return &AdvancedRobotDisplay{
		manager:        manager,
		currentScreen:  ScreenStatus,
		screenNames:    []string{"Status", "Joints", "Angles", "Diagnostics", "Sequence", "Drivers"},
		lastSwitch:     time.Now(),
		autoSwitch:     true,
		switchInterval: 5 * time.Second,
	}
}

// Start begins the advanced display system
func (ad *AdvancedRobotDisplay) Start() {
	ad.manager.Start()
	go ad.screenManager()
}

// Stop stops the advanced display system
func (ad *AdvancedRobotDisplay) Stop() {
	ad.manager.Stop()
}

// SetScreen manually sets the current screen
func (ad *AdvancedRobotDisplay) SetScreen(screen DisplayScreen) {
	ad.currentScreen = screen
	ad.lastSwitch = time.Now()
}

// EnableAutoSwitch enables/disables automatic screen switching
func (ad *AdvancedRobotDisplay) EnableAutoSwitch(enable bool, interval time.Duration) {
	ad.autoSwitch = enable
	ad.switchInterval = interval
}

// screenManager handles automatic screen switching
func (ad *AdvancedRobotDisplay) screenManager() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ad.autoSwitch && time.Since(ad.lastSwitch) > ad.switchInterval {
				ad.nextScreen()
			}
		}
	}
}

// nextScreen advances to the next screen
func (ad *AdvancedRobotDisplay) nextScreen() {
	ad.currentScreen = DisplayScreen((int(ad.currentScreen) + 1) % len(ad.screenNames))
	ad.lastSwitch = time.Now()
}

// UpdateStatus updates robot status and renders current screen
func (ad *AdvancedRobotDisplay) UpdateStatus(status RobotStatus) {
	ad.manager.lastStatus = status
	ad.renderCurrentScreen()
}

// renderCurrentScreen renders the currently selected screen
func (ad *AdvancedRobotDisplay) renderCurrentScreen() {
	switch ad.currentScreen {
	case ScreenStatus:
		ad.renderStatusScreen()
	case ScreenJoints:
		ad.renderJointsScreen()
	case ScreenAngles:
		ad.renderAnglesScreen()
	case ScreenDiagnostics:
		ad.renderDiagnosticsScreen()
	case ScreenSequence:
		ad.renderSequenceScreen()
	case ScreenDriverInfo:
		ad.renderDriverScreen()
	}
}

// Screen 1: Enhanced Status Overview
func (ad *AdvancedRobotDisplay) renderStatusScreen() {
	display := ad.manager.display
	status := ad.manager.lastStatus
	
	display.Clear()
	
	// Header with screen indicator
	timeStr := time.Now().Format("15:04:05")
	display.DrawText(fmt.Sprintf("STS %s [1/6]", timeStr), 2, 10, basicfont.Face7x13)
	
	// System status with icons
	y := 25
	statusIcons := ad.getStatusIcons()
	display.DrawText(statusIcons, 2, y, basicfont.Face7x13)
	
	// Current operation with progress
	if status.CurrentOp != "" {
		y += 15
		opText := status.CurrentOp
		if len(opText) > 16 {
			opText = opText[:13] + "..."
		}
		display.DrawText(opText, 2, y, basicfont.Face7x13)
	}
	
	// System temperature simulation (in real system, read from sensors)
	y += 15
	temp := 42 + int(math.Sin(float64(time.Now().Unix())/10)*3) // Simulated temp
	display.DrawText(fmt.Sprintf("Temp: %d°C", temp), 2, y, basicfont.Face7x13)
	
	// System uptime
	y += 12
	uptime := time.Since(status.Timestamp).Round(time.Second)
	display.DrawText(fmt.Sprintf("Up: %v", uptime), 2, y, basicfont.Face7x13)
	
	display.Update()
}

// Screen 2: Joint Positions with Enhanced Visualization
func (ad *AdvancedRobotDisplay) renderJointsScreen() {
	display := ad.manager.display
	status := ad.manager.lastStatus
	
	display.Clear()
	
	// Header
	display.DrawText("JOINTS [2/6]", 2, 10, basicfont.Face7x13)
	
	// Show 3 joints at a time with detailed info
	for i := 0; i < 3 && i < len(status.Joints); i++ {
		joint := status.Joints[i]
		y := 25 + (i * 15)
		
		// Joint name and position
		text := fmt.Sprintf("%s:%d", joint.ID[:3], joint.Position)
		display.DrawText(text, 2, y, basicfont.Face7x13)
		
		// Enhanced progress bar with target
		progress := ad.manager.calculateJointProgress(joint)
		barX := 50
		barY := y - 8
		
		// Draw progress bar background
		display.DrawRectangle(barX, barY, 70, 8, false)
		
		// Fill current position
		fillWidth := int(68 * progress)
		if fillWidth > 0 {
			display.DrawRectangle(barX+1, barY+1, fillWidth, 6, true)
		}
		
		// Target marker
		if joint.Target != joint.Position {
			targetProgress := ad.manager.calculateJointTargetProgress(joint)
			targetX := barX + int(68*targetProgress)
			// Draw target as a vertical line
			for ty := barY; ty < barY+8; ty++ {
				display.SetPixel(targetX, ty, true)
			}
		}
		
		// Position percentage
		percentage := int(progress * 100)
		display.DrawText(fmt.Sprintf("%d%%", percentage), barX+72, y, basicfont.Face7x13)
	}
	
	display.Update()
}

// Screen 3: Angle Display (Convert steps to degrees)
func (ad *AdvancedRobotDisplay) renderAnglesScreen() {
	display := ad.manager.display
	status := ad.manager.lastStatus
	
	display.Clear()
	
	// Header
	display.DrawText("ANGLES [3/6]", 2, 10, basicfont.Face7x13)
	
	// Show angles for pan/tilt joints
	for i := 0; i < 3 && i < len(status.Joints); i++ {
		joint := status.Joints[i]
		y := 25 + (i * 15)
		
		// Convert steps to degrees (assuming 200 steps per revolution, 32x microstepping)
		stepsPerDegree := (200 * 32) / 360.0 // 6400 steps = 360°
		degrees := float64(joint.Position) / stepsPerDegree
		
		// Joint name and angle
		display.DrawText(fmt.Sprintf("%s:", joint.ID[:3]), 2, y, basicfont.Face7x13)
		display.DrawText(fmt.Sprintf("%.1f°", degrees), 35, y, basicfont.Face7x13)
		
		// Target angle if different
		if joint.Target != joint.Position {
			targetDegrees := float64(joint.Target) / stepsPerDegree
			display.DrawText(fmt.Sprintf("→%.1f°", targetDegrees), 80, y, basicfont.Face7x13)
		}
	}
	
	// Show workspace visualization (simple arm diagram)
	ad.drawArmDiagram(45)
	
	display.Update()
}

// Screen 4: Diagnostics and Performance
func (ad *AdvancedRobotDisplay) renderDiagnosticsScreen() {
	display := ad.manager.display
	status := ad.manager.lastStatus
	
	display.Clear()
	
	// Header
	display.DrawText("DIAG [4/6]", 2, 10, basicfont.Face7x13)
	
	y := 25
	
	// Motor status summary
	enabled := 0
	for _, joint := range status.Joints {
		if joint.Enabled {
			enabled++
		}
		// In a real system, you'd track if motors are currently moving
	}
	
	display.DrawText(fmt.Sprintf("Motors: %d/%d ON", enabled, len(status.Joints)), 2, y, basicfont.Face7x13)
	
	y += 12
	// Communication status
	display.DrawText("I2C: OK  GPIO: OK", 2, y, basicfont.Face7x13)
	
	y += 12
	// Error count (simulated)
	display.DrawText("Errors: 0", 2, y, basicfont.Face7x13)
	
	y += 12
	// Last movement time
	display.DrawText(fmt.Sprintf("Last Move: %s", status.Timestamp.Format("15:04")), 2, y, basicfont.Face7x13)
	
	display.Update()
}

// Screen 5: Sequence Progress
func (ad *AdvancedRobotDisplay) renderSequenceScreen() {
	display := ad.manager.display
	status := ad.manager.lastStatus
	
	display.Clear()
	
	// Header
	display.DrawText("SEQ [5/6]", 2, 10, basicfont.Face7x13)
	
	y := 25
	
	// Current operation
	if status.CurrentOp != "" {
		display.DrawText("Current:", 2, y, basicfont.Face7x13)
		y += 12
		opText := status.CurrentOp
		if len(opText) > 18 {
			opText = opText[:15] + "..."
		}
		display.DrawText(opText, 2, y, basicfont.Face7x13)
	}
	
	y += 15
	// Sequence progress (simulated)
	progress := 0.6 // In real system, track actual progress
	display.DrawText("Progress:", 2, y, basicfont.Face7x13)
	y += 10
	display.DrawProgressBar(2, y, 124, 8, progress)
	
	y += 15
	display.DrawText(fmt.Sprintf("Step: 3/5 (%.0f%%)", progress*100), 2, y, basicfont.Face7x13)
	
	display.Update()
}

// Screen 6: Driver Information
func (ad *AdvancedRobotDisplay) renderDriverScreen() {
	display := ad.manager.display
	
	display.Clear()
	
	// Header
	display.DrawText("DRIVERS [6/6]", 2, 10, basicfont.Face7x13)
	
	y := 25
	
	// Driver assignments
	drivers := []string{"DM556T x2", "DM320T x1", "DRV8825 x3"}
	for i, driver := range drivers {
		display.DrawText(driver, 2, y+(i*12), basicfont.Face7x13)
	}
	
	y += 40
	// Current limits (simulated)
	display.DrawText("Current:", 2, y, basicfont.Face7x13)
	display.DrawText("5.6A/3.0A/2.2A", 50, y, basicfont.Face7x13)
	
	display.Update()
}

// Helper methods
func (ad *AdvancedRobotDisplay) getStatusIcons() string {
	status := ad.manager.lastStatus
	icons := "["
	
	if status.IsHomed {
		icons += "H"
	} else {
		icons += "-"
	}
	
	if status.IsEnabled {
		icons += "E"
	} else {
		icons += "-"
	}
	
	if status.IsMoving {
		icons += "M"
	} else {
		icons += "-"
	}
	
	icons += "] "
	
	// Add connection status
	icons += "WiFi:● I2C:●"
	
	return icons
}

// Draw simple arm diagram
func (ad *AdvancedRobotDisplay) drawArmDiagram(startY int) {
	display := ad.manager.display
	
	// Simple stick figure representation of robot arm
	centerX := 64
	baseY := startY + 15
	
	// Base
	display.DrawRectangle(centerX-5, baseY, 10, 3, true)
	
	// Joint 1 (pan) - represented as a line
	display.DrawLine(centerX, baseY, centerX+15, baseY-10)
	
	// Joint 2 (tilt) - another line
	display.DrawLine(centerX+15, baseY-10, centerX+25, baseY-15)
	
	// End effector
	display.SetPixel(centerX+25, baseY-15, true)
	display.SetPixel(centerX+26, baseY-15, true)
}

// Interface methods to maintain compatibility
func (ad *AdvancedRobotDisplay) ShowStartupScreen() {
	ad.manager.ShowStartupScreen()
}

func (ad *AdvancedRobotDisplay) ShowError(errorMsg string) {
	ad.manager.ShowError(errorMsg)
}

func (ad *AdvancedRobotDisplay) ShowCalibration(jointID string, progress float64) {
	ad.manager.ShowCalibration(jointID, progress)
}
