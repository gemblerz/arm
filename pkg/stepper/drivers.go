package stepper

import (
	"fmt"
	"time"
)

// DriverType represents different stepper motor driver types
type DriverType string

const (
	DriverDRV8825 DriverType = "DRV8825"
	DriverDM320T  DriverType = "DM320T" 
	DriverDM556T  DriverType = "DM556T"
	DriverUnknown DriverType = "UNKNOWN"
)

// DriverConfig holds driver-specific configuration
type DriverConfig struct {
	Type          DriverType
	Microstepping int           // 1, 2, 4, 8, 16, 32
	StepMode      StepMode      // Full, Half, Quarter, etc.
	CurrentLimit  float32       // Current limit in Amperes
	DecayMode     DecayMode     // Fast, Slow, Mixed
	EnableLogic   bool          // true = active high, false = active low
	StepPulseWidth time.Duration // Minimum step pulse width
}

// StepMode represents microstepping modes
type StepMode int

const (
	FullStep StepMode = iota
	HalfStep
	QuarterStep
	EighthStep
	SixteenthStep
	ThirtySecondStep
)

// DecayMode represents current decay modes
type DecayMode int

const (
	FastDecay DecayMode = iota
	SlowDecay
	MixedDecay
	AutoDecay
)

// StepperDriver interface defines driver-specific operations
type StepperDriver interface {
	// Initialize configures the driver with specific settings
	Initialize(config DriverConfig) error
	
	// SetMicrostepping configures microstepping mode
	SetMicrostepping(mode StepMode) error
	
	// SetCurrentLimit sets the current limit (if supported)
	SetCurrentLimit(current float32) error
	
	// GetDriverInfo returns driver capabilities and current settings
	GetDriverInfo() DriverInfo
	
	// ValidateConfig checks if configuration is valid for this driver
	ValidateConfig(config DriverConfig) error
}

// DriverInfo contains driver capabilities and status
type DriverInfo struct {
	Type              DriverType
	SupportedModes    []StepMode
	MaxCurrent        float32
	MinStepPulseWidth time.Duration
	SupportsCurrentControl bool
	SupportsDecayControl   bool
	RequiresModeSet        bool // Some drivers need mode pins set
}

// GetDRV8825Config returns optimal configuration for DRV8825
func GetDRV8825Config() DriverConfig {
	return DriverConfig{
		Type:           DriverDRV8825,
		Microstepping:  16,                        // 1/16 step for smooth movement
		StepMode:       SixteenthStep,
		CurrentLimit:   1.5,                       // 1.5A typical for NEMA 17
		DecayMode:      MixedDecay,
		EnableLogic:    false,                     // Active low enable
		StepPulseWidth: time.Microsecond * 2,      // 2µs minimum
	}
}

// GetDM320TConfig returns optimal configuration for DM320T
func GetDM320TConfig() DriverConfig {
	return DriverConfig{
		Type:           DriverDM320T,
		Microstepping:  32,                        // DM320T supports up to 1/32
		StepMode:       ThirtySecondStep,
		CurrentLimit:   2.0,                       // DM320T handles higher current
		DecayMode:      AutoDecay,                 // DM320T has intelligent decay
		EnableLogic:    true,                      // Active high enable
		StepPulseWidth: time.Microsecond * 1,      // 1µs minimum (faster)
	}
}

// GetDM556TConfig returns optimal configuration for DM556T
func GetDM556TConfig() DriverConfig {
	return DriverConfig{
		Type:           DriverDM556T,
		Microstepping:  128,                       // DM556T supports up to 1/128 microstepping
		StepMode:       ThirtySecondStep,          // Start conservative, can go higher
		CurrentLimit:   5.6,                       // DM556T handles up to 5.6A
		DecayMode:      AutoDecay,                 // DM556T has intelligent auto-decay
		EnableLogic:    true,                      // Active high enable
		StepPulseWidth: time.Nanosecond * 500,     // 0.5µs minimum (very fast)
	}
}

// DRV8825Driver implements StepperDriver for DRV8825
type DRV8825Driver struct {
	config DriverConfig
	m0Pin  int // Mode pin 0
	m1Pin  int // Mode pin 1  
	m2Pin  int // Mode pin 2
}

