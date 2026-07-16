package provisioning

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
)

// Provisioning service interface for getting and creating device settings
// used by handlers
type provisioningService interface {
	GetConfig(ctx context.Context, rInfo reqInfo) ([]byte, error)
}

// Vendor resolver interface for resolving the vendor of a device based on its MAC address
type vendorResolver interface {
	Get(ctx context.Context, filename, userAgetn string) (*domain.DeviceInfo, error)
}

// Setting service interface for getting and creating device settings structs
type settingService interface {
	// Get the settings for the device with the given MAC address
	// returns: if no settings are found nil and an error
	// error: If settings are not found, returns domain.ErrNotFound
	// error: if the settings are not complete, returns domain.ErrIncomplete
	Get(ctx context.Context, mac string) (*domain.PfoneSettings, error)
	// Create new settings for the device
	// parameters: device info - contains model, vendor, mac  and ip addresses
	Create(*domain.DeviceInfo) error
}

// Generator factory interface for creating Generators based on model and vendor
type GeneratorFactory interface {
	Get(model, vendor string) (Generator, error)
}

// Generator interface for generating device configuration files
type Generator interface {
	Generate(ctx context.Context, settings *domain.PfoneSettings) ([]byte, error)
}

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}
