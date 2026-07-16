package provisioning

import (
	"context"
	"errors"
	"github/sta-zot/pcs/internal/domain"
)

// UseCase is the use case for provisioning
// It contains the business logic for provisioning a device
// 1. Use identifier service for getting device vendor and model
// 2. Use Settings service for getting device settings
// 3. Use the generator factory to get a generator depending on the device.
// 4. Use the generator to generate the device configuration.
type UseCase struct {
	settingSvc     settingService
	logger         logger
	genFactory     GeneratorFactory
	vendorResolver vendorResolver
}

// Provision provisions a device
// It takes a context and request info as parameters
// Returns the device configuration as a byte slice and an error
func (uc *UseCase) Provision(ctx context.Context, reqInfo reqInfo) ([]byte, error) {
	deviceInfo, err := uc.vendorResolver.Get(ctx, reqInfo.Filename, reqInfo.UserAgent)
	if err != nil {
		return nil, err
	}
	settings, err := uc.settingSvc.Get(ctx, deviceInfo.MAC())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			deviceInfo.SetIP(reqInfo.IP)
			go uc.settingSvc.Create(deviceInfo)
			return nil, err
		}

		return nil, err
	}

	gen, err := uc.genFactory.Get(deviceInfo.Model(), deviceInfo.Vendor())
	if err != nil {
		return nil, err
	}
	cFile, err := gen.Generate(ctx, settings)
	if err != nil {
		return nil, err
	}

	return cFile, nil
}

// NewUseCase creates a new UseCase
// It takes a setting service, logger, generator factory, and vendor resolver as parameters
// Returns a pointer to the new UseCase
// parameters:
// - settingService: interface for the setting service to use for getting device settings
// - logger: interface for the logger to use for logging
// - generatorFactory: interface for the generator factory to use for getting a generator
// - vendorResolver: interface for the vendor resolver to use for resolving device vendor and model
func NewUseCase(
	settingService settingService,
	logger logger,
	generatorFactory GeneratorFactory,
	vendorResolver vendorResolver,
) *UseCase {
	return &UseCase{
		settingSvc:     settingService,
		logger:         logger,
		genFactory:     generatorFactory,
		vendorResolver: vendorResolver,
	}
}