// NewDRV8825Driver creates a new DRV8825 driver
func NewDRV8825Driver(m0Pin, m1Pin, m2Pin int) *DRV8825Driver {
	return &DRV8825Driver{
		config: GetDRV8825Config(),
		m0Pin:  m0Pin,
		m1Pin:  m1Pin,
		m2Pin:  m2Pin,
	}
}

// Initialize configures the DRV8825 driver
func (d *DRV8825Driver) Initialize(config DriverConfig) error {
	d.config = config
	
	// Set microstepping mode pins
	return d.SetMicrostepping(config.StepMode)
}

// SetMicrostepping configures DRV8825 microstepping via mode pins
func (d *DRV8825Driver) SetMicrostepping(mode StepMode) error {
	// DRV8825 microstepping truth table:
	// M0 M1 M2 | Step Mode
	// 0  0  0  | Full step
	// 1  0  0  | Half step  
	// 0  1  0  | 1/4 step
	// 1  1  0  | 1/8 step
	// 0  0  1  | 1/16 step
	// 1  0  1  | 1/32 step
	
	var m0, m1, m2 bool
	
	switch mode {
	case FullStep:
		m0, m1, m2 = false, false, false
	case HalfStep:
		m0, m1, m2 = true, false, false
	case QuarterStep:
		m0, m1, m2 = false, true, false
	case EighthStep:
		m0, m1, m2 = true, true, false
	case SixteenthStep:
		m0, m1, m2 = false, false, true
	case ThirtySecondStep:
		m0, m1, m2 = true, false, true
	default:
		return fmt.Errorf("unsupported step mode for DRV8825: %v", mode)
	}
	
	// In real implementation, set GPIO pins:
	// gpio.SetPin(d.m0Pin, m0)
	// gpio.SetPin(d.m1Pin, m1)
	// gpio.SetPin(d.m2Pin, m2)
	
	fmt.Printf("[DRV8825] Set microstepping mode: M0=%t, M1=%t, M2=%t (mode: %v)\n", 
		m0, m1, m2, mode)
	
	return nil
}

// SetCurrentLimit sets current limit via VREF (external resistor/potentiometer)
func (d *DRV8825Driver) SetCurrentLimit(current float32) error {
	// DRV8825 current is set via VREF voltage, not digitally controllable
	// This would require external DAC or digital potentiometer
	fmt.Printf("[DRV8825] Current limit setting requires external VREF adjustment to %.2fA\n", current)
	return nil
}

// GetDriverInfo returns DRV8825 capabilities
func (d *DRV8825Driver) GetDriverInfo() DriverInfo {
	return DriverInfo{
		Type:              DriverDRV8825,
		SupportedModes:    []StepMode{FullStep, HalfStep, QuarterStep, EighthStep, SixteenthStep, ThirtySecondStep},
		MaxCurrent:        2.2, // 2.2A max with proper cooling
		MinStepPulseWidth: time.Microsecond * 2,
		SupportsCurrentControl: false, // Requires external VREF adjustment
		SupportsDecayControl:   true,  // Fixed/Fast decay modes via DECAY pin
		RequiresModeSet:        true,  // Needs mode pins configured
	}
}

// ValidateConfig validates DRV8825 configuration
func (d *DRV8825Driver) ValidateConfig(config DriverConfig) error {
	if config.Type != DriverDRV8825 {
		return fmt.Errorf("config type mismatch: expected %s, got %s", DriverDRV8825, config.Type)
	}
	
	if config.CurrentLimit > 2.2 {
		return fmt.Errorf("current limit %.2fA exceeds DRV8825 maximum of 2.2A", config.CurrentLimit)
	}
	
	return nil
}

// DM320TDriver implements StepperDriver for DM320T
type DM320TDriver struct {
	config   DriverConfig
	sw1Pin   int // DIP switch 1 (or digital control pin)
	sw2Pin   int // DIP switch 2
	sw3Pin   int // DIP switch 3
	sw4Pin   int // DIP switch 4
}

// NewDM320TDriver creates a new DM320T driver
func NewDM320TDriver(sw1Pin, sw2Pin, sw3Pin, sw4Pin int) *DM320TDriver {
	return &DM320TDriver{
		config: GetDM320TConfig(),
		sw1Pin: sw1Pin,
		sw2Pin: sw2Pin,
		sw3Pin: sw3Pin,
		sw4Pin: sw4Pin,
	}
}

