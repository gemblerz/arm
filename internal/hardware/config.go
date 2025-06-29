package hardware

import (
	"time"
	"github.com/gemblerz/arm/pkg/stepper"
)

// BoardConfig holds configuration for different board types
type BoardConfig struct {
	Name     string
	Type     string
	GPIO     GPIOConfig
	Motors   []MotorConfig
}

// GPIOConfig contains GPIO-specific configuration
type GPIOConfig struct {
	Numbering string // "BCM" or "BOARD" for Raspberry Pi
	Voltage   string // "3.3V" or "5V"
}

// MotorConfig contains pre-configured motor settings for specific joints
type MotorConfig struct {
	JointID string
	stepper.Config
	MinPos int
	MaxPos int
}

// GetRaspberryPiConfig returns a default configuration for Raspberry Pi
func GetRaspberryPiConfig() BoardConfig {
	return BoardConfig{
		Name: "Raspberry Pi 4B",
		Type: "raspberry-pi",
		GPIO: GPIOConfig{
			Numbering: "BCM",
			Voltage:   "3.3V",
		},
		Motors: []MotorConfig{
			{
				JointID: "base",
				Config: stepper.Config{
					StepPin:      18, // GPIO 18 (PWM0)
					DirPin:       19, // GPIO 19
					EnablePin:    20, // GPIO 20
					StepsPerRev:  200,
					MaxSpeed:     1000,
					DefaultSpeed: 200,
				},
				MinPos: -1800, // -1800 steps = -9 full rotations
				MaxPos: 1800,  // +1800 steps = +9 full rotations
			},
			{
				JointID: "shoulder",
				Config: stepper.Config{
					StepPin:      21, // GPIO 21
					DirPin:       22, // GPIO 22
					EnablePin:    23, // GPIO 23
					StepsPerRev:  200,
					MaxSpeed:     800,
					DefaultSpeed: 150,
				},
				MinPos: -800, // Limited range for shoulder
				MaxPos: 800,
			},
			{
				JointID: "elbow",
				Config: stepper.Config{
					StepPin:      24, // GPIO 24
					DirPin:       25, // GPIO 25
					EnablePin:    26, // GPIO 26
					StepsPerRev:  200,
					MaxSpeed:     600,
					DefaultSpeed: 100,
				},
				MinPos: -600,
				MaxPos: 600,
			},
			{
				JointID: "wrist_pitch",
				Config: stepper.Config{
					StepPin:      27, // GPIO 27
					DirPin:       17, // GPIO 17
					EnablePin:    16, // GPIO 16
					StepsPerRev:  200,
					MaxSpeed:     400,
					DefaultSpeed: 80,
				},
				MinPos: -400,
				MaxPos: 400,
			},
			{
				JointID: "wrist_roll",
				Config: stepper.Config{
					StepPin:      13, // GPIO 13
					DirPin:       12, // GPIO 12
					EnablePin:    6,  // GPIO 6
					StepsPerRev:  200,
					MaxSpeed:     400,
					DefaultSpeed: 80,
				},
				MinPos: -400,
				MaxPos: 400,
			},
			{
				JointID: "gripper",
				Config: stepper.Config{
					StepPin:      5,  // GPIO 5
					DirPin:       4,  // GPIO 4
					EnablePin:    3,  // GPIO 3
					StepsPerRev:  200,
					MaxSpeed:     200,
					DefaultSpeed: 50,
				},
				MinPos: 0,   // Fully open
				MaxPos: 100, // Fully closed
			},
		},
	}
}

// GetCoralDevBoardConfig returns a default configuration for Google Coral Dev Board
func GetCoralDevBoardConfig() BoardConfig {
	return BoardConfig{
		Name: "Google Coral Dev Board",
		Type: "coral",
		GPIO: GPIOConfig{
			Numbering: "Coral",
			Voltage:   "3.3V",
		},
		Motors: []MotorConfig{
			{
				JointID: "base",
				Config: stepper.Config{
					StepPin:      142, // Coral GPIO pin 142
					DirPin:       143, // Coral GPIO pin 143
					EnablePin:    144, // Coral GPIO pin 144
					StepsPerRev:  200,
					MaxSpeed:     1000,
					DefaultSpeed: 200,
					StepDelay:    time.Microsecond * 1000,
				},
				MinPos: -1800, // -1800 steps = -9 full rotations
				MaxPos: 1800,  // +1800 steps = +9 full rotations
			},
			{
				JointID: "shoulder",
				Config: stepper.Config{
					StepPin:      145, // Coral GPIO pin 145
					DirPin:       146, // Coral GPIO pin 146
					EnablePin:    147, // Coral GPIO pin 147
					StepsPerRev:  200,
					MaxSpeed:     800,
					DefaultSpeed: 150,
					StepDelay:    time.Microsecond * 1200,
				},
				MinPos: -800, // Limited range for shoulder
				MaxPos: 800,
			},
			{
				JointID: "elbow",
				Config: stepper.Config{
					StepPin:      148, // Coral GPIO pin 148
					DirPin:       149, // Coral GPIO pin 149
					EnablePin:    150, // Coral GPIO pin 150
					StepsPerRev:  200,
					MaxSpeed:     600,
					DefaultSpeed: 100,
					StepDelay:    time.Microsecond * 1500,
				},
				MinPos: -600,
				MaxPos: 600,
			},
			{
				JointID: "wrist_pitch",
				Config: stepper.Config{
					StepPin:      151, // Coral GPIO pin 151
					DirPin:       152, // Coral GPIO pin 152
					EnablePin:    153, // Coral GPIO pin 153
					StepsPerRev:  200,
					MaxSpeed:     400,
					DefaultSpeed: 80,
					StepDelay:    time.Microsecond * 2000,
				},
				MinPos: -400,
				MaxPos: 400,
			},
			{
				JointID: "wrist_roll",
				Config: stepper.Config{
					StepPin:      154, // Coral GPIO pin 154
					DirPin:       155, // Coral GPIO pin 155
					EnablePin:    156, // Coral GPIO pin 156
					StepsPerRev:  200,
					MaxSpeed:     400,
					DefaultSpeed: 80,
					StepDelay:    time.Microsecond * 2000,
				},
				MinPos: -400,
				MaxPos: 400,
			},
			{
				JointID: "gripper",
				Config: stepper.Config{
					StepPin:      157, // Coral GPIO pin 157
					DirPin:       158, // Coral GPIO pin 158
					EnablePin:    -1,  // No enable pin for gripper (always enabled)
					StepsPerRev:  200,
					MaxSpeed:     200,
					DefaultSpeed: 50,
					StepDelay:    time.Microsecond * 3000,
				},
				MinPos: 0,   // Fully open
				MaxPos: 100, // Fully closed
			},
		},
	}
}

// GetGenericLinuxConfig returns a generic configuration
func GetGenericLinuxConfig() BoardConfig {
	return BoardConfig{
		Name: "Generic Linux Board",
		Type: "linux-gpio",
		GPIO: GPIOConfig{
			Numbering: "Generic",
			Voltage:   "3.3V",
		},
		Motors: []MotorConfig{
			// Use safe GPIO pin numbers that are common across platforms
			{
				JointID: "base",
				Config: stepper.Config{
					StepPin:      2,
					DirPin:       3,
					EnablePin:    4,
					StepsPerRev:  200,
					MaxSpeed:     1000,
					DefaultSpeed: 200,
				},
				MinPos: -1800,
				MaxPos: 1800,
			},
			// Add more generic configurations...
		},
	}
}
