# DM556T Driver Update Summary

## Overview
Updated the robot arm controller to properly support the newly identified DM556T stepper motor drivers. The configuration has been optimized for a 6-DOF robot with specific motor sizing and driver assignments.

## New Robot Configuration (6-DOF Pan-Tilt System)

### Motor & Driver Assignment
| Motor | Joint | Size | Driver | Current | Microstepping | Notes |
|-------|-------|------|--------|---------|---------------|-------|
| 1 | Pan (Base) | Large | DM556T #1 | 3.0-5.6A | 1/32 | High torque base rotation |
| 2 | Tilt (Elevation) | Large | DM556T #2 | 3.0-5.6A | 1/32 | Precision vertical movement |
| 3 | 3rd Joint | Medium | DM320T #1 | 2.0A | 1/32 | Intelligent mid-arm joint |
| 4 | 4th Joint | Small | DRV8825 #1 | 1.5A | 1/16 | Cost-effective control |
| 5 | 5th Joint | Small | DRV8825 #2 | 1.5A | 1/16 | Wrist articulation |
| 6 | Gripper | Small | DRV8825 #3 | 1.5A | 1/16 | Precision gripping |

### Pin Assignments (Google Coral + AW9523)

#### Critical Timing Signals (Native Coral GPIO)
- **Pan Motor (DM556T #1)**: Step=142, Dir=143, Enable=1000 (Ext)
- **Tilt Motor (DM556T #2)**: Step=144, Dir=145, Enable=1001 (Ext)  
- **Joint3 (DM320T)**: Step=146, Dir=147, Enable=1002 (Ext)
- **Joint5 (DRV8825 #2)**: Step=148, Dir=1006 (Ext), Enable=1007 (Ext)

#### Mode Control Pins
- **DM556T #1 (Pan)**: SW1=149, SW2=150, SW3=151 (Coral GPIO)
- **DM556T #2 (Tilt)**: SW1=152, SW2=153, SW3=154 (Coral GPIO)
- **DM320T (Joint3)**: SW1=155, SW2=156, SW3=157 (Coral GPIO)
- **DRV8825 drivers**: Mode pins via AW9523 extender

## Updated Files

### 1. Configuration Files
- **`internal/hardware/coral-mixed-config.go`**
  - Updated joint names: pan, tilt, joint3, joint4, joint5, gripper
  - Optimized pin assignments for DM556T performance
  - Updated driver assignments and microstepping settings
  - Added proper current limit configurations

### 2. Documentation
- **`docs/CORAL_MIXED_DRIVERS.md`**
  - Complete DM556T DIP switch configuration table (SW1-SW6)
  - Updated robot joint descriptions and motor sizing
  - Optimized power distribution recommendations
  - Updated performance expectations table

- **`README.md`**
  - Updated feature list to highlight mixed driver support
  - New joint naming convention
  - Added GPIO extender support information

### 3. Test Files  
- **`pkg/stepper/motor_test.go`**
  - Added comprehensive DM556T driver tests
  - Updated joint names in configuration tests
  - Added driver capability validation tests

### 4. Examples
- **`examples/6dof-pan-tilt-config.go`**
  - Complete working example of the new configuration
  - Demonstrates pan-tilt movements with large motors
  - Shows full arm coordination sequence
  - Motor characteristic comparisons

## DM556T Driver Features

### Microstepping Configuration (SW1-SW3)
```
SW1 SW2 SW3 | Microstepping
ON  ON  ON  | 1 (Full step)
OFF ON  ON  | 2 (Half step)
ON  OFF ON  | 4 (Quarter step)  
OFF OFF ON  | 8 (Eighth step)
ON  ON  OFF | 16 (Sixteenth step)
OFF ON  OFF | 32 (Thirty-second step) ← RECOMMENDED
ON  OFF OFF | 64 (Sixty-fourth step)
OFF OFF OFF | 128 (One-twenty-eighth step)
```

### Current Limit Configuration (SW4-SW6)
```
SW4 SW5 SW6 | Current Limit
ON  ON  ON  | 1.0A
OFF ON  ON  | 1.5A
ON  OFF ON  | 2.1A
OFF OFF ON  | 3.0A ← RECOMMENDED for Motors 1-2
ON  ON  OFF | 3.5A
OFF ON  OFF | 4.2A  
ON  OFF OFF | 4.8A
OFF OFF OFF | 5.6A (Maximum)
```

### Performance Characteristics
- **Step Pulse Width**: 0.5µs minimum (very fast)
- **Maximum Current**: 5.6A (excellent for large motors)
- **Enable Logic**: Active HIGH
- **Decay Mode**: Intelligent auto-decay
- **Microstepping**: Up to 1/128 (exceptional precision)

## Key Optimizations

1. **Large Motors (Pan/Tilt)**: DM556T drivers provide high current (5.6A) and precision (1/32) for heavy loads
2. **GPIO Assignment**: Critical step/dir signals use native Coral GPIO for maximum speed
3. **Current Settings**: 3.0A recommended for Motors 1-2 (good balance of torque and heat)
4. **Microstepping**: 1/32 for large motors (smooth movement), 1/16 for smaller motors (adequate precision)
5. **Power Distribution**: 24V for large motors, 12V for smaller motors

## Testing Status

The updated configuration includes:
- ✅ Complete DM556T driver implementation
- ✅ Updated hardware abstraction layer  
- ✅ Comprehensive test coverage
- ✅ Working example code
- ✅ Complete documentation
- ✅ Optimal pin assignments for Coral + AW9523

The system is ready for deployment on Google Coral Dev Board with mixed stepper motor drivers for your 6-DOF pan-tilt robot system.
