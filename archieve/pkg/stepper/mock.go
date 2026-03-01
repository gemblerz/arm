package stepper

import (
	"fmt"
	"log"
	"time"
)

// MockMotor is a mock implementation of StepperMotor for testing and development
type MockMotor struct {
	*Motor
	logOutput bool
}

// NewMockMotor creates a new mock stepper motor
func NewMockMotor(config Config, logOutput bool) *MockMotor {
	motor := NewMotor(config)
	mock := &MockMotor{
		Motor:     motor,
		logOutput: logOutput,
	}
	
	// Set mock hardware functions
	motor.SetHardwareFunctions(
		mock.mockStepFunc,
		mock.mockDirFunc,
		mock.mockEnableFunc,
	)
	
	return mock
}

// mockStepFunc simulates stepping the motor
func (m *MockMotor) mockStepFunc(high bool) error {
	if m.logOutput {
		if high {
			log.Printf("[MOCK] Step pin HIGH (pin %d)", m.config.StepPin)
		} else {
			log.Printf("[MOCK] Step pin LOW (pin %d)", m.config.StepPin)
		}
	}
	return nil
}

// mockDirFunc simulates setting the direction
func (m *MockMotor) mockDirFunc(direction Direction) error {
	if m.logOutput {
		dirStr := "CLOCKWISE"
		if direction == CounterClockwise {
			dirStr = "COUNTER-CLOCKWISE"
		}
		log.Printf("[MOCK] Direction set to %s (pin %d)", dirStr, m.config.DirPin)
	}
	return nil
}

// mockEnableFunc simulates enabling/disabling the motor
func (m *MockMotor) mockEnableFunc(enable bool) error {
	if m.logOutput {
		state := "ENABLED"
		if !enable {
			state = "DISABLED"
		}
		log.Printf("[MOCK] Motor %s (pin %d)", state, m.config.EnablePin)
	}
	return nil
}

// StepWithDelay performs a step operation with visible delay for demonstration
func (m *MockMotor) StepWithDelay(steps int, direction Direction, delay time.Duration) error {
	if !m.enabled {
		return ErrMotorNotEnabled
	}
	
	if m.dirFunc != nil {
		if err := m.dirFunc(direction); err != nil {
			return err
		}
	}
	
	fmt.Printf("Stepping %d steps in direction %v...\n", steps, direction)
	
	for i := 0; i < steps; i++ {
		if m.stepFunc != nil {
			// Pulse the step pin
			if err := m.stepFunc(true); err != nil {
				return err
			}
			time.Sleep(delay / 2)
			if err := m.stepFunc(false); err != nil {
				return err
			}
			time.Sleep(delay / 2)
		}
		
		// Update position
		if direction == Clockwise {
			m.position++
		} else {
			m.position--
		}
		
		if m.logOutput {
			fmt.Printf("  Step %d/%d completed, position: %d\n", i+1, steps, m.position)
		}
	}
	
	return nil
}

// GetStatus returns a formatted status string
func (m *MockMotor) GetStatus() string {
	status := "DISABLED"
	if m.enabled {
		status = "ENABLED"
	}
	
	return fmt.Sprintf("Motor Status: %s, Position: %d, Speed: %d steps/sec", 
		status, m.position, m.speed)
}