// Initialize configures the DM320T driver
func (d *DM320TDriver) Initialize(config DriverConfig) error {
	d.config = config
	return d.SetMicrostepping(config.StepMode)
}

// SetMicrostepping configures DM320T microstepping
func (d *DM320TDriver) SetMicrostepping(mode StepMode) error {
	// DM320T microstepping (if digitally controllable):
	// Different DIP switch combinations for different modes
	// This varies by specific DM320T model
	
	fmt.Printf("[DM320T] Set microstepping mode: %v\n", mode)
	fmt.Printf("[DM320T] Note: May require physical DIP switch configuration\n")
	
	return nil
}

// SetCurrentLimit sets current limit (DM320T supports digital current control)
func (d *DM320TDriver) SetCurrentLimit(current float32) error {
	// DM320T often supports digital current control
	fmt.Printf("[DM320T] Set current limit to %.2fA\n", current)
	return nil
}

// GetDriverInfo returns DM320T capabilities
func (d *DM320TDriver) GetDriverInfo() DriverInfo {
	return DriverInfo{
		Type:              DriverDM320T,
		SupportedModes:    []StepMode{FullStep, HalfStep, QuarterStep, EighthStep, SixteenthStep, ThirtySecondStep},
		MaxCurrent:        3.0, // Higher current capability
		MinStepPulseWidth: time.Microsecond * 1,
		SupportsCurrentControl: true,  // Digital current control
		SupportsDecayControl:   true,  // Intelligent decay modes
		RequiresModeSet:        false, // May have auto-configuration
	}
}

// ValidateConfig validates DM320T configuration
func (d *DM320TDriver) ValidateConfig(config DriverConfig) error {
	if config.Type != DriverDM320T {
		return fmt.Errorf("config type mismatch: expected %s, got %s", DriverDM320T, config.Type)
	}
	
	if config.CurrentLimit > 3.0 {
		return fmt.Errorf("current limit %.2fA exceeds DM320T maximum of 3.0A", config.CurrentLimit)
	}
	
	return nil
}

// DM556TDriver implements StepperDriver for DM556T
type DM556TDriver struct {
	config DriverConfig
	sw1Pin int // DIP switch 1 or digital control pin
	sw2Pin int // DIP switch 2  
	sw3Pin int // DIP switch 3
	sw4Pin int // DIP switch 4
	sw5Pin int // DIP switch 5 (for higher microstepping)
	sw6Pin int // DIP switch 6
}

// NewDM556TDriver creates a new DM556T driver
func NewDM556TDriver(sw1Pin, sw2Pin, sw3Pin, sw4Pin, sw5Pin, sw6Pin int) *DM556TDriver {
	return &DM556TDriver{
		config: GetDM556TConfig(),
		sw1Pin: sw1Pin,
		sw2Pin: sw2Pin,
		sw3Pin: sw3Pin,
		sw4Pin: sw4Pin,
		sw5Pin: sw5Pin,
		sw6Pin: sw6Pin,
	}
}

// Initialize configures the DM556T driver
func (d *DM556TDriver) Initialize(config DriverConfig) error {
	d.config = config
	return d.SetMicrostepping(config.StepMode)
}

// SetMicrostepping configures DM556T microstepping via DIP switches
func (d *DM556TDriver) SetMicrostepping(mode StepMode) error {
	// DM556T microstepping configuration via SW1-SW3:
	// SW1 SW2 SW3 | Microstepping
	// ON  ON  ON  | 1 (Full step)
	// OFF ON  ON  | 2 (Half step)
	// ON  OFF ON  | 4 (Quarter step)
	// OFF OFF ON  | 8 (Eighth step)
	// ON  ON  OFF | 16 (Sixteenth step)
	// OFF ON  OFF | 32 (Thirty-second step)
	// ON  OFF OFF | 64 (Sixty-fourth step)
	// OFF OFF OFF | 128 (One-twenty-eighth step)
	
	var sw1, sw2, sw3 bool
	
	switch mode {
	case FullStep:
		sw1, sw2, sw3 = true, true, true      // ON ON ON
	case HalfStep:
		sw1, sw2, sw3 = false, true, true     // OFF ON ON
	case QuarterStep:
		sw1, sw2, sw3 = true, false, true     // ON OFF ON
	case EighthStep:
		sw1, sw2, sw3 = false, false, true    // OFF OFF ON
	case SixteenthStep:
		sw1, sw2, sw3 = true, true, false     // ON ON OFF
	case ThirtySecondStep:
		sw1, sw2, sw3 = false, true, false    // OFF ON OFF
	default:
		// For higher microstepping modes not defined in enum
		sw1, sw2, sw3 = false, true, false    // Default to 1/32
	}
	
	// In real implementation, set GPIO pins or physical DIP switches:
	// gpio.SetPin(d.sw1Pin, sw1)
	// gpio.SetPin(d.sw2Pin, sw2)
	// gpio.SetPin(d.sw3Pin, sw3)
	
	fmt.Printf("[DM556T] Set microstepping mode: SW1=%t, SW2=%t, SW3=%t (mode: %v)\n", 
		sw1, sw2, sw3, mode)
	
	return nil
}

