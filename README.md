# 6-DOF Robot Arm Controller (Go Implementation)

A Go-based controller for operating a 6-degrees-of-freedom robot arm using stepper motors. This implementation provides a clean, modular architecture with hardware abstraction for different development boards, mixed stepper motor drivers, and integrated SSD1306 OLED display support.

## Features

- 🤖 **6-DOF Robot Arm Control**: Complete control system for 6 joints (Pan, Tilt, Joint3, Joint4, Joint5, Gripper)
- 🔧 **Mixed Driver Support**: Supports DM556T (high-performance), DM320T (intelligent), and DRV8825 (cost-effective) drivers
- 🎯 **Optimized Motor Assignment**: Large motors (Pan/Tilt) use DM556T, medium motor (Joint3) uses DM320T, smaller motors use DRV8825
- 📺 **SSD1306 OLED Display**: Real-time status display with joint positions, angles, and system status
- 🔄 **Sequence Execution**: Execute complex movement sequences with configurable delays
- 🏠 **Homing System**: Automatic homing sequence with visual progress feedback
- 🧪 **Mock Implementation**: Complete mock hardware for development and testing
- 📊 **Status Monitoring**: Real-time status reporting for all joints and motors
- 🛡️ **Safety Features**: Position limits, enable/disable controls, and error handling
- 🌐 **GPIO Extender Support**: AW9523 I2C GPIO extender for pin expansion
-  kinematics: Forward and inverse kinematics for Cartesian control (work in progress).

## SSD1306 Display Features

- **Real-time Status**: Joint positions, targets, and system status
- **Visual Progress Bars**: Joint position relative to range with target indicators
- **Homing Progress**: Visual calibration progress for each joint
- **Error Display**: Clear error messages and status indicators
- **Startup Screen**: Professional boot sequence with progress animation
- **Customizable**: Easy to modify display layouts and add custom graphics

### Display Layout

```
Robot 15:04:05
STS:HEM
Op: Homing robot
pan:100  [████████░░]
tlt:-50  [███░░░░░░░]
jt3:200  [██████████]
```

Status indicators: H=Homed, E=Enabled, M=Moving

## Supported Boards

- **Mock Board**: For development and testing (no hardware required)
- **Raspberry Pi**: Full GPIO support with optimized pin mappings
- **Google Coral Dev Board**: Hardware-specific implementation
- **Generic Linux GPIO**: For other Linux boards with GPIO support
- Extensible architecture for additional boards

## Quick Start on Linux Boards

### Automated Setup (Recommended)

For Raspberry Pi, Coral, or other Linux boards:

```bash
# Clone the repository
git clone https://github.com/gemblerz/arm.git
cd arm

# Run the automated setup script
chmod +x setup-linux.sh
sudo ./setup-linux.sh
```

The setup script will:
- Detect your board type automatically
- Install Go if not present
- Configure GPIO permissions
- Build the project
- Create startup scripts

### Manual Installation

```bash
# Ensure Go 1.21+ is installed
go version

# Install system dependencies (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install git build-essential

# For Raspberry Pi: enable GPIO access
sudo usermod -a -G gpio $USER
# (logout and login again, or reboot)

# Build the project
go mod tidy
go build -o arm-controller ./cmd/arm
```

```
pkg/
├── stepper/          # Stepper motor control package
│   ├── motor.go      # Core stepper motor implementation
│   ├── mock.go       # Mock motor for testing
│   └── errors.go     # Error definitions
├── robot/            # Robot arm control package
│   └── arm.go        # 6-DOF arm implementation
internal/
├── hardware/         # Hardware abstraction layer
│   └── board.go      # Board interfaces and implementations
cmd/
└── arm/              # Main application
    └── main.go       # Demo application
```

## Installation

### Prerequisites
- Go 1.21 or later
- Git

### Clone and Build

```bash
git clone https://github.com/gemblerz/arm.git
cd arm
go mod tidy
go build -o arm-controller ./cmd/arm
```

## Usage

### Command Line Options

```bash
# Run with different board types
./arm-controller -board=mock              # Mock hardware (default)
./arm-controller -board=raspberry-pi      # Raspberry Pi
./arm-controller -board=coral             # Google Coral Dev Board
./arm-controller -board=linux-gpio        # Generic Linux GPIO

# Configuration options
./arm-controller -board=raspberry-pi -config=raspberry-pi -verbose=true
./arm-controller -board=coral -config=coral-mixed -verbose=true
```

### Quick Run Commands

```bash
# Development with mock hardware
make run

# Run in interactive mode
make interactive
```

### Run as System Service

After running the setup script:

```bash
# Install and start as system service
./install-service.sh

# Check service status
sudo systemctl status robot-arm

# View live logs
sudo journalctl -u robot-arm -f

# Stop/start service
sudo systemctl stop robot-arm
sudo systemctl start robot-arm
```

### GPIO Pin Mappings

#### Raspberry Pi (BCM numbering)
| Joint | Step Pin | Dir Pin | Enable Pin |
|-------|----------|---------|------------|
| Base | 18 | 19 | 20 |
| Shoulder | 21 | 22 | 23 |
| Elbow | 24 | 25 | 26 |
| Wrist Pitch | 27 | 17 | 16 |
| Wrist Roll | 13 | 12 | 6 |
| Gripper | 5 | 4 | 3 |

**Note**: GPIO pins 0, 1, 14, 15 are reserved for I2C and UART.

#### Google Coral Dev Board
| Joint | Step Pin | Dir Pin | Enable Pin |
|-------|----------|---------|------------|
| Base | 142 | 143 | 144 |
| Shoulder | 145 | 146 | 147 |
| Elbow | 148 | 149 | 150 |
| Wrist Pitch | 151 | 152 | 153 |
| Wrist Roll | 154 | 155 | 156 |
| Gripper | 157 | 158 | - |

