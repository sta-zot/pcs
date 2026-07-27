package provisioning

import (
	"context"
	"errors"
	"github/sta-zot/pcs/internal/domain"
	"time"

	"testing"
)

// MockVendorResolver================================
type MockVendorResolver struct {
	ReceivedFilename  string
	ReceivedUserAgent string
	DeviceInfo        *domain.DeviceInfo
	Error             error
}

func (mvr *MockVendorResolver) Get(ctx context.Context, filename, userAgent string) (*domain.DeviceInfo, error) {
	mvr.ReceivedFilename = filename
	mvr.ReceivedUserAgent = userAgent
	return mvr.DeviceInfo, mvr.Error
}

//================================================================================

// MockGenerator =================================
type MockGenerator struct {
	ReceivedSettings *domain.PhoneSettings
	GeneratedContent []byte
	Error            error
}

func (mg *MockGenerator) Generate(ctx context.Context, settings *domain.PhoneSettings) ([]byte, error) {
	mg.ReceivedSettings = settings
	return mg.GeneratedContent, mg.Error
}

//================================================================================

// MockLogger =================================
type MockLogger struct {
	InfoMsg  string
	ErrorMsg string
	DebugMsg string
	WarnMsg  string
}

func (ml *MockLogger) Info(msg string, args ...any) {
	ml.InfoMsg = msg
}

func (ml *MockLogger) Error(msg string, args ...any) {
	ml.ErrorMsg = msg
}

func (ml *MockLogger) Debug(msg string, args ...any) {
	ml.DebugMsg = msg
}

func (ml *MockLogger) Warn(msg string, args ...any) {
	ml.WarnMsg = msg
}

//================================================================================

// MockGeneratorFactory =================================
type MockGeneratorFactory struct {
	Generator         *MockGenerator
	Error             error
	ReceivedDevModel  string
	ReceivedDevVendor string
}

func (mgf *MockGeneratorFactory) Get(model, vendor string) (Generator, error) {
	mgf.ReceivedDevModel = model
	mgf.ReceivedDevVendor = vendor
	return mgf.Generator, mgf.Error
}

//================================================================================

// MockSettingService =================================
type MockSettingService struct {
	ReceivedDeviceInfo *domain.DeviceInfo
	ReceivedMacAddress string
	Error              error
	PSettings          *domain.PhoneSettings
}

func (mss *MockSettingService) Get(ctx context.Context, mac string) (*domain.PhoneSettings, error) {
	mss.ReceivedMacAddress = mac
	return mss.PSettings, mss.Error
}

func (mss *MockSettingService) Create(deviceInfo *domain.DeviceInfo) error {
	mss.ReceivedDeviceInfo = deviceInfo
	return mss.Error
}

var (
	_ = vendorResolver(&MockVendorResolver{})
	_ = settingService(&MockSettingService{})
	_ = GeneratorFactory(&MockGeneratorFactory{})
	_ = logger(&MockLogger{})
	_ = Generator(&MockGenerator{})
)

