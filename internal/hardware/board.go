package hardware

import (
	"fmt"
	"log"

	"github.com/gemblerz/arm/pkg/stepper"
)

// Board represents a hardware board interface
type Board interface {
	Initialize() error
	CreateStepperMotor(config stepper.Config) (stepper.StepperMotor, error)
	Cleanup() error
	GetName() string
}

// MockBoard is a mock implementation for testing and development
type MockBoard struct {
	name        string
	initialized bool
	logOutput   bool
}

// NewMockBoard creates a new mock board
func NewMockBoard(name string, logOutput bool) *MockBoard {
	return &MockBoard{
		name:      name,
		logOutput: logOutput,
	}
}

// Initialize initializes the mock board
func (b *MockBoard) Initialize() error {
	if b.logOutput {
		log.Printf("[MOCK BOARD] Initializing %s", b.name)
	}
	b.initialized = true
	return nil
}

// CreateStepperMotor creates a new mock stepper motor
func (b *MockBoard) CreateStepperMotor(config stepper.Config) (stepper.StepperMotor, error) {
	if !b.initialized {
		return nil, fmt.Errorf("board %s not initialized", b.name)
	}
	
	if b.logOutput {
		log.Printf("[MOCK BOARD] Creating stepper motor with config: Step=%d, Dir=%d, Enable=%d",
			config.StepPin, config.DirPin, config.EnablePin)
	}
	
	return stepper.NewMockMotor(config, b.logOutput), nil
}

// Cleanup performs cleanup operations
func (b *MockBoard) Cleanup() error {
	if b.logOutput {
		log.Printf("[MOCK BOARD] Cleaning up %s", b.name)
	}
	b.initialized = false
	return nil
}

// GetName returns the board name
func (b *MockBoard) GetName() string {
	return b.name
}

// CoralBoard represents the Google Coral development board
type CoralBoard struct {
	name        string
	initialized bool
	pins        map[int]interface{} // Store pin references
}

// NewCoralBoard creates a new Coral board instance
func NewCoralBoard() *CoralBoard {
	return &CoralBoard{
		name: "Google Coral Dev Board",
		pins: make(map[int]interface{}),
	}
}

// Initialize initializes the Coral board
func (b *CoralBoard) Initialize() error {
	// In a real implementation with periph.io:
	// _, err := host.Init()
	// if err != nil {
	//     return fmt.Errorf("failed to initialize periph.io: %w", err)
	// }
	
	log.Printf("Initializing %s with Coral-specific GPIO", b.name)
	
	// Coral Dev Board specific initialization
	// The Coral uses specific GPIO mappings different from Raspberry Pi
	log.Printf("GPIO Controller: Coral Edge TPU")
	log.Printf("Available GPIO pins: 40-pin header compatible")
	
	b.initialized = true
	return nil
}

// CreateStepperMotor creates a stepper motor for the Coral board
func (b *CoralBoard) CreateStepperMotor(config stepper.Config) (stepper.StepperMotor, error) {
	if !b.initialized {
		return nil, fmt.Errorf("board %s not initialized", b.name)
	}
	
	// Validate Coral-specific GPIO pins
	if err := b.validateCoralGPIOPin(config.StepPin); err != nil {
		return nil, fmt.Errorf("invalid step pin: %w", err)
	}
	if err := b.validateCoralGPIOPin(config.DirPin); err != nil {
		return nil, fmt.Errorf("invalid direction pin: %w", err)
	}
	if config.EnablePin >= 0 {
		if err := b.validateCoralGPIOPin(config.EnablePin); err != nil {
			return nil, fmt.Errorf("invalid enable pin: %w", err)
		}
	}
	
	log.Printf("Creating Coral GPIO stepper motor: Step=%d, Dir=%d, Enable=%d",
		config.StepPin, config.DirPin, config.EnablePin)
	
	// In a real implementation, this would create actual GPIO pins for Coral:
	// stepPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.StepPin))
	// dirPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.DirPin))
	// var enablePin gpio.PinOut
	// if config.EnablePin >= 0 {
	//     enablePin = gpioreg.ByName(fmt.Sprintf("GPIO%d", config.EnablePin))
	// }
	
	// For development, return enhanced mock motor with Coral-specific features
	motor := stepper.NewMockMotor(config, true)
	return &CoralStepperMotor{
		MockMotor: motor,
		board:     b,
		config:    config,
	}, nil
}

