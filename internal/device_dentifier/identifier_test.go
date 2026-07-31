package identifier

import (
	"context"
	"errors"
	"fmt"
	"github/sta-zot/pcs/internal/domain"
	"testing"
)

// Вданном случе интерфейс для клиентов для внешних API
type mockvendorChecker struct {
	vendor      string
	receivedMac string
	err         error
}

func (mc *mockvendorChecker) GetVendor(ctx context.Context, mac string) (string, error) {
	mc.receivedMac = mac
	return mc.vendor, mc.err
}

func TestIdnetifier(t *testing.T) {
	tests := []struct {
		name         string
		mockVendor   string
		mockErr      error
		wantErr      bool
		ua           string
		filename     string
		expectVendor string
		expectModel  string
		expectMAC    string
		expectErr    error
	}{
		{
			// Yealink
			name:         "Yealint with valid data",
			mockVendor:   "yealink",
			mockErr:      nil,
			wantErr:      false,
			ua:           "Yealink SIP-T46U 108.86.0.20",
			filename:     "001565aabbcc.cfg",
			expectMAC:    "001565aabbcc",
			expectVendor: "yealink",
			expectModel:  "sip-t46u",
			expectErr:    nil,
		},
		{
			name:         "Yealink with invalid data",
			mockVendor:   "",
			mockErr:      errors.New("no data found"),
			wantErr:      true,
			ua:           "Noname SIP-T46U 108.86.0.20",
			filename:     "001565aabbcc.cfg",
			expectMAC:    "001565aabbcc",
			expectVendor: "",
			expectModel:  "",
			expectErr:    domain.ErrVendorNotFound,
		},
		{
			name:         "Cisco with valid data",
			mockVendor:   "cisco",
			mockErr:      nil,
			wantErr:      false,
			ua:           "Cisco/SPA504G-7.6.2c",
			filename:     "spa001122AABBCC.xml",
			expectMAC:    "001122aabbcc",
			expectVendor: "cisco",
			expectModel:  "spa504g",
			expectErr:    nil,
		},
		{
			name:         "Masmached vendors",
			mockVendor:   "cisco",
			mockErr:      nil,
			wantErr:      true,
			ua:           "Yealink SIP-T46U 108.86.0.20",
			filename:     "spa001122AABBCC.xml",
			expectMAC:    "001122aabbcc",
			expectVendor: "cisco",
			expectModel:  "spa504g",
			expectErr:    domain.ErrVendorMismatch,
		},
		{
			name:         "Cisco with valid data but OUI API not available",
			mockVendor:   "",
			mockErr:      errors.New("all providers failed"),
			wantErr:      false,
			ua:           "Cisco/SPA504G-7.6.2c",
			filename:     "spa001122AABBCC.xml",
			expectMAC:    "001122aabbcc",
			expectVendor: "cisco",
			expectModel:  "spa504g",
		},
	}
	for idx, tt := range tests {
		testName := fmt.Sprintf("%d. %s", idx, tt.name)
		t.Run(testName, func(t *testing.T) {
			vendorChecker := &mockvendorChecker{
				vendor: tt.mockVendor,
				err:    tt.mockErr,
			}
			devIdentifier := NewDeviceIdentifier(vendorChecker)
			devInfo, err := devIdentifier.Identify(context.Background(), tt.ua, tt.filename)
			if tt.wantErr && err == nil {
				t.Fatalf("Expected error (%s) but received nil", tt.expectErr.Error())
			}
			if err != nil && !tt.wantErr {
				t.Fatalf("Received unexpected error (%s) ", err.Error())
			}
			if tt.wantErr {
				if err != nil {
					if errors.Is(err, tt.expectErr) {
						t.Logf("Recieved expected error (%s)", err.Error())
						return
					} else {

						t.Fatalf("Recieved error(%s) missmatch expected (%s)", err.Error(), tt.expectErr.Error())
					}
				} else {
					t.Logf("Expected error (%s), but rrceived nil ", tt.expectErr.Error())
				}
			}
			if devInfo.MAC() != tt.expectMAC {
				t.Errorf("MAC addresses is missmatch. Received (%s) - Expected (%s)", devInfo.MAC(), tt.expectMAC)
			}
			if devInfo.Model() != tt.expectModel {
				t.Errorf("Models is missmatch. Received (%s) - Expected (%s)", devInfo.Model(), tt.expectModel)
			}
			if devInfo.Vendor() != tt.expectVendor {
				t.Errorf("Vendors is missmatch. Received (%s) - Expected (%s)", devInfo.Vendor(), tt.expectVendor)
			}
		})

	}
}
