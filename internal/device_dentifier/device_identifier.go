package identifier

import (
	"context"
	"fmt"
	"github/sta-zot/pcs/internal/domain"
	"log"
	"strings"
)

type Logger interface {
	Info(msg string)
}

type vendorChecker interface {
	GetVendor(ctx context.Context, mac string) (string, error)
}

type deviceIdentifier struct {
	vendorChecker vendorChecker
	loger         Logger
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
		log.Println(err)
	}

	uaVendor, model := getVendorAndModel(userAgent)
	ouiVendor = normalizeVendor(ouiVendor)
	uaVendor = normalizeVendor(uaVendor)
	model = strings.ToLower(model)
	normalMAC = strings.ToLower(normalMAC)
	fmt.Println("uaVendor: " + uaVendor)
	fmt.Println("ouiVendor: " + ouiVendor)
	if uaVendor == "" && ouiVendor == "" {
		return nil, fmt.Errorf("%w. UserAgent: %s", domain.ErrVendorNotFound, userAgent)
	}
	if (uaVendor != "") && (ouiVendor != "") && (ouiVendor != uaVendor) {
		return nil, fmt.Errorf("%w user-agent (%s) does not match OUI (%s)", domain.ErrVendorMismatch, uaVendor, ouiVendor)
	}
	if ouiVendor == "" {
		ouiVendor = uaVendor
	}
	return domain.NewDeviceInfo(model, ouiVendor, normalMAC, ""), nil

}
