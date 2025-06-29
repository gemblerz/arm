package hardware

import (
	"time"
	
	"github.com/gemblerz/arm/pkg/stepper"
)

// CoralMixedDriverConfig creates optimal configuration for Coral + GPIO extender + mixed drivers
func GetCoralMixedDriverConfig() BoardConfig {
	return BoardConfig{
		Name: "Coral Dev Board + AW9523 GPIO Extender + 6-DOF Robot (Pan-Tilt-Joint3-Joint4-Joint5-Gripper)",
		Type: "coral-extended-6dof",
		GPIO: GPIOConfig{
			Numbering: "Coral + AW9523",
			Voltage:   "3.3V",
		},
		Motors: []MotorConfig{
			// Motor 1: Pan (Base rotation) - DM556T #1 (Large motor, high torque)
			{
				JointID: "pan",
				Config: stepper.Config{
					StepPin:      142, // Coral native GPIO for high-speed stepping
					DirPin:       143, // Coral native GPIO  
					EnablePin:    1000, // AW9523 extender pin 0 (1000 + 0)
					StepsPerRev:  6400, // 200 * 32 microstepping (DM556T precision)
					MaxSpeed:     2000, // High speed for pan movement
					DefaultSpeed: 500,
					StepDelay:    time.Nanosecond * 500, // DM556T fast pulse
				},
				MinPos: -12800, // ±2 full rotations (360° each direction)
				MaxPos: 12800,
			},
			
			// Motor 2: Tilt (Shoulder/Elevation) - DM556T #2 (Large motor, high precision)
			{
				JointID: "tilt", 
				Config: stepper.Config{
					StepPin:      144, // Coral native GPIO
					DirPin:       145, // Coral native GPIO
					EnablePin:    1001, // AW9523 extender pin 1
					StepsPerRev:  6400, // 200 * 32 microstepping (DM556T precision)
					MaxSpeed:     1600, // High precision tilt movement
					DefaultSpeed: 400,
					StepDelay:    time.Nanosecond * 500, // DM556T fast pulse
				},
				MinPos: -3200, // ±90° tilt range (conservative for safety)
				MaxPos: 3200,
			},
			
			// Motor 3: 3rd Joint from base - DM320T (Medium motor)
			{
				JointID: "joint3",
				Config: stepper.Config{
					StepPin:      146, // Coral native GPIO
					DirPin:       147, // Coral native GPIO
					EnablePin:    1002, // AW9523 extender pin 2
					StepsPerRev:  6400, // 200 * 32 microstepping (DM320T)
					MaxSpeed:     1200,
					DefaultSpeed: 300,
					StepDelay:    time.Microsecond * 1, // DM320T timing
				},
				MinPos: -3200, // ±180° range
				MaxPos: 3200,
			},
			
			// Motor 4: 4th Joint - DRV8825 #1 (Smaller motor)
			{
				JointID: "joint4",
				Config: stepper.Config{
					StepPin:      1003, // AW9523 extender pin 3
					DirPin:       1004, // AW9523 extender pin 4
					EnablePin:    1005, // AW9523 extender pin 5
					StepsPerRev:  3200, // 200 * 16 microstepping (DRV8825)
					MaxSpeed:     1000,
					DefaultSpeed: 250,
					StepDelay:    time.Microsecond * 2, // Slower due to I2C overhead
				},
				MinPos: -3200, // ±180° range
				MaxPos: 3200,
			},
			
			// Motor 5: 5th Joint - DRV8825 #2 (Smaller motor)
			{
				JointID: "joint5",
				Config: stepper.Config{
					StepPin:      148, // Coral native GPIO for better performance
					DirPin:       1006, // AW9523 extender pin 6
					EnablePin:    1007, // AW9523 extender pin 7
					StepsPerRev:  3200, // 200 * 16 microstepping
					MaxSpeed:     800,
					DefaultSpeed: 200,
					StepDelay:    time.Microsecond * 2,
				},
				MinPos: -6400, // ±360° (can rotate freely)
				MaxPos: 6400,
			},
			
			// Motor 6: End Effector/Gripper - DRV8825 #3 (Smallest motor)
			// Motor 6: End Effector/Gripper - DRV8825 #3 (Smallest motor)
			{
				JointID: "gripper",
				Config: stepper.Config{
					StepPin:      1008, // AW9523 extender pin 8
					DirPin:       1009, // AW9523 extender pin 9
					EnablePin:    1010, // AW9523 extender pin 10
					StepsPerRev:  3200, // 200 * 16 microstepping (adequate for gripper)
					MaxSpeed:     400,  // Slow speed for precision control
					DefaultSpeed: 100,
					StepDelay:    time.Microsecond * 2, // Via extender, conservative timing
				},
				MinPos: 0,    // Fully open
				MaxPos: 1600, // Fully closed (0.25 rotation, 90°)
			},
		},
	}
}

// DriverAssignment represents which driver controls which joint
type DriverAssignment struct {
	JointID    string
	DriverType stepper.DriverType
	DriverID   int // For multiple drivers of same type
	ModePin1   int // Mode/config pin 1
	ModePin2   int // Mode/config pin 2  
	ModePin3   int // Mode/config pin 3
	Notes      string
}

