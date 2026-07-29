package identifier

import (
	"strings"
	"testing"
)

func TestGetVendorAndModel(t *testing.T) {
	tests := []struct {
		name       string
		ua         string
		wantModel  string
		wantVendor string
	}{
		{
			name:       "Yealink \"Yealink SIP-T46U 108.86.0.20\"",
			ua:         "Yealink SIP-T46U 108.86.0.20",
			wantModel:  "SIP-T46U",
			wantVendor: "Yealink",
		},
		{
			name:       "Grandstream \"Grandstream GXP2170 1.0.11.57\"",
			ua:         "Grandstream GXP2170 1.0.11.57",
			wantModel:  "GXP2170",
			wantVendor: "Grandstream",
		},
		{
			name:       "Fanvil \"Fanvil X4U 2.12.7.6\"",
			ua:         "Fanvil X4U 2.12.7.6",
			wantModel:  "X4U",
			wantVendor: "Fanvil",
		},
		{
			name:       "Cisco \"Cisco/SPA504G-7.6.2c\"",
			ua:         "Cisco/SPA504G-7.6.2c",
			wantModel:  "SPA504G",
			wantVendor: "cisco",
		},
		{
			name:       "snom \"User-Agent:snomD785-SIP 10.1.169.16\"",
			ua:         "snomD785-SIP 10.1.169.16",
			wantModel:  "D785",
			wantVendor: "snom",
		},
		{

			name:       "Polycom \"User-Agent: PolycomVVX-VVX_450-UA/6.4.5.1000\"",
			ua:         "PolycomVVX-VVX450-UA/6.4.5.1000",
			wantModel:  "VVX450",
			wantVendor: "Polycom",
		},
		{
			name:       "Htek  \"User-Agent: Htek UC924E V3.0.4.5\"",
			ua:         "Htek UC924E V3.0.4.5",
			wantModel:  "UC924E",
			wantVendor: "Htek",
		},
		{
			name:       "Unknown UA or NoName vendor",
			ua:         "User-Agent: NoName Unknown",
			wantModel:  "",
			wantVendor: "",
		},
		{
			name:       "Empty UserAgent",
			ua:         "",
			wantModel:  "",
			wantVendor: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vendor, model := getVendorAndModel(tt.ua)
			wantModel := strings.ToLower(tt.wantModel)
			wantVendor := strings.ToLower(tt.wantVendor)
			t.Logf("Received vedor: (%s), model: (%s)\n", vendor, model)

			if vendor != wantVendor {
				t.Errorf("vendor: got %q, want %q", vendor, wantVendor)
			}

			if model != wantModel {
				t.Errorf("model: got %q, want %q", model, wantModel)
			}
		})
	}

}
