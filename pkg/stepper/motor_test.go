package stepper

import (
	"testing"
	"time"
)

func TestMockMotor(t *testing.T) {
	config := Config{
		StepPin:      2,
		DirPin:       3,
		EnablePin:    4,
		StepsPerRev:  200,
		MaxSpeed:     1000,
		DefaultSpeed: 200,
		StepDelay:    time.Microsecond * 1000,
	}

	motor := NewMockMotor(config, false)

	if motor.IsEnabled() {
		t.Error("Motor should be disabled initially")
	}

	if motor.GetPosition() != 0 {
		t.Error("Motor position should be 0 initially")
	}

	if err := motor.Enable(); err != nil {
		t.Fatalf("Failed to enable motor: %v", err)
	}

	if !motor.IsEnabled() {
		t.Error("Motor should be enabled after Enable()")
	}

	if err := motor.SetSpeed(500); err != nil {
		t.Fatalf("Failed to set speed: %v", err)
	}

	if err := motor.Step(10, Clockwise); err != nil {
		t.Fatalf("Failed to step motor: %v", err)
	}

	if motor.GetPosition() != 10 {
		t.Errorf("Expected position 10, got %d", motor.GetPosition())
	}

	if err := motor.Disable(); err != nil {
		t.Fatalf("Failed to disable motor: %v", err)
	}
}

func TestDriverConfigs(t *testing.T) {
	dm556tConfig := GetDM556TConfig()
	if dm556tConfig.Type != DriverDM556T {
		t.Error("DM556T config should have correct type")
	}

	dm320tConfig := GetDM320TConfig()
	if dm320tConfig.Type != DriverDM320T {
		t.Error("DM320T config should have correct type")
	}

	drv8825Config := GetDRV8825Config()
	if drv8825Config.Type != DriverDRV8825 {
		t.Error("DRV8825 config should have correct type")
	}
}
