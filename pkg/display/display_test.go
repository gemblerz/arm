package display

import (
	"testing"
	"time"
)

func TestSSD1306Creation(t *testing.T) {
	config := DisplayConfig{
		Width:      128,
		Height:     64,
		I2CAddress: 0x3C,
		I2CBus:     1,
		Mock:       true,
		Verbose:    false,
	}

	display, err := NewSSD1306(config)
	if err != nil {
		t.Fatalf("Failed to create SSD1306: %v", err)
	}

	if display.Width != 128 {
		t.Errorf("Expected width 128, got %d", display.Width)
	}

	if display.Height != 64 {
		t.Errorf("Expected height 64, got %d", display.Height)
	}

	if display.I2CAddress != 0x3C {
		t.Errorf("Expected I2C address 0x3C, got 0x%02X", display.I2CAddress)
	}
}

func TestSSD1306BasicOperations(t *testing.T) {
	config := DisplayConfig{
		Width:      128,
		Height:     64,
		I2CAddress: 0x3C,
		I2CBus:     1,
		Mock:       true,
		Verbose:    false,
	}

	display, err := NewSSD1306(config)
	if err != nil {
		t.Fatalf("Failed to create SSD1306: %v", err)
	}

	// Test clear
	display.Clear()

	// Test pixel setting
	display.SetPixel(10, 10, true)
	display.SetPixel(20, 20, false)

	// Test text drawing
	display.DrawText("Test", 0, 10, nil)

	// Test line drawing
	display.DrawLine(0, 0, 10, 10)

	// Test rectangle drawing
	display.DrawRectangle(50, 50, 20, 10, false)
	display.DrawRectangle(75, 50, 20, 10, true)

	// Test progress bar
	display.DrawProgressBar(10, 30, 50, 8, 0.75)

	// Test update
	err = display.Update()
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Test brightness
	err = display.SetBrightness(128)
	if err != nil {
		t.Errorf("SetBrightness failed: %v", err)
	}

	// Test close
	err = display.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestRobotDisplayManager(t *testing.T) {
	config := DisplayConfig{
		Width:      128,
		Height:     64,
		I2CAddress: 0x3C,
		I2CBus:     1,
		Mock:       true,
		Verbose:    false,
	}

	display, err := NewSSD1306(config)
	if err != nil {
		t.Fatalf("Failed to create SSD1306: %v", err)
	}

	dm := NewRobotDisplayManager(display, 100*time.Millisecond)

	// Test status update
	status := RobotStatus{
		IsHomed:   true,
		IsEnabled: true,
		IsMoving:  false,
		CurrentOp: "Testing",
		Joints: []JointStatus{
			{ID: "pan", Position: 100, Target: 100, MinPos: -1000, MaxPos: 1000, Enabled: true},
			{ID: "tilt", Position: -50, Target: 0, MinPos: -500, MaxPos: 500, Enabled: true},
			{ID: "joint3", Position: 200, Target: 300, MinPos: 0, MaxPos: 500, Enabled: true},
		},
	}

	dm.UpdateStatus(status)

	// Test startup screen
	dm.ShowStartupScreen()

	// Test error display
	dm.ShowError("Test error message that is quite long and should wrap")

	// Test calibration display
	dm.ShowCalibration("pan", 0.75)

	// Start and stop display manager
	dm.Start()
	time.Sleep(200 * time.Millisecond)
	dm.Stop()

	display.Close()
}

func TestDisplayConfigurations(t *testing.T) {
	configurations := []DisplayConfig{
		{Width: 128, Height: 64, I2CAddress: 0x3C, I2CBus: 1, Mock: true},
		{Width: 128, Height: 32, I2CAddress: 0x3D, I2CBus: 1, Mock: true},
		{Width: 64, Height: 48, I2CAddress: 0x3C, I2CBus: 0, Mock: true},
	}

	for i, config := range configurations {
		display, err := NewSSD1306(config)
		if err != nil {
			t.Errorf("Configuration %d failed: %v", i, err)
			continue
		}

		// Test basic operations on each configuration
		display.Clear()
		display.DrawText("Test", 0, 10, nil)
		display.Update()
		display.Close()
	}
}
