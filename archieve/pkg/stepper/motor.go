package stepper

import (
	"time"
)

// Direction represents the rotation direction of the stepper motor
type Direction int

const (
	Clockwise Direction = iota
	CounterClockwise
)

// StepperMotor interface defines the basic operations for a stepper motor
type StepperMotor interface {
	// Step moves the motor by the specified number of steps in the given direction
	Step(steps int, direction Direction) error
	
	// SetSpeed sets the stepping speed in steps per second
	SetSpeed(stepsPerSecond int) error
	
	// Enable enables the motor (powers the coils)
	Enable() error
	
	// Disable disables the motor (unpowers the coils)
	Disable() error
	
	// IsEnabled returns true if the motor is currently enabled
	IsEnabled() bool
	
	// GetPosition returns the current position in steps from the reference point
	GetPosition() int
	
	// SetPosition sets the current position (useful for homing)
	SetPosition(position int)
	
	// MoveTo moves the motor to an absolute position
	MoveTo(position int) error
}

// Config holds the configuration for a stepper motor
type Config struct {
	StepPin      int           // GPIO pin for step signal
	DirPin       int           // GPIO pin for direction signal
	EnablePin    int           // GPIO pin for enable signal (optional, -1 if not used)
	StepsPerRev  int           // Number of steps per full revolution
	MaxSpeed     int           // Maximum speed in steps per second
	DefaultSpeed int           // Default speed in steps per second
	StepDelay    time.Duration // Minimum delay between steps
}

// Motor represents a stepper motor with its configuration and state
type Motor struct {
	config   Config
	position int
	enabled  bool
	speed    int
	stepFunc func(bool) error  // Function to pulse the step pin
	dirFunc  func(Direction) error // Function to set direction
	enableFunc func(bool) error    // Function to enable/disable motor
}

// NewMotor creates a new stepper motor instance
func NewMotor(config Config) *Motor {
	return &Motor{
		config:   config,
		position: 0,
		enabled:  false,
		speed:    config.DefaultSpeed,
	}
}

// SetHardwareFunctions sets the hardware interface functions for the motor
func (m *Motor) SetHardwareFunctions(stepFunc func(bool) error, dirFunc func(Direction) error, enableFunc func(bool) error) {
	m.stepFunc = stepFunc
	m.dirFunc = dirFunc
	m.enableFunc = enableFunc
}

// Step implements the StepperMotor interface
func (m *Motor) Step(steps int, direction Direction) error {
	if !m.enabled {
		return ErrMotorNotEnabled
	}
	
	if m.dirFunc != nil {
		if err := m.dirFunc(direction); err != nil {
			return err
		}
	}
	
	stepDelay := time.Second / time.Duration(m.speed)
	if stepDelay < m.config.StepDelay {
		stepDelay = m.config.StepDelay
	}
	
	for i := 0; i < steps; i++ {
		if m.stepFunc != nil {
			// Pulse the step pin
			if err := m.stepFunc(true); err != nil {
				return err
			}
			time.Sleep(stepDelay / 2)
			if err := m.stepFunc(false); err != nil {
				return err
			}
			time.Sleep(stepDelay / 2)
		}
		
		// Update position
		if direction == Clockwise {
			m.position++
		} else {
			m.position--
		}
	}
	
	return nil
}

// SetSpeed implements the StepperMotor interface
func (m *Motor) SetSpeed(stepsPerSecond int) error {
	if stepsPerSecond <= 0 || stepsPerSecond > m.config.MaxSpeed {
		return ErrInvalidSpeed
	}
	m.speed = stepsPerSecond
	return nil
}

// Enable implements the StepperMotor interface
func (m *Motor) Enable() error {
	if m.enableFunc != nil {
		if err := m.enableFunc(true); err != nil {
			return err
		}
	}
	m.enabled = true
	return nil
}

// Disable implements the StepperMotor interface
func (m *Motor) Disable() error {
	if m.enableFunc != nil {
		if err := m.enableFunc(false); err != nil {
			return err
		}
	}
	m.enabled = false
	return nil
}

// IsEnabled implements the StepperMotor interface
func (m *Motor) IsEnabled() bool {
	return m.enabled
}

// GetPosition implements the StepperMotor interface
func (m *Motor) GetPosition() int {
	return m.position
}

// SetPosition implements the StepperMotor interface
func (m *Motor) SetPosition(position int) {
	m.position = position
}

// MoveTo implements the StepperMotor interface
func (m *Motor) MoveTo(position int) error {
	if !m.enabled {
		return ErrMotorNotEnabled
	}
	
	stepsDiff := position - m.position
	if stepsDiff == 0 {
		return nil
	}
	
	direction := Clockwise
	steps := stepsDiff
	if stepsDiff < 0 {
		direction = CounterClockwise
		steps = -stepsDiff
	}
	
	return m.Step(steps, direction)
}
