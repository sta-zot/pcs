package domain

import "errors"

var (
	ErrVendorResolvingError  = errors.New("vendor resolving error")
	ErrFileGeneratorNotExist = errors.New("file generator not exist")
	ErrFileGenerationError   = errors.New("file generation error")
	ErrNotFound              = errors.New("not found")
	ErrInvalidSlot           = errors.New("invalid slot number")
	ErrConfigIncomplete      = errors.New("config incomplete")
	ErrInvalidMac            = errors.New("invalid mac address")
)