// validateCoralGPIOPin validates GPIO pin numbers for Coral Dev Board
func (b *CoralBoard) validateCoralGPIOPin(pin int) error {
	// Coral Dev Board GPIO pin validation
	// The Coral has specific GPIO pins available on the 40-pin header
	validPins := map[int]bool{
		// Coral Dev Board GPIO pins (subset that are safe to use)
		138: true, 140: true, 141: true, 142: true, 143: true,
		144: true, 145: true, 146: true, 147: true, 148: true,
		149: true, 150: true, 151: true, 152: true, 153: true,
		154: true, 155: true, 156: true, 157: true, 158: true,
	}
	
	if !validPins[pin] {
		return fmt.Errorf("GPIO pin %d not available on Coral Dev Board", pin)
	}
	
	// Check for pins that might conflict with system functions
	systemPins := map[int]string{
		138: "SPI0_MOSI", 140: "SPI0_MISO", 141: "SPI0_SCLK",
	}
	
	if sysFunc, exists := systemPins[pin]; exists {
		log.Printf("Warning: GPIO pin %d is used for %s on Coral", pin, sysFunc)
	}
	
	return nil
}

// Cleanup performs cleanup operations
func (b *CoralBoard) Cleanup() error {
	log.Printf("Cleaning up %s GPIO resources", b.name)
	
	// In a real implementation, cleanup GPIO pins:
	// for _, pin := range b.pins {
	//     if gpioPin, ok := pin.(gpio.PinOut); ok {
	//         gpioPin.Out(gpio.Low)
	//     }
	// }
	
	b.initialized = false
	b.pins = make(map[int]interface{})
	return nil
}

// GetName returns the board name
func (b *CoralBoard) GetName() string {
	return b.name
}

// CoralStepperMotor wraps a mock motor with Coral-specific functionality
type CoralStepperMotor struct {
	*stepper.MockMotor
	board  *CoralBoard
	config stepper.Config
}

// Enable overrides the mock enable to simulate Coral GPIO control
func (m *CoralStepperMotor) Enable() error {
	log.Printf("[CORAL GPIO] Enabling motor on pins: Step=%d, Dir=%d, Enable=%d",
		m.config.StepPin, m.config.DirPin, m.config.EnablePin)
	
	// In a real implementation for Coral:
	// if m.enablePin != nil {
	//     m.enablePin.Out(gpio.High)
	// }
	// Configure pin drive strength and slew rate for stepper motors
	
	return m.MockMotor.Enable()
}

// Disable overrides the mock disable to simulate Coral GPIO control
func (m *CoralStepperMotor) Disable() error {
	log.Printf("[CORAL GPIO] Disabling motor on pins: Step=%d, Dir=%d, Enable=%d",
		m.config.StepPin, m.config.DirPin, m.config.EnablePin)
	
	// In a real implementation for Coral:
	// if m.enablePin != nil {
	//     m.enablePin.Out(gpio.Low)
	// }
	
	return m.MockMotor.Disable()
}

// LinuxGPIOBoard represents a generic Linux board with GPIO support (Raspberry Pi, Coral, etc.)
type LinuxGPIOBoard struct {
	name        string
	initialized bool
	pins        map[int]interface{} // Store pin references
}

// NewLinuxGPIOBoard creates a new Linux GPIO board instance
func NewLinuxGPIOBoard(name string) *LinuxGPIOBoard {
	return &LinuxGPIOBoard{
		name: name,
		pins: make(map[int]interface{}),
	}
}

// Initialize initializes the Linux GPIO board using periph.io
func (b *LinuxGPIOBoard) Initialize() error {
	// In a real implementation with periph.io:
	// _, err := host.Init()
	// if err != nil {
	//     return fmt.Errorf("failed to initialize periph.io: %w", err)
	// }
	
	// For now, simulate initialization
	log.Printf("Initializing %s with GPIO support", b.name)
	b.initialized = true
	return nil
}

