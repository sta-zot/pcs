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
	msp.receivedSettings = settings
	return msp.err
}

// MOCK for vendor identifier
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
	return domain.NewDeviceInfo(mvi.model, mvi.vendor, "112233aabbcc", ""), mvi.err
}

// MOCK for Config generator
type mConfigGenerator struct {
	receivedSettings *domain.PhoneSettings
	err              error
	content          string // вместо buf храним строку
}

func (mcg *mConfigGenerator) Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error) {
	mcg.receivedSettings = settings
	if mcg.err != nil {
		return nil, mcg.err
	}
	return strings.NewReader(mcg.content), nil
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

	const expectedContent = "Test file entry"
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

		confGenErr      error
		expectedContent string // ожидаемое содержимое файла

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
			expectedContent:            expectedContent,
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
			expectedContent:            expectedContent,
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
			expectedContent:            expectedContent,
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
			expectedContent:            expectedContent,
		},
		{
			name:        "Ошибка генерации файла",
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
			expectedContent:            expectedContent,
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("#%d. %s", idx+1, tt.name),
			func(t *testing.T) {
				vIdentifier := &mVendorIdentifier{
					err:    tt.identifierErr,
					vendor: tt.identifierVendor,
					model:  tt.identifierModel,
				}
				sProvider := &mSettingsProvider{
					pSettings: tt.providerSettings,
					err:       tt.sProviderErr,
				}
				cGenerator := &mConfigGenerator{
					err:     tt.confGenErr,
					content: tt.expectedContent,
				}
				logger := &mLogger{}

				uc := NewUseCase(sProvider, vIdentifier, cGenerator, logger)
				file, err := uc.Provision(context.Background(), rInfo)

				if tt.wantErr {
					if err == nil {
						t.Errorf("Expected error(%s), but received nil", tt.expectedErr.Error())
					}
					if !errors.Is(err, tt.expectedErr) {
						t.Fatalf("Received error(%s) mismatch expected error(%s)", err.Error(), tt.expectedErr.Error())
					}
					if errors.Is(err, domain.ErrSettingsNotFound) {
						t.Logf("Received error(%s)", err.Error())
						if !sProvider.called {
							t.Errorf("Create function not called")
						}
						if tt.providerSettings == nil || sProvider.receivedSettings == nil {
							t.Errorf("Received settings(%v) does not match transmitted(%v)", tt.providerSettings, sProvider.receivedSettings)
						}
					}
					// Для ошибочных кейсов файл должен быть nil, проверять содержимое не нужно
					if file != nil {
						t.Errorf("Expected file to be nil, but got %v", file)
					}
					return // выходим, остальные проверки не актуальны
				}

				// Успешный случай
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if file == nil {
					t.Fatal("File is nil")
				}

				// Проверяем, что в vendor identifier пришли правильные параметры
				if vIdentifier.receivedFile != rInfo.Filename {
					t.Errorf("Received file name(%s) does not match transmitted(%s)", rInfo.Filename, vIdentifier.receivedFile)
				}
				if vIdentifier.receivedUA != rInfo.UserAgent {
					t.Errorf("Received UserAgent(%s) does not match transmitted(%s)", rInfo.UserAgent, vIdentifier.receivedUA)
				}

				// Проверяем содержимое файла
				received := stringFromByteReader(file)
				if received != tt.expectedContent {
					t.Errorf("File content mismatch:\nexpected:\n%s\ngot:\n%s", tt.expectedContent, received)
				}
			})
	}
}

func stringFromByteReader(br io.ByteReader) string {
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
