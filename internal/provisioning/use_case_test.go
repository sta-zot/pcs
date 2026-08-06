package provisioning

import (
	"context"
	"errors"
	"fmt"
	"github/sta-zot/pcs/internal/domain"
	"io"
	"strings"
	"testing"
)

//mocks

// MOCK for settings provider
type mSettingsProvider struct {
	receivedMac      string
	receivedSettings *domain.PhoneSettings
	pSettings        *domain.PhoneSettings
	called           bool
	err              error
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
	called       bool
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
	buf              io.ByteReader
	called           bool
}

func (mcg *mConfigGenerator) Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error) {
	mcg.receivedSettings = settings

	return mcg.buf, mcg.err
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

	createMocks := func() (*mVendorIdentifier, *mSettingsProvider, *mConfigGenerator, *mLogger) {
		return &mVendorIdentifier{called: true},
			&mSettingsProvider{called: true},
			&mConfigGenerator{called: true},
			&mLogger{}
	}
	_, _, _, _ = createMocks()

	// Позже удалить, поставил что бы убрать ошибки
	// #############################################

	_ = rInfo
	//##############################################
	buf := strings.NewReader("Test file entry")
	vendor := "yealink"
	model := "sip-t46u"
	mac := "112233aabbcc"

	tPhoneSettings := domain.NewPhoneSettings(*domain.NewDeviceInfo(model, vendor, mac, rInfo.IP))

	tests := []struct {
		name                       string
		identifierErr              error
		identifierVendor           string
		identifierModel            string
		identifierMAC              string
		identifierReceivedUA       string
		identifierReceivedFileName string

		sProviderErr     error
		providerSettings *domain.PhoneSettings

		confGenErr error
		confGenBuf io.ByteReader

		wantErr     bool
		expectedErr error
	}{
		{
			name:    "Полностью валидный",
			wantErr: false,

			identifierErr:              nil,
			identifierVendor:           vendor,
			identifierModel:            model,
			identifierMAC:              mac,
			identifierReceivedFileName: rInfo.Filename,
			identifierReceivedUA:       rInfo.UserAgent,
			providerSettings:           tPhoneSettings,
			sProviderErr:               nil,
			confGenErr:                 nil,
			confGenBuf:                 buf,
		},
		{
			name:          "Ошибка идентификации устройства",
			wantErr:       true,
			expectedErr:   domain.ErrVendorNotFound,
			identifierErr: domain.ErrVendorNotFound,

			identifierVendor:           vendor,
			identifierModel:            model,
			identifierMAC:              mac,
			identifierReceivedFileName: rInfo.Filename,
			identifierReceivedUA:       rInfo.UserAgent,
			providerSettings:           tPhoneSettings,
			sProviderErr:               nil,
			confGenErr:                 nil,
			confGenBuf:                 buf,
		},
		{
			name:        "Конфигурации нет в БД. Первое обращение",
			wantErr:     true,
			expectedErr: domain.ErrSettingsNotFound,

			identifierErr:              nil,
			identifierVendor:           vendor,
			identifierModel:            model,
			identifierMAC:              mac,
			identifierReceivedFileName: rInfo.Filename,
			identifierReceivedUA:       rInfo.UserAgent,
			providerSettings:           tPhoneSettings,
			sProviderErr:               domain.ErrSettingsNotFound,
			confGenErr:                 nil,
			confGenBuf:                 buf,
		},
		{
			name:        "Конфигурация есть, но не полная",
			wantErr:     true,
			expectedErr: domain.ErrConfigIncomplete,

			identifierErr:              nil,
			identifierVendor:           vendor,
			identifierModel:            model,
			identifierMAC:              mac,
			identifierReceivedFileName: rInfo.Filename,
			identifierReceivedUA:       rInfo.UserAgent,
			providerSettings:           tPhoneSettings,
			sProviderErr:               domain.ErrConfigIncomplete,
			confGenErr:                 nil,
			confGenBuf:                 buf,
		},
		{
			name:        "Щшибка генерации файла",
			wantErr:     true,
			expectedErr: domain.ErrFileGenerationError,

			identifierErr:              nil,
			identifierVendor:           vendor,
			identifierModel:            model,
			identifierMAC:              mac,
			identifierReceivedFileName: rInfo.Filename,
			identifierReceivedUA:       rInfo.UserAgent,
			providerSettings:           tPhoneSettings,
			sProviderErr:               nil,
			confGenErr:                 domain.ErrFileGenerationError,
			confGenBuf:                 buf,
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%d. %s", idx, tt.name),
			func(t *testing.T) {
				// Confirure vendor identifier
				vIdentifier := &mVendorIdentifier{
					err:    tt.identifierErr,
					vendor: tt.identifierVendor,
					model:  tt.identifierModel,
				}
				// Confirure setting provider

				sProvider := &mSettingsProvider{
					pSettings: tPhoneSettings,
					err:       tt.sProviderErr,
				}
				// Confirure config file generator
				cGenerator := &mConfigGenerator{
					err: tt.confGenErr,
					buf: tt.confGenBuf,
				}
				logger := &mLogger{}

				uc := NewUseCase(sProvider, vIdentifier, cGenerator, logger)
				file, err := uc.Provision(context.Background(), rInfo)

				if tt.wantErr {
					if err == nil {
						t.Errorf("Expected error(%s), but received nil", tt.expectedErr.Error())
					}
					if errors.Is(err, tt.expectedErr) {
						t.Logf("Received error(%s) macth to expected error(%s)", err.Error(), tt.expectedErr.Error())
						return
					} else {
						t.Fatalf("Received error(%s) mismacth to expected error(%s)", err.Error(), tt.expectedErr.Error())
					}
					if errors.Is(err, domain.ErrSettingsNotFound) {
						if !sProvider.called {
							t.Errorf("Create function not called")
						}
					}
					if vIdentifier.receivedFile != rInfo.Filename {
						t.Errorf("Received file name(%s) does not match transmited(%s)", rInfo.Filename, vIdentifier.receivedFile)
					}
					if vIdentifier.receivedUA != rInfo.UserAgent {
						t.Errorf("Received UserAgent(%s) does not match transmited(%s)", rInfo.UserAgent, vIdentifier.receivedUA)
					}

					if tt.providerSettings != sProvider.receivedSettings {
						t.Errorf("Received settings(%v) does not match transmited(%v)", tt.providerSettings, sProvider.receivedSettings)
					}

					expectedFile := stringFromByteReade(tt.confGenBuf)
					if expectedFile == "" {
						t.Fatalf("Буфер с данными для файла не удалось преобразовать в строку")
					}
					receivedFile := stringFromByteReade(file)

					if len(receivedFile) <= 0 {
						t.Errorf("Received file is empty, but expected\n[%s]", string(expectedFile))
					}

					if expectedFile != receivedFile {
						t.Errorf("Received file (%s) does not match expected(%s)", string(receivedFile), string(expectedFile))
					}

				}
			})
	}
}

func stringFromByteReade(br io.ByteReader) string {
	var data []byte

	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ""
		}
		data = append(data, b)
	}
	return string(data)

}