func TestProvision(t *testing.T) {
	testMac := "aabbccddeeff"
	testFileName := "aabbccddeeff.cfg"
	testUserAgent := "Yealink/t34"
	testDeviceInfo := domain.NewDeviceInfo("t34", "Yealink", testMac, "1.2.3.4")
	testPhoneSettings := domain.NewPhoneSettings(*testDeviceInfo)
	testReqInfo := reqInfo{
		UserAgent: testUserAgent,
		Filename:  testFileName,
		IP:        "1.2.3.4",
	}

	tests := []struct {
		name         string
		resolverErr  error
		generatorErr error
		factoryErr   error
		settingsErr  error
		wantErr      error

		// Явные ожидания: какие методы ДОЛЖНЫ быть вызваны в этом сценарии
		expectResolverCall  bool
		expectSettingGet    bool
		expectSettingCreate bool
		expectFactoryCall   bool
	}{
		{
			name:                "success",
			wantErr:             nil,
			expectResolverCall:  true,
			expectSettingGet:    true,
			expectFactoryCall:   true,
			expectSettingCreate: false,
		},
		{
			name:                "resolver error",
			resolverErr:         domain.ErrVendorResolvingError,
			wantErr:             domain.ErrVendorResolvingError,
			expectResolverCall:  true,
			expectSettingGet:    false, // Прервалось раньше, не вызываем
			expectFactoryCall:   false, // Прервалось раньше, не вызываем
			expectSettingCreate: false,
		},
		{
			name:                "generator error",
			generatorErr:        domain.ErrFileGenerationError,
			wantErr:             domain.ErrFileGenerationError,
			expectResolverCall:  true,
			expectSettingGet:    true,
			expectFactoryCall:   true,
			expectSettingCreate: false,
		},
		{
			name:                "Generator factory error",
			factoryErr:          domain.ErrFileGeneratorNotExist,
			wantErr:             domain.ErrFileGeneratorNotExist,
			expectResolverCall:  true,
			expectSettingGet:    true,
			expectFactoryCall:   true,
			expectSettingCreate: false,
		},
		{
			name:                "Settings not found error",
			settingsErr:         domain.ErrNotFound,
			wantErr:             domain.ErrNotFound,
			expectResolverCall:  true,
			expectSettingGet:    true,
			expectSettingCreate: true,  // Должен вызваться асинхронно!
			expectFactoryCall:   false, // Прервалось после Create
		},
		{
			name:                "Settings not complete",
			settingsErr:         domain.ErrConfigIncomplete,
			wantErr:             domain.ErrConfigIncomplete,
			expectResolverCall:  true,
			expectSettingGet:    true,
			expectSettingCreate: false, // Ошибка не ErrNotFound, Create не вызывается
			expectFactoryCall:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Инициализация моков
			resolver := &MockVendorResolver{
				DeviceInfo: testDeviceInfo,
				Error:      tc.resolverErr,
			}
			settingService := &MockSettingService{
				PSettings: testPhoneSettings,
				Error:     tc.settingsErr,
			}
			generatorFactory := &MockGeneratorFactory{
				Generator: &MockGenerator{
					GeneratedContent: []byte("File content!"),
					Error:            tc.generatorErr,
				},
				Error: tc.factoryErr,
			}
			logger := &MockLogger{}

			useCase := NewUseCase(settingService, logger, generatorFactory, resolver)

			// 2. Вызов
			result, err := useCase.Provision(context.Background(), testReqInfo)

			// 3. Проверка ошибки
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error: %v, got: %v", tc.wantErr, err)
			}

			// 4. Проверка результата (ИСПРАВЛЕНО: проверяем наличие результата ТОЛЬКО при успехе)
			if tc.wantErr == nil && len(result) == 0 {
				t.Errorf("expected result on success, got empty")
			}

			// 5. Проверка взаимодействий (ТОЛЬКО если ожидаем вызов)

			if tc.expectResolverCall {
				if resolver.ReceivedFilename != testReqInfo.Filename {
					t.Errorf("expected filename: %s, got: %s", testReqInfo.Filename, resolver.ReceivedFilename)
				}
				if resolver.ReceivedUserAgent != testReqInfo.UserAgent {
					t.Errorf("expected user agent: %s, got: %s", testReqInfo.UserAgent, resolver.ReceivedUserAgent)
				}
			}

			if tc.expectSettingGet {
				if settingService.ReceivedMacAddress != testMac {
					t.Errorf("expected mac address: %s, got: %s", testMac, settingService.ReceivedMacAddress)
				}
			}

			if tc.expectSettingCreate {
				// ИСПРАВЛЕНО: 5 секунд - это вечность для unit-теста. 10-50 мс достаточно для планировщика Go.
				time.Sleep(50 * time.Millisecond)

				if settingService.ReceivedDeviceInfo == nil {
					t.Errorf("expected Create to be called with device info, got nil")
				} else if settingService.ReceivedDeviceInfo.IP() != testReqInfo.IP {
					// ИСПРАВЛЕНО: Сравнивать указатели (!= testDeviceInfo) опасно, лучше сравнивать поля
					t.Errorf("expected device info IP: %s, got: %s", testReqInfo.IP, settingService.ReceivedDeviceInfo.IP())
				}
			}

			if tc.expectFactoryCall {
				if generatorFactory.ReceivedDevModel != testDeviceInfo.Model() {
					t.Errorf("expected dev model: %s, got: %s", testDeviceInfo.Model(), generatorFactory.ReceivedDevModel)
				}
				if generatorFactory.ReceivedDevVendor != testDeviceInfo.Vendor() {
					t.Errorf("expected dev vendor: %s, got: %s", testDeviceInfo.Vendor(), generatorFactory.ReceivedDevVendor)
				}
			}
		})
	}
}
