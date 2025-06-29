# Google Coral Dev Board + Mixed Stepper Drivers Setup Guide

## Your 6-DOF Robot Configuration

### Robot Joint Layout
- **Motor 1**: Pan (Base rotation) - Large motor, high torque for base movement
- **Motor 2**: Tilt (Elevation/Shoulder) - Large motor, precision vertical movement  
- **Motor 3**: 3rd Joint from base - Medium motor, arm articulation
- **Motor 4**: 4th Joint - Smaller motor, wrist/forearm movement
- **Motor 5**: 5th Joint - Smaller motor, wrist rotation
- **Motor 6**: Gripper/End Effector - Smallest motor, precision gripping

### Stepper Motor Drivers
- 2x **DM556T** drivers → Motors 1-2 (Large motors: Pan & Tilt)
  - 128 microstepping, 5.6A max, high-performance for heavy loads
- 1x **DM320T** driver → Motor 3 (Medium motor: 3rd joint)
  - 32 microstepping, 3A max, intelligent features
- 3x **DRV8825** drivers → Motors 4-6 (Smaller motors: joints 4-5, gripper)
  - 16 microstepping, 2.2A max, cost-effective for lighter loads

### GPIO Expansion
- **AW9523 GPIO Extender** (16 additional GPIO pins via I2C)
- **Google Coral Dev Board** (native GPIO pins 138-158)

## Optimal Pin Assignment Strategy

### 🎯 High-Speed Signals → Native Coral GPIO
Critical timing signals (STEP/DIR) use Coral's native GPIO for maximum performance:

| Joint | Motor | Driver | Step Pin | Dir Pin | Speed Priority |
|-------|-------|--------|----------|---------|----------------|
| Pan | Motor 1 | DM556T #1 | **142** (Coral) | **143** (Coral) | **CRITICAL** |
| Tilt | Motor 2 | DM556T #2 | **144** (Coral) | **145** (Coral) | **CRITICAL** |
| Joint3 | Motor 3 | DM320T #1 | **146** (Coral) | **147** (Coral) | **HIGH** |
| Joint5 | Motor 5 | DRV8825 #2 | **148** (Coral) | 1006 (Extender) | **MEDIUM** |

### 🔧 Control Signals → GPIO Extender  
Enable and mode pins can use extender (I2C latency acceptable):

| Joint | Motor | Enable Pin | Mode Pins | Notes |
|-------|-------|------------|-----------|-------|
| Pan | Motor 1 | **1000** (Ext-0) | 149,150,151 (Coral) | DM556T config |
| Tilt | Motor 2 | **1001** (Ext-1) | 152,153,154 (Coral) | DM556T config |
| Joint3 | Motor 3 | **1002** (Ext-2) | 155,156,157 (Coral) | DM320T config |
| Joint4 | Motor 4 | **1005** (Ext-5) | 1011,1012,1013 (Ext) | DRV8825 mode |
| Joint5 | Motor 5 | **1007** (Ext-7) | 1014,1015,158 | Mixed mode pins |
| Gripper | Motor 6 | **1010** (Ext-10) | 1016,1011,1012 (Ext) | DRV8825 mode |

## 🔌 Wiring Recommendations

### Power Distribution
```
Main Power Supply (24V, 20A recommended for large motors)
├── Motors 1-3 Power (24V) → Large/Medium motors (DM556T, DM320T)
├── Motors 4-6 Power (12V) → Smaller motors (DRV8825)
├── Driver Logic Power (5V, 2A) → Level shifters
└── Coral Dev Board (5V, 3A) → Separate official adapter

Coral 3.3V Rail
└── AW9523 GPIO Extender (3.3V, 50mA)
```

### I2C Connection (AW9523)
```
Coral Dev Board    AW9523 GPIO Extender
Pin 3 (SDA)    →   SDA
Pin 5 (SCL)    →   SCL  
3.3V           →   VCC
GND            →   GND
```

### Signal Routing Priority
1. **Critical timing signals** → Direct Coral GPIO (142-148)
2. **Driver enable signals** → AW9523 extender (1000-1010)
3. **Mode configuration** → Mixed (Coral for DM320T, Extender for DRV8825)

## 🚀 Performance Optimization

### Driver-Specific Settings