**Note**: Coral uses GPIO pins 138-158 range. Pins 138-141 are reserved for SPI.

## Google Coral Dev Board Setup

### Hardware Requirements
- Google Coral Dev Board
- 6x Stepper motors (NEMA 17 or similar)
- 6x Stepper motor drivers (A4988, DRV8825, TMC2209, etc.)
- Power supply for motors (12V/24V depending on motors)
- Jumper wires and breadboard/PCB

### Coral-Specific Features
- **Edge TPU Integration**: Ready for AI-enhanced motion planning
- **High-Performance GPIO**: Optimized for real-time motor control
- **3.3V Logic**: Compatible with most stepper motor drivers
- **40-pin Header**: Standard GPIO connector

### Quick Start for Coral

1. **Flash Mendel Linux** (if not already done):
   ```bash
   # Follow Google's official Coral setup guide
   # https://coral.ai/docs/dev-board/get-started/
   ```

2. **Run the automated setup**:
   ```bash
   git clone https://github.com/gemblerz/arm.git
   cd arm
   sudo ./setup-linux.sh
   ```

3. **Start the robot controller**:
   ```bash
   sudo ./arm-controller -board=coral
   ```

### Coral GPIO Considerations

- **GPIO Range**: Uses pins 138-158 for motor control
- **Reserved Pins**: 138-141 are used for SPI (avoid for motors)
- **Performance**: Coral's CPU can handle high-frequency stepping
- **Edge TPU**: Available for advanced motion planning algorithms

### Wiring Example for Coral

```
Coral Dev Board → Stepper Driver → Stepper Motor
GPIO 142        → STEP           → Motor 1 (Base)
GPIO 143        → DIR            → Motor 1 (Base)  
GPIO 144        → ENABLE         → Motor 1 (Base)
```

### Troubleshooting Coral

1. **Permission Issues**:
   ```bash
   sudo usermod -a -G gpio $USER
   sudo chmod 666 /dev/gpiochip*
   ```

2. **Check GPIO Status**:
   ```bash
   gpioinfo  # Show all GPIO chips and pins
   ```

3. **Monitor Resource Usage**:
   ```bash
   htop  # Check CPU usage during operation
   ```
## Development

### Adding New Hardware Boards

1. Implement the `Board` interface in `internal/hardware/board.go`
2. Add your board type to the `BoardFactory`
3. Implement hardware-specific GPIO/stepper motor controls

### Testing with Mock Hardware

The mock implementation provides full functionality without requiring physical hardware:

```go
// Create a mock board
board := hardware.NewMockBoard("Test Board", true)

// Create mock stepper motor
config := stepper.Config{
    StepPin: 2, DirPin: 3, EnablePin: 4,
    StepsPerRev: 200, MaxSpeed: 1000, DefaultSpeed: 200,
}
motor := stepper.NewMockMotor(config, true)
```

## API Examples

### Basic Motor Control

```go
// Create and configure a motor
motor := stepper.NewMockMotor(config, true)
motor.Enable()
motor.SetSpeed(100) // steps per second

// Move motor
motor.Step(200, stepper.Clockwise)
motor.MoveTo(500) // absolute position
```

### Robot Arm Control

```go
// Create robot arm with joints
arm := robot.NewArm(joints)
arm.EnableAll()
arm.Home()

// Move individual joint
arm.MoveJoint("base", 100)

// Coordinated movement
targets := map[string]int{
    "base": 200,
    "shoulder": 100,
    "elbow": -150,
}
arm.MoveJoints(targets)

### Cartesian Control (Work in Progress)
// Move the end-effector to a specific (x, y, z) coordinate.
// Note: This requires a fully implemented and calibrated kinematics solver.
arm.MoveToXYZ(150.0, 50.0, 100.0)
```

## Configuration

Motor configurations can be customized for different hardware setups:

```go
config := stepper.Config{
    StepPin:      2,           // GPIO pin for step signal
    DirPin:       3,           // GPIO pin for direction
    EnablePin:    4,           // GPIO pin for enable
    StepsPerRev:  200,         // Steps per revolution
    MaxSpeed:     1000,        // Max steps per second
    DefaultSpeed: 200,         // Default speed
    StepDelay:    time.Microsecond * 1000, // Min delay between steps
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is open source. Please check the LICENSE file for details.

## Display Integration

### SSD1306 Setup

#### Hardware Connections

For **Raspberry Pi** or **Coral Dev Board**:
```
SSD1306    Board
VCC   →    3.3V
GND   →    GND  
SDA   →    GPIO 2 (I2C1 SDA)
SCL   →    GPIO 3 (I2C1 SCL)
```

#### Configuration

```go
// Create SSD1306 display
displayConfig := display.DisplayConfig{
    Width:      128,        // 128x64 or 128x32
    Height:     64,
    I2CAddress: 0x3C,      // or 0x3D
    I2CBus:     1,         // I2C bus number
    Mock:       false,     // Set to true for testing
    Verbose:    true,
}

ssd1306, err := display.NewSSD1306(displayConfig)
if err != nil {
    log.Fatalf("Failed to create display: %v", err)
}

// Create display adapter for robot
displayAdapter := robot.NewDisplayAdapter(ssd1306, 100*time.Millisecond)

// Set display for robot arm
arm.SetDisplay(displayAdapter)
```

### Examples with Display

```bash
# Run robot with display demo
make run-display-demo

# Run 6DOF demo (without display)
make run-6dof-demo

# Run coral configuration example
make run-coral-example
```