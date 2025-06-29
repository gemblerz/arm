package hardware

import (
	"fmt"
	"log"
	
	"github.com/gemblerz/arm/pkg/stepper"
)

// GPIOExtender represents a GPIO expansion chip
type GPIOExtender interface {
	Initialize() error
	SetPin(pin int, value bool) error
	GetPin(pin int) (bool, error)
	SetPinMode(pin int, mode PinMode) error
	GetPinCount() int
	Cleanup() error
}

// PinMode represents GPIO pin modes
type PinMode int

const (
	PinModeInput PinMode = iota
	PinModeOutput
	PinModeInputPullup
	PinModeInputPulldown
)

// AW9523GPIOExtender implements GPIO extender for Adafruit AW9523
type AW9523GPIOExtender struct {
	address     uint8
	initialized bool
	pinStates   map[int]bool
	pinModes    map[int]PinMode
}

// NewAW9523GPIOExtender creates a new AW9523 GPIO extender
func NewAW9523GPIOExtender(i2cAddress uint8) *AW9523GPIOExtender {
	return &AW9523GPIOExtender{
		address:   i2cAddress,
		pinStates: make(map[int]bool),
		pinModes:  make(map[int]PinMode),
	}
}

// Initialize configures the AW9523 GPIO extender
func (e *AW9523GPIOExtender) Initialize() error {
	// In real implementation:
	// - Initialize I2C communication
	// - Configure AW9523 registers
	// - Set default pin modes
	
	log.Printf("Initializing AW9523 GPIO Extender at I2C address 0x%02X", e.address)
	
	// AW9523 has 16 GPIO pins (P0_0 to P1_7)
	for pin := 0; pin < 16; pin++ {
		e.pinModes[pin] = PinModeOutput
		e.pinStates[pin] = false
	}
	
	e.initialized = true
	return nil
}

// SetPin sets the state of a GPIO pin
func (e *AW9523GPIOExtender) SetPin(pin int, value bool) error {
	if !e.initialized {
		return fmt.Errorf("GPIO extender not initialized")
	}
	
	if pin < 0 || pin >= 16 {
		return fmt.Errorf("pin %d out of range (0-15)", pin)
	}
	
	// In real implementation:
	// Write to AW9523 output register via I2C
	
	e.pinStates[pin] = value
	log.Printf("[AW9523] Set pin %d to %t", pin, value)
	
	return nil
}

// GetPin reads the state of a GPIO pin
func (e *AW9523GPIOExtender) GetPin(pin int) (bool, error) {
	if !e.initialized {
		return false, fmt.Errorf("GPIO extender not initialized")
	}
	
	if pin < 0 || pin >= 16 {
		return false, fmt.Errorf("pin %d out of range (0-15)", pin)
	}
	
	// In real implementation:
	// Read from AW9523 input register via I2C
	
	return e.pinStates[pin], nil
}

// SetPinMode configures pin mode (input/output)
func (e *AW9523GPIOExtender) SetPinMode(pin int, mode PinMode) error {
	if !e.initialized {
		return fmt.Errorf("GPIO extender not initialized")
	}
	
	if pin < 0 || pin >= 16 {
		return fmt.Errorf("pin %d out of range (0-15)", pin)
	}
	
	e.pinModes[pin] = mode
	log.Printf("[AW9523] Set pin %d mode to %v", pin, mode)
	
	return nil
}

// GetPinCount returns the number of available GPIO pins
func (e *AW9523GPIOExtender) GetPinCount() int {
	return 16
}

// Cleanup performs cleanup operations
func (e *AW9523GPIOExtender) Cleanup() error {
	// Set all pins to safe state
	for pin := 0; pin < 16; pin++ {
		e.SetPin(pin, false)
	}
	
	e.initialized = false
	return nil
}

// ExtendedCoralBoard represents Coral board with GPIO extender
type ExtendedCoralBoard struct {
	*CoralBoard
	extender GPIOExtender
	extenderPinOffset int // Offset for extender pins in addressing
}

// NewExtendedCoralBoard creates Coral board with GPIO extender
func NewExtendedCoralBoard(extender GPIOExtender) *ExtendedCoralBoard {
	return &ExtendedCoralBoard{
		CoralBoard:        NewCoralBoard(),
		extender:          extender,
		extenderPinOffset: 1000, // Use 1000+ for extender pins
	}
}

// Initialize initializes both Coral board and GPIO extender
func (b *ExtendedCoralBoard) Initialize() error {
	// Initialize base Coral board
	if err := b.CoralBoard.Initialize(); err != nil {
		return err
	}
	
	// Initialize GPIO extender
	if err := b.extender.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize GPIO extender: %w", err)
	}
	
	log.Printf("Extended Coral board initialized with GPIO extender")
	return nil
}