// SetCurrentLimit sets current limit via DIP switches SW4-SW6
func (d *DM556TDriver) SetCurrentLimit(current float32) error {
	// DM556T current settings via SW4-SW6:
	// SW4 SW5 SW6 | Peak Current (A)
	// ON  ON  ON  | 1.0
	// OFF ON  ON  | 1.5  
	// ON  OFF ON  | 2.1
	// OFF OFF ON  | 2.5
	// ON  ON  OFF | 3.0
	// OFF ON  OFF | 3.5
	// ON  OFF OFF | 4.2
	// OFF OFF OFF | 5.6
	
	var sw4, sw5, sw6 bool
	var actualCurrent float32
	
	switch {
	case current <= 1.0:
		sw4, sw5, sw6 = true, true, true    // 1.0A
		actualCurrent = 1.0
	case current <= 1.5:
		sw4, sw5, sw6 = false, true, true   // 1.5A
		actualCurrent = 1.5
	case current <= 2.1:
		sw4, sw5, sw6 = true, false, true   // 2.1A
		actualCurrent = 2.1
	case current <= 2.5:
		sw4, sw5, sw6 = false, false, true  // 2.5A
		actualCurrent = 2.5
	case current <= 3.0:
		sw4, sw5, sw6 = true, true, false   // 3.0A
		actualCurrent = 3.0
	case current <= 3.5:
		sw4, sw5, sw6 = false, true, false  // 3.5A
		actualCurrent = 3.5
	case current <= 4.2:
		sw4, sw5, sw6 = true, false, false  // 4.2A
		actualCurrent = 4.2
	default:
		sw4, sw5, sw6 = false, false, false // 5.6A (maximum)
		actualCurrent = 5.6
	}
	
	// In real implementation:
	// gpio.SetPin(d.sw4Pin, sw4)
	// gpio.SetPin(d.sw5Pin, sw5)
	// gpio.SetPin(d.sw6Pin, sw6)
	
	fmt.Printf("[DM556T] Set current limit: SW4=%t, SW5=%t, SW6=%t (%.1fA requested, %.1fA actual)\n", 
		sw4, sw5, sw6, current, actualCurrent)
	
	return nil
}

// GetDriverInfo returns DM556T capabilities
func (d *DM556TDriver) GetDriverInfo() DriverInfo {
	return DriverInfo{
		Type:              DriverDM556T,
		SupportedModes:    []StepMode{FullStep, HalfStep, QuarterStep, EighthStep, SixteenthStep, ThirtySecondStep},
		MaxCurrent:        5.6, // 5.6A maximum current
		MinStepPulseWidth: time.Nanosecond * 500, // 0.5µs minimum
		SupportsCurrentControl: true,  // Digital current control via DIP switches
		SupportsDecayControl:   true,  // Advanced auto-decay control
		RequiresModeSet:        true,  // Requires DIP switch configuration
	}
}

// ValidateConfig validates DM556T configuration
func (d *DM556TDriver) ValidateConfig(config DriverConfig) error {
	if config.Type != DriverDM556T {
		return fmt.Errorf("config type mismatch: expected %s, got %s", DriverDM556T, config.Type)
	}
	
	if config.CurrentLimit > 5.6 {
		return fmt.Errorf("current limit %.2fA exceeds DM556T maximum of 5.6A", config.CurrentLimit)
	}
	
	return nil
}