#### DM556T Drivers (Motors 1-2: Pan & Tilt)
```go
// Large motors with high torque requirements
// Microstepping: 1/32 (balance of precision and speed)
// Current: 3.0-5.6A (high torque for heavy loads)
// Step pulse: 0.5µs minimum (very fast)
// Enable: Active HIGH

DIP Switch Configuration:
SW1 SW2 SW3 | Microstepping
ON  ON  ON  | 1 (Full step)
OFF ON  ON  | 2 (Half step)  
ON  OFF ON  | 4 (Quarter step)
OFF OFF ON  | 8 (Eighth step)
ON  ON  OFF | 16 (Sixteenth step)
OFF ON  OFF | 32 (Thirty-second step) ← recommended for Motors 1-2
ON  OFF OFF | 64 (Sixty-fourth step)
OFF OFF OFF | 128 (One-twenty-eighth step) ← max precision

SW4 SW5 SW6 | Current Limit  
ON  ON  ON  | 1.0A (low power)
OFF ON  ON  | 1.5A (typical NEMA 17)
ON  OFF ON  | 2.1A (medium torque) 
OFF OFF ON  | 3.0A (recommended for Motors 1-2) ← good balance
ON  ON  OFF | 3.5A (higher torque)
OFF ON  OFF | 4.2A (high torque)
ON  OFF OFF | 4.8A (very high)
OFF OFF OFF | 5.6A (maximum for heavy loads)
```

#### DM320T Driver (Motor 3: 3rd Joint) 
```go
// Medium motor with precision requirements
// Microstepping: 1/32 (maximum precision)
// Current: 2.0A (good torque for medium motor)
// Step pulse: 1µs minimum  
// Enable: Active HIGH
// Intelligent current decay and auto-configuration
```

#### DRV8825 Drivers (Motors 4-6: Smaller joints & gripper)
```go
// Smaller motors, cost-effective control
// Microstepping: 1/16 (smooth movement, adequate precision)
// Current: 1.5A (typical NEMA 17)
// Step pulse: 2µs minimum
// Enable: Active LOW

Mode Pins (M0,M1,M2):
- 1/16 step: LOW, LOW, HIGH ← recommended
- Configure via AW9523 extender pins
```

SW4 SW5 SW6 | Current Limit  
ON  ON  ON  | 1.0A (low power)
OFF ON  ON  | 1.5A (typical NEMA 17)
ON  OFF ON  | 2.1A (medium torque) 
OFF OFF ON  | 3.0A (good for most NEMA 17) ← recommended
ON  ON  OFF | 3.5A (higher torque)
OFF ON  OFF | 4.2A (high torque)
ON  OFF OFF | 4.8A (very high)
OFF OFF OFF | 5.6A (maximum)
```

## 📋 Setup Commands

### Run with Mixed Driver Configuration
```bash
# Use the coral-extended board type with mixed driver config
./arm-controller -board=coral-extended -config=coral-mixed -verbose=true

# Or via Makefile
make run-coral-mixed  # (add this target)
```

### Test Individual Drivers
```bash
# Test each joint separately to verify driver compatibility
./arm-controller -board=coral-extended -test-joint=base
./arm-controller -board=coral-extended -test-joint=shoulder
```

## 🔧 Driver Identification for Unknown Drivers

### Visual Inspection
Look for these markings on your unknown drivers:
- **DM320T**: Often has "DM320T" or "Leadshine" marking
- **TB6600**: Common alternative, blue/green PCB
- **A4988**: Smaller, black IC, different pin layout
- **TMC2209**: Advanced, UART communication

### Test Procedure
```bash
# Start with conservative settings for unknown drivers
# Gradually increase microstepping and current
# Monitor for:
# - Smooth operation
# - Minimal noise/vibration  
# - Proper torque delivery
```

## ⚡ Performance Expectations

### Speed Capabilities
| Joint | Motor | Max Speed | Microstepping | Driver | Notes |
|-------|-------|-----------|---------------|--------|-------|
| Pan | Motor 1 | 2000 steps/s | 1/32 | DM556T | High-speed base rotation |
| Tilt | Motor 2 | 1600 steps/s | 1/32 | DM556T | Precision elevation |
| Joint3 | Motor 3 | 1200 steps/s | 1/32 | DM320T | Medium motor precision |
| Joint4 | Motor 4 | 1000 steps/s | 1/16 | DRV8825 | Via extender |
| Joint5 | Motor 5 | 800 steps/s | 1/16 | DRV8825 | Mixed GPIO |
| Gripper | Motor 6 | 400 steps/s | 1/16 | DRV8825 | Slow precision |

### I2C Extender Limitations
- **Latency**: ~1ms additional delay for extender pins
- **Speed Impact**: 10-20% reduction for extender-based stepping
- **Mitigation**: Use native Coral GPIO for high-speed axes

## 🛠️ Troubleshooting

### Common Issues
1. **Driver not responding**: Check enable pin polarity (HIGH vs LOW)
2. **Jerky movement**: Verify microstepping configuration
3. **Insufficient torque**: Increase current limit (check driver specs)
4. **Overheating**: Add cooling, reduce current, check wiring

### Debugging Commands
```bash
# Check GPIO extender communication
i2cdetect -y 1  # Should show device at 0x58

# Monitor pin states
gpioinfo | grep "142\|143\|144"  # Check Coral pins

# Test extender pins individually
# (Custom test commands in debug mode)
```

This configuration gives you the optimal balance of performance and pin utilization for your mixed driver setup! 🎯
