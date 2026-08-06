package provisioning

import (
	"context"
	"errors"
	"github/sta-zot/pcs/internal/domain"
	"io"
)

// UseCase is the use case for provisioning
// It contains the business logic for provisioning a device
// 1. Use identifier service for getting device vendor and model
// 2. Use Settings service for getting device settings
// 3. Use the generator factory to get a generator depending on the device.
// 4. Use the generator to generate the device configuration.
type UseCase struct {
	sProv            settingsProvider
	logger           logger
	generator        configGenerator
	vendorIdentifier vendorIdentifier
}

// Provision provisions a device
// It takes a context and request info as parameters
// Returns the device configuration as a byte slice and an error
func (uc *UseCase) Provision(ctx context.Context, reqInfo reqInfo) (io.ByteReader, error) {
	deviceInfo, err := uc.vendorIdentifier.Identify(ctx, reqInfo.Filename, reqInfo.UserAgent)
	if err != nil {
		return nil, err
	}
	settings, err := uc.sProv.Get(ctx, deviceInfo.MAC())
	if err != nil {
		if errors.Is(err, domain.ErrSettingsNotFound) {
			deviceInfo.SetIP(reqInfo.IP)
			settings := domain.NewPhoneSettings(*deviceInfo)
			go uc.sProv.Create(ctx, settings)
			return nil, err
		}
		return nil, err
	}
	if settings != nil {
		return uc.generator.Generate(ctx, settings)
	}
	return nil, errors.New("setting is nil")
}

// NewUseCase creates a new UseCase
// It takes a setting service, logger, generator factory, and vendor resolver as parameters
// Returns a pointer to the new UseCase
// parameters:
// - settingService: interface for the setting service to use for getting device settings
// - logger: interface for the logger to use for logging
// - generatorFactory: interface for the generator factory to use for getting a generator
// - vendorIdentifier: interface for the vendor resolver to use for resolving device vendor and model
func NewUseCase(
	settingsProvider settingsProvider,
	vendorIdentifier vendorIdentifier,
	generator configGenerator,
	logger logger,
) *UseCase {
	return &UseCase{
		sProv:            settingsProvider,
		logger:           logger,
		generator:        generator,
		vendorIdentifier: vendorIdentifier,
	}
}