// CreateStepperMotor creates motor with mixed Coral/extender GPIO
func (b *ExtendedCoralBoard) CreateStepperMotor(config stepper.Config) (stepper.StepperMotor, error) {
	if !b.initialized {
		return nil, fmt.Errorf("board not initialized")
	}
	
	// Determine which pins are on extender vs native Coral
	stepOnExtender := b.isExtenderPin(config.StepPin)
	dirOnExtender := b.isExtenderPin(config.DirPin)
	enableOnExtender := config.EnablePin >= 0 && b.isExtenderPin(config.EnablePin)
	
	log.Printf("Creating motor: Step=%d(%s), Dir=%d(%s), Enable=%d(%s)",
		config.StepPin, b.pinLocation(stepOnExtender),
		config.DirPin, b.pinLocation(dirOnExtender),
		config.EnablePin, b.pinLocation(enableOnExtender))
	
	// Validate pins
	if err := b.validateMixedPins(config); err != nil {
		return nil, err
	}
	
	// Create enhanced motor with extender support
	motor := stepper.NewMockMotor(config, true)
	return &ExtendedStepperMotor{
		MockMotor: motor,
		board:     b,
		config:    config,
	}, nil
}

// isExtenderPin checks if pin number refers to extender
func (b *ExtendedCoralBoard) isExtenderPin(pin int) bool {
	return pin >= b.extenderPinOffset
}

// pinLocation returns human-readable pin location
func (b *ExtendedCoralBoard) pinLocation(onExtender bool) string {
	if onExtender {
		return "Extender"
	}
	return "Coral"
}

// validateMixedPins validates pin assignments across Coral and extender
func (b *ExtendedCoralBoard) validateMixedPins(config stepper.Config) error {
	// Validate Coral pins
	if !b.isExtenderPin(config.StepPin) {
		if err := b.CoralBoard.validateCoralGPIOPin(config.StepPin); err != nil {
			return fmt.Errorf("step pin validation failed: %w", err)
		}
	}
	
	if !b.isExtenderPin(config.DirPin) {
		if err := b.CoralBoard.validateCoralGPIOPin(config.DirPin); err != nil {
			return fmt.Errorf("direction pin validation failed: %w", err)
		}
	}
	
	if config.EnablePin >= 0 && !b.isExtenderPin(config.EnablePin) {
		if err := b.CoralBoard.validateCoralGPIOPin(config.EnablePin); err != nil {
			return fmt.Errorf("enable pin validation failed: %w", err)
		}
	}
	
	// Validate extender pins
	if b.isExtenderPin(config.StepPin) {
		extPin := config.StepPin - b.extenderPinOffset
		if extPin < 0 || extPin >= b.extender.GetPinCount() {
			return fmt.Errorf("step pin %d invalid for extender", extPin)
		}
	}
	
	// Similar validation for other pins...
	
	return nil
}

// Cleanup cleans up both Coral and extender
func (b *ExtendedCoralBoard) Cleanup() error {
	// Cleanup extender first
	if err := b.extender.Cleanup(); err != nil {
		log.Printf("Warning: extender cleanup failed: %v", err)
	}
	
	// Cleanup base Coral board
	return b.CoralBoard.Cleanup()
}

// ExtendedStepperMotor wraps motor with extender pin support
type ExtendedStepperMotor struct {
	*stepper.MockMotor
	board  *ExtendedCoralBoard
	config stepper.Config
}

// Enable enables motor using mixed GPIO
func (m *ExtendedStepperMotor) Enable() error {
	if m.config.EnablePin >= 0 {
		if m.board.isExtenderPin(m.config.EnablePin) {
			extPin := m.config.EnablePin - m.board.extenderPinOffset
			if err := m.board.extender.SetPin(extPin, true); err != nil {
				return fmt.Errorf("failed to enable via extender: %w", err)
			}
		} else {
			// Use native Coral GPIO
			log.Printf("[CORAL] Enable motor on native pin %d", m.config.EnablePin)
		}
	}
	
	return m.MockMotor.Enable()
}

// Disable disables motor using mixed GPIO  
func (m *ExtendedStepperMotor) Disable() error {
	if m.config.EnablePin >= 0 {
		if m.board.isExtenderPin(m.config.EnablePin) {
			extPin := m.config.EnablePin - m.board.extenderPinOffset
			if err := m.board.extender.SetPin(extPin, false); err != nil {
				return fmt.Errorf("failed to disable via extender: %w", err)
			}
		} else {
			// Use native Coral GPIO
			log.Printf("[CORAL] Disable motor on native pin %d", m.config.EnablePin)
		}
	}
	
	return m.MockMotor.Disable()
}
