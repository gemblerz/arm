package robot

import (
	"time"

	"github.com/gemblerz/arm/pkg/display"
)

// DisplayAdapter adapts between robot and display interfaces
type DisplayAdapter struct {
	manager *display.RobotDisplayManager
}

// NewDisplayAdapter creates a new display adapter
func NewDisplayAdapter(displayInterface display.DisplayInterface, updateRate time.Duration) *DisplayAdapter {
	manager := display.NewRobotDisplayManager(displayInterface, updateRate)
	return &DisplayAdapter{
		manager: manager,
	}
}

// UpdateStatus updates the robot status for display
func (da *DisplayAdapter) UpdateStatus(status DisplayStatus) {
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

	da.manager.UpdateStatus(robotStatus)
}

// ShowStartupScreen displays the startup screen
func (da *DisplayAdapter) ShowStartupScreen() {
	da.manager.ShowStartupScreen()
}

// ShowError displays an error message
func (da *DisplayAdapter) ShowError(errorMsg string) {
	da.manager.ShowError(errorMsg)
}

// ShowCalibration displays calibration progress
func (da *DisplayAdapter) ShowCalibration(jointID string, progress float64) {
	da.manager.ShowCalibration(jointID, progress)
}

// Start starts the display update loop
func (da *DisplayAdapter) Start() {
	da.manager.Start()
}

// Stop stops the display update loop
func (da *DisplayAdapter) Stop() {
	da.manager.Stop()
}
