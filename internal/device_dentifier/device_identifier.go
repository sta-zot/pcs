package identifier

import (
	"context"
	"fmt"
	"github/sta-zot/pcs/internal/domain"
	"strings"
)

type vendorChecker interface {
	GetVendor(ctx context.Context, mac string) (string, error)
}

type deviceIdentifier struct {
	vendorChecker vendorChecker
}

func NewDeviceIdentifier(vendorChecker vendorChecker) *deviceIdentifier {
	return &deviceIdentifier{
		vendorChecker: vendorChecker,
	}
}

func (d *deviceIdentifier) Identify(ctx context.Context, userAgent, fileName string) (*domain.DeviceInfo, error) {

	normalMAC, err := normalizeMacAddr(fileName)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", domain.ErrInvalidMac, err)
	}
	ouiVendor, err := d.vendorChecker.GetVendor(ctx, normalMAC)
	if err != nil {
		return nil, err
	}

	uaVendor, model := getVendorAndModel(userAgent)
	ouiVendor = strings.Fields(strings.ToLower(ouiVendor))[0]
	uaVendor = strings.Fields(strings.ToLower(uaVendor))[0]
	model = strings.ToLower(model)
	normalMAC = strings.ToLower(normalMAC)
	if uaVendor == "" && ouiVendor == "" {
		return nil, fmt.Errorf("%w. UserAgent: %s", domain.ErrVendorNotFound, userAgent)
	}
	if (uaVendor != "") && (ouiVendor != "") && (ouiVendor != uaVendor) {
		return nil, fmt.Errorf("%w user-agent (%s) does not match OUI (%s)", domain.ErrVendorMismatch, uaVendor, ouiVendor)
	}

	return domain.NewDeviceInfo(model, ouiVendor, normalMAC, ""), nil

}
