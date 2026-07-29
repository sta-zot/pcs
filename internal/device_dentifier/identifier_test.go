package identifier

import "context"

type mockOUIClient struct {
	vendor string
	err    error
}

func (mc *mockOUIClient) GetVendor(ctx context.Context, mac string) (string, error) {
	return mc.vendor, mc.err
}
