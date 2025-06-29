package stepper

import "errors"

var (
	// ErrMotorNotEnabled is returned when trying to move a disabled motor
	ErrMotorNotEnabled = errors.New("motor is not enabled")
	
	// ErrInvalidSpeed is returned when setting an invalid speed
	ErrInvalidSpeed = errors.New("invalid speed: must be positive and within motor limits")
	
	// ErrHardwareNotInitialized is returned when hardware functions are not set
	ErrHardwareNotInitialized = errors.New("hardware functions not initialized")
	
	// ErrInvalidPin is returned when a GPIO pin number is invalid
	ErrInvalidPin = errors.New("invalid GPIO pin number")
	
	// ErrMotorTimeout is returned when a motor operation times out
	ErrMotorTimeout = errors.New("motor operation timed out")
)