// GetCoralDriverAssignments returns the recommended driver assignments
func GetCoralDriverAssignments() []DriverAssignment {
	return []DriverAssignment{
		{
			JointID:    "pan",
			DriverType: stepper.DriverDM556T,
			DriverID:   1,
			ModePin1:   149,  // Coral native GPIO (SW1)
			ModePin2:   150,  // Coral native GPIO (SW2)
			ModePin3:   151,  // Coral native GPIO (SW3)
			Notes:      "DM556T for large pan motor, 1/32 microstepping, up to 5.6A current",
		},
		{
			JointID:    "tilt", 
			DriverType: stepper.DriverDM556T,
			DriverID:   2,
			ModePin1:   152,  // Coral native GPIO (SW1)
			ModePin2:   153,  // Coral native GPIO (SW2) 
			ModePin3:   154,  // Coral native GPIO (SW3)
			Notes:      "DM556T for large tilt motor, 1/32 microstepping, high precision",
		},
		{
			JointID:    "joint3",
			DriverType: stepper.DriverDM320T,
			DriverID:   1,
			ModePin1:   155,  // Coral native GPIO (SW1)
			ModePin2:   156,  // Coral native GPIO (SW2)
			ModePin3:   157,  // Coral native GPIO (SW3)
			Notes:      "DM320T for medium motor, 1/32 microstepping, intelligent features",
		},
		{
			JointID:    "joint4",
			DriverType: stepper.DriverDRV8825,
			DriverID:   1,
			ModePin1:   1011, // AW9523 pin 11 (M0)
			ModePin2:   1012, // AW9523 pin 12 (M1)
			ModePin3:   1013, // AW9523 pin 13 (M2)
			Notes:      "DRV8825 for smaller motor, 1/16 microstepping",
		},
		{
			JointID:    "joint5",
			DriverType: stepper.DriverDRV8825,
			DriverID:   2,
			ModePin1:   1014, // AW9523 pin 14 (M0)
			ModePin2:   1015, // AW9523 pin 15 (M1)
			ModePin3:   158,  // Coral native GPIO (M2)
			Notes:      "DRV8825 for smaller motor, 1/16 microstepping",
		},
		{
			JointID:    "gripper",
			DriverType: stepper.DriverDRV8825,
			DriverID:   3,
			ModePin1:   1016, // AW9523 extender pin 16 (if available, or reuse)
			ModePin2:   1011, // Reuse AW9523 pin 11 (M1)
			ModePin3:   1012, // Reuse AW9523 pin 12 (M2)
			Notes:      "DRV8825 for gripper motor, 1/16 microstepping, shared mode pins",
		},
	}
}

// WiringRecommendations provides optimal wiring setup
type WiringRecommendation struct {
	Component    string
	Connection   string
	CoralPin     string
	ExtenderPin  string
	Notes        string
}

// GetCoralWiringRecommendations returns wiring recommendations
func GetCoralWiringRecommendations() []WiringRecommendation {
	return []WiringRecommendation{
		{
			Component:   "AW9523 GPIO Extender",
			Connection:  "I2C",
			CoralPin:    "SDA/SCL (pins 3/5)",
			ExtenderPin: "SDA/SCL",
			Notes:       "I2C address 0x58 or 0x59, 3.3V power from Coral",
		},
		{
			Component:   "Base Motor (DRV8825 #1)",
			Connection:  "Step/Dir",
			CoralPin:    "GPIO 142/143",
			ExtenderPin: "Pin 0 (Enable)",
			Notes:       "High-speed signals on native Coral GPIO",
		},
		{
			Component:   "Shoulder Motor (DM320T #1)",
			Connection:  "Step/Dir",
			CoralPin:    "GPIO 144/145",
			ExtenderPin: "Pin 1 (Enable)",
			Notes:       "Precision control, native GPIO for timing",
		},
		{
			Component:   "Elbow Motor (DRV8825 #2)",
			Connection:  "Step/Dir/Enable",
			CoralPin:    "-",
			ExtenderPin: "Pins 2/3/4",
			Notes:       "All signals via extender, acceptable for medium speed",
		},
		{
			Component:   "Power Distribution",
			Connection:  "Motor Power",
			CoralPin:    "5V (logic) from Coral",
			ExtenderPin: "3.3V logic via extender",
			Notes:       "Separate 12V/24V supply for motor power",
		},
		{
			Component:   "Mode Pins",
			Connection:  "Microstepping",
			CoralPin:    "GPIO 149-158 (DM320T)",
			ExtenderPin: "Pins 11-15 (DRV8825)",
			Notes:       "Configure microstepping at startup",
		},
	}
}

// PowerRequirements calculates power needs
type PowerRequirements struct {
	Component      string
	Voltage        string
	Current        string
	Notes          string
}

// GetCoralPowerRequirements returns power requirements
func GetCoralPowerRequirements() []PowerRequirements {
	return []PowerRequirements{
		{
			Component: "Google Coral Dev Board",
			Voltage:   "5V",
			Current:   "3A",
			Notes:     "Official Coral power adapter required",
		},
		{
			Component: "AW9523 GPIO Extender", 
			Voltage:   "3.3V",
			Current:   "50mA",
			Notes:     "Powered from Coral 3.3V rail",
		},
		{
			Component: "6x Stepper Motors",
			Voltage:   "12V or 24V",
			Current:   "1.5-2A each",
			Notes:     "Separate power supply, 10-15A capacity recommended",
		},
		{
			Component: "Driver Logic (6x)",
			Voltage:   "3.3V/5V",
			Current:   "10mA each",
			Notes:     "Logic power from Coral or separate 5V supply",
		},
	}
}
