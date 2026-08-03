package provisioning

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
	"io"
)

// Provisioning service interface for getting and creating device settings
// used by handlers
type provisioningService interface {
	GetConfig(ctx context.Context, rInfo reqInfo) ([]byte, error)
}

// Vendor resolver interface for resolving the vendor of a device based on its MAC address
type vendorIdentifier interface {
	Identify(ctx context.Context, filename, userAgent string) (*domain.DeviceInfo, error)
}

// Setting service interface for getting and creating device settings structs
type settingsProvider interface {
	// Get the settings for the device with the given MAC address
	// returns: if no settings are found nil and an error
	// error: If settings are not found, returns domain.ErrNotFound
	// error: if the settings are not complete, returns domain.ErrIncomplete
	Get(ctx context.Context, mac string) (*domain.PhoneSettings, error)
	// Create new settings for the device
	// parameters: device info - contains model, vendor, mac  and ip addresses
	Create(ctx context.Context, settings *domain.PhoneSettings) error
}

type configGenerator interface {
	Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error)
}

// Logger
type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}
