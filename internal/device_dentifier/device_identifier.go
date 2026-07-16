package identifier

import "context"

type macChecker interface {
	GetVendor(ctx context.Context, mac string) (string, error)
}

type deviceIdentifier struct {
	macChecker macChecker
}

func NewDeviceIdentifier(macChecker macChecker) *deviceIdentifier {
	return &deviceIdentifier{
		macChecker: macChecker,
	}
}
