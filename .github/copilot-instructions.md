<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Robot Arm Controller - Go Implementation

This is a Go project for controlling a 6-DOF robot arm using stepper motors. The project uses a clean, modular architecture with hardware abstraction.

## Code Style Guidelines

- Follow standard Go conventions and formatting (use `gofmt`)
- Use meaningful variable and function names
- Add comprehensive documentation with godoc comments
- Implement proper error handling with wrapped errors
- Use interfaces for hardware abstraction
- Include unit tests for all public functions
- Use structured logging for debugging and monitoring

## Architecture Patterns

- **Hardware Abstraction**: Use the `Board` interface for different hardware platforms
- **Interface-Based Design**: Stepper motors implement the `StepperMotor` interface
- **Concurrent Safety**: Use mutexes for shared state in the robot arm
- **Error Handling**: Return detailed errors with context
- **Configuration**: Use struct-based configuration for motors and joints

## Key Components

- `pkg/stepper`: Core stepper motor control and interfaces
- `pkg/robot`: High-level robot arm control and coordination
- `internal/hardware`: Hardware-specific implementations and board abstraction
- `cmd/arm`: Main application and demo sequences

## Testing

- Use mock implementations for hardware-independent testing
- Test both individual components and integration scenarios
- Include benchmarks for performance-critical motor control code
- Test error conditions and edge cases

## Dependencies

- Minimal external dependencies preferred
- Use standard library when possible
- For real hardware: consider `periph.io` for GPIO control
- For logging: use standard `log` package or structured logging library
