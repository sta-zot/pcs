package identifier

const (
	VendorYealink     = "yealink"
	VendorCisco       = "cisco"
	VendorGrandstream = "grandstream"
	VendorFanvil      = "fanvil"
	VendorSnom        = "snom"
	VendorHtek        = "htek"
	VendorPoly        = "polycom"
	VendorAudioCodes  = "audiocodes"
	VendorFlyingvoice = "flyingvoice"
	VendorUnknown     = ""
)

var vendorPrefixes = map[string]string{
	"yealink":     VendorYealink,
	"cisco":       VendorCisco,
	"grandstream": VendorGrandstream,
	"fanvil":      VendorFanvil,
	"snom":        VendorSnom,
	"htek":        VendorHtek,
	"poly":        VendorPoly,
	"polycom":     VendorPoly,
	"audiocodes":  VendorAudioCodes,
	"flyingvoice": VendorFlyingvoice,
}
