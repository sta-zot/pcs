package domain

import "errors"

var (
	ErrVendorResolvingError  = errors.New("vendor resolving error")
	ErrFileGeneratorNotExist = errors.New("file generator not exist")
	ErrFileGenerationError   = errors.New("file generation error")
	ErrSettingsNotFound      = errors.New("settings not found")
	ErrInvalidSlot           = errors.New("invalid slot number")
	ErrConfigIncomplete      = errors.New("config incomplete")
	ErrInvalidMac            = errors.New("invalid mac address")
	ErrVendorNotFound        = errors.New("no manufacturer could be identified")
	ErrVendorMismatch        = errors.New("user-agent vendor does not match vendor from OUI")
	ErrPhoneNotSupported     = errors.New("vendor or model of phone not supported")
)
