package display

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/font/basicfont"
)

// ArmVisualization provides 2D arm visualization on the display
type ArmVisualization struct {
	display      DisplayInterface
	centerX      int
	centerY      int
	scale        float64
	joints       []JointPosition
	lastUpdate   time.Time
}

// JointPosition represents a joint's position and angle
type JointPosition struct {
	X, Y  int
	Angle float64
	Length int
}

// NewArmVisualization creates a new arm visualization
func NewArmVisualization(display DisplayInterface) *ArmVisualization {
	return &ArmVisualization{
		display: display,
		centerX: 64,  // Center of 128px display
		centerY: 32,  // Center of 64px display
		scale:   0.5,
		joints:  make([]JointPosition, 6),
	}
}

// UpdateFromRobotStatus updates the visualization from robot status
func (av *ArmVisualization) UpdateFromRobotStatus(status RobotStatus) {
	av.lastUpdate = time.Now()
	
	// Convert joint positions to angles and update visualization
	for i, joint := range status.Joints {
		if i >= len(av.joints) {
			break
		}
		
		// Convert steps to radians
		angle := av.stepsToRadians(joint.Position)
		av.joints[i].Angle = angle
		
		// Set arm segment lengths (proportional)
		lengths := []int{20, 18, 15, 12, 10, 8} // Base to tip
		if i < len(lengths) {
			av.joints[i].Length = lengths[i]
		}
	}
}

// DrawArmVisualization draws the complete arm visualization
func (av *ArmVisualization) DrawArmVisualization() {
	av.display.Clear()
	
	// Header
	av.display.DrawText("ARM VIEW", 45, 10, basicfont.Face7x13)
	
	// Draw coordinate system
	av.drawCoordinateSystem()
	
	// Calculate forward kinematics
	av.calculateJointPositions()
	
	// Draw arm segments
	av.drawArmSegments()
	
	// Draw joint angles as text
	av.drawJointAngles()
	
	av.display.Update()
}

// drawCoordinateSystem draws X-Y coordinate reference
func (av *ArmVisualization) drawCoordinateSystem() {
	// Draw base circle
	av.drawCircle(av.centerX, av.centerY, 3, true)
	
	// Draw X and Y axes (small)
	av.display.DrawLine(av.centerX-8, av.centerY, av.centerX+8, av.centerY) // X-axis
	av.display.DrawLine(av.centerX, av.centerY-6, av.centerX, av.centerY+6) // Y-axis
	
	// Label axes
	av.display.DrawText("X", av.centerX+10, av.centerY+5, basicfont.Face7x13)
	av.display.DrawText("Y", av.centerX-5, av.centerY-8, basicfont.Face7x13)
}

// calculateJointPositions performs forward kinematics
func (av *ArmVisualization) calculateJointPositions() {
	currentX := float64(av.centerX)
	currentY := float64(av.centerY)
	cumulativeAngle := 0.0
	
	for i := range av.joints {
		if i >= 2 { // Only calculate for first few joints for 2D view
			break
		}
		
		cumulativeAngle += av.joints[i].Angle
		
		// Calculate joint position
		length := float64(av.joints[i].Length) * av.scale
		av.joints[i].X = int(currentX + length*math.Cos(cumulativeAngle))
		av.joints[i].Y = int(currentY + length*math.Sin(cumulativeAngle))
		
		currentX = float64(av.joints[i].X)
		currentY = float64(av.joints[i].Y)
	}
}

// drawArmSegments draws the arm segments connecting joints
func (av *ArmVisualization) drawArmSegments() {
	prevX := av.centerX
	prevY := av.centerY
	
	for i := range av.joints {
		if i >= 2 { // Only draw first few segments
			break
		}
		
		// Draw link
		av.display.DrawLine(prevX, prevY, av.joints[i].X, av.joints[i].Y)
		
		// Draw joint circle
		av.drawCircle(av.joints[i].X, av.joints[i].Y, 2, false)
		
		prevX = av.joints[i].X
		prevY = av.joints[i].Y
	}
	
	// Draw end effector
	if len(av.joints) > 1 {
		av.drawCircle(av.joints[1].X, av.joints[1].Y, 3, true)
	}
}

// drawJointAngles displays angle values as text
func (av *ArmVisualization) drawJointAngles() {
	y := 50
	for i := 0; i < 2 && i < len(av.joints); i++ {
		angle := av.joints[i].Angle * 180.0 / math.Pi // Convert to degrees
		text := fmt.Sprintf("J%d:%.0f°", i+1, angle)
		av.display.DrawText(text, 2, y+(i*10), basicfont.Face7x13)
	}
}

// Helper method to draw circles
func (av *ArmVisualization) drawCircle(centerX, centerY, radius int, filled bool) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				if filled || x*x+y*y >= (radius-1)*(radius-1) {
					av.display.SetPixel(centerX+x, centerY+y, true)
				}
			}
		}
	}
}

// stepsToRadians converts motor steps to radians
func (av *ArmVisualization) stepsToRadians(steps int) float64 {
	stepsPerRevolution := 200.0 * 32.0 // 200 steps * 32 microstepping
	return float64(steps) * 2.0 * math.Pi / stepsPerRevolution
}

// WorkspaceVisualization shows the reachable workspace
type WorkspaceVisualization struct {
	display DisplayInterface
}

// NewWorkspaceVisualization creates workspace visualization
func NewWorkspaceVisualization(display DisplayInterface) *WorkspaceVisualization {
	return &WorkspaceVisualization{display: display}
}

// DrawWorkspace shows the robot's reachable workspace
func (wv *WorkspaceVisualization) DrawWorkspace() {
	wv.display.Clear()
	
	// Header
	wv.display.DrawText("WORKSPACE", 35, 10, basicfont.Face7x13)
	
	centerX := 64
	centerY := 35
	
	// Draw workspace boundaries (circles representing reach)
	wv.drawCircle(centerX, centerY, 25, false) // Max reach
	wv.drawCircle(centerX, centerY, 15, false) // Min reach
	wv.drawCircle(centerX, centerY, 2, true)   // Base
	
	// Draw grid
	for i := centerX - 20; i <= centerX + 20; i += 10 {
		wv.display.DrawLine(i, centerY-20, i, centerY+20)
	}
	for i := centerY - 20; i <= centerY + 20; i += 10 {
		wv.display.DrawLine(centerX-20, i, centerX+20, i)
	}
	
	// Label dimensions
	wv.display.DrawText("25cm", centerX+15, centerY-20, basicfont.Face7x13)
	wv.display.DrawText("15cm", centerX+8, centerY-12, basicfont.Face7x13)
	
	wv.display.Update()
}

// Helper method for workspace visualization
func (wv *WorkspaceVisualization) drawCircle(centerX, centerY, radius int, filled bool) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				if filled || x*x+y*y >= (radius-1)*(radius-1) {
					wv.display.SetPixel(centerX+x, centerY+y, true)
				}
			}
		}
	}
}
