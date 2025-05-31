package config

import "errors"

// Configuration related errors
var (
	// ErrMissingODAConverterPath is returned when the ODA Converter path is not set
	ErrMissingODAConverterPath = errors.New("ODA Converter path is not set")
	// ErrODAConverterNotFound is returned when the ODA Converter executable is not found
	ErrODAConverterNotFound = errors.New("ODA Converter executable not found")
	// ErrInvalidODAConverterPath is returned when the ODA Converter path is not a valid executable
	ErrInvalidODAConverterPath = errors.New("ODA Converter path is not a valid executable")
)
