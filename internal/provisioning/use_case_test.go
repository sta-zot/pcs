package provisioning

import (
	"context"
	"fmt"
	"github/sta-zot/pcs/internal/domain"
	"io"
	"strings"
	"testing"
)

//mocks

// MOCK for settings provider
type mSettingsProvider struct {
	receivedMac string
	pSettings   *domain.PhoneSettings
	called      bool
	err         error
}

func (msp *mSettingsProvider) Get(ctx context.Context, mac string) (*domain.PhoneSettings, error) {
	return msp.pSettings, msp.err
}
func (msp *mSettingsProvider) Create(ctx context.Context, settings *domain.PhoneSettings) error {
	msp.called = true
	return msp.err
}

type mVendorIdentifier struct {
	vendor       string
	model        string
	err          error
	receivedFile string
	receivedUA   string
}

func (mvi *mVendorIdentifier) Identify(ctx context.Context, filename, userAgent string) (*domain.DeviceInfo, error) {
	mvi.receivedFile = filename
	mvi.receivedUA = userAgent
	return domain.NewDeviceInfo(mvi.model, mvi.vendor, "112233aabbcc", ""),
		mvi.err

}

// MOCK for Config generator

type mConfigGenerator struct {
	receivedSettings *domain.PhoneSettings
	err              error
	buf              string
}

func (mcg *mConfigGenerator) Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error) {
	mcg.receivedSettings = settings
	reader := strings.NewReader(mcg.buf)
	return reader, mcg.err
}

type mLogger struct{}

func (mLogger) Info(msg string, args ...any)  {}
func (mLogger) Error(msg string, args ...any) {}
func (mLogger) Debug(msg string, args ...any) {}
func (mLogger) Warn(msg string, args ...any)  {}

func TestUseCase(t *testing.T) {
	rInfo := reqInfo{
		IP:        "10.0.0.1",
		UserAgent: "Yealink SIP-T46U 108.86.0.20",
		Filename:  "112233aabbcc.cfg",
	}
	createMocks := func() (vendorIdentifier, settingsProvider, configGenerator, logger) {
		return &mVendorIdentifier{}, &mSettingsProvider{}, &mConfigGenerator{}, &mLogger{}
	}
	_, _, _, _ = createMocks()
	tests := []struct {
		name         string
		idetifierErr error
		sProviderErr error
		confGenErr   error
	}{}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%d. %s", idx, tt.name),
			func(t testing.T))
	}
}