// CreateStepperMotor creates a GPIO-backed stepper motor
func (b *LinuxGPIOBoard) CreateStepperMotor(config stepper.Config) (stepper.StepperMotor, error) {
	if !b.initialized {
		return nil, fmt.Errorf("board %s not initialized", b.name)
	}
	
	// Validate GPIO pin numbers (typical range for Raspberry Pi)
	if err := b.validateGPIOPin(config.StepPin); err != nil {
		return nil, fmt.Errorf("invalid step pin: %w", err)
	}
	if err := b.validateGPIOPin(config.DirPin); err != nil {
		return nil, fmt.Errorf("invalid direction pin: %w", err)
	}
	if config.EnablePin >= 0 {
		if err := b.validateGPIOPin(config.EnablePin); err != nil {
			return nil, fmt.Errorf("invalid enable pin: %w", err)
		}
	}
	
	log.Printf("Creating GPIO stepper motor on %s: Step=%d, Dir=%d, Enable=%d",
		b.name, config.StepPin, config.DirPin, config.EnablePin)
	
	// In a real implementation, this would create actual GPIO pins:
	// stepPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.StepPin))
	// dirPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.DirPin))
	// var enablePin gpio.PinOut
	// if config.EnablePin >= 0 {
	//     enablePin = gpioreg.ByName(fmt.Sprintf("GPIO%d", config.EnablePin))
	// }
	
	// For development, return mock motor with GPIO simulation
	motor := stepper.NewMockMotor(config, true)
	return &GPIOStepperMotor{
		MockMotor: motor,
		board:     b,
		config:    config,
	}, nil
}

// validateGPIOPin validates GPIO pin numbers for common Linux boards
func (b *LinuxGPIOBoard) validateGPIOPin(pin int) error {
	// Raspberry Pi GPIO pin validation (BCM numbering)
	// Valid pins: 0-27 (depending on model)
	if pin < 0 || pin > 27 {
		return fmt.Errorf("GPIO pin %d out of range (0-27)", pin)
	}
	
	// Check for reserved pins that shouldn't be used
	reservedPins := []int{0, 1, 14, 15} // I2C and UART pins
	for _, reserved := range reservedPins {
		if pin == reserved {
			log.Printf("Warning: GPIO pin %d is typically reserved for system use", pin)
		}
	}
	
	return nil
}

// Cleanup performs cleanup operations
func (b *LinuxGPIOBoard) Cleanup() error {
	log.Printf("Cleaning up %s GPIO resources", b.name)
	
	// In a real implementation, cleanup GPIO pins:
	// for _, pin := range b.pins {
	//     if gpioPin, ok := pin.(gpio.PinOut); ok {
	//         gpioPin.Out(gpio.Low)
	//     }
	// }
	
	b.initialized = false
	b.pins = make(map[int]interface{})
	return nil
}

// GetName returns the board name
func (b *LinuxGPIOBoard) GetName() string {
	return b.name
}

// GPIOStepperMotor wraps a mock motor with GPIO-specific functionality
type GPIOStepperMotor struct {
	*stepper.MockMotor
	board  *LinuxGPIOBoard
	config stepper.Config
}

// Enable overrides the mock enable to simulate GPIO control
func (m *GPIOStepperMotor) Enable() error {
	log.Printf("[GPIO] Enabling motor on pins: Step=%d, Dir=%d, Enable=%d",
		m.config.StepPin, m.config.DirPin, m.config.EnablePin)
	
	// In a real implementation:
	// if m.enablePin != nil {
	//     m.enablePin.Out(gpio.High)
	// }
	
	return m.MockMotor.Enable()
}

// Disable overrides the mock disable to simulate GPIO control
func (m *GPIOStepperMotor) Disable() error {
	log.Printf("[GPIO] Disabling motor on pins: Step=%d, Dir=%d, Enable=%d",
		m.config.StepPin, m.config.DirPin, m.config.EnablePin)
	
	// In a real implementation:
	// if m.enablePin != nil {
	//     m.enablePin.Out(gpio.Low)
	// }
	
	return m.MockMotor.Disable()
}

// BoardFactory creates boards based on type
type BoardFactory struct{}

// CreateBoard creates a board instance based on the board type
func (f *BoardFactory) CreateBoard(boardType string, logOutput bool) (Board, error) {
	switch boardType {
	case "mock":
		return NewMockBoard("Mock Development Board", logOutput), nil
	case "coral":
		return NewCoralBoard(), nil
	case "coral-extended":
		// Create Coral board with AW9523 GPIO extender
		extender := NewAW9523GPIOExtender(0x58) // Default I2C address
		return NewExtendedCoralBoard(extender), nil
	case "raspberry-pi", "rpi":
		return NewLinuxGPIOBoard("Raspberry Pi"), nil
	case "linux-gpio":
		return NewLinuxGPIOBoard("Generic Linux GPIO Board"), nil
	default:
		return nil, fmt.Errorf("unknown board type: %s. Supported types: mock, coral, coral-extended, raspberry-pi, linux-gpio", boardType)
	}
}
