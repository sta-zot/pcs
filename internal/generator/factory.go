package generator

import (
	"fmt"
	"github/sta-zot/pcs/internal/domain"
)

var supportedVendors = map[string]Generator{
	"yealink": NewYealinkGenerator(),
	// "cisco":       "cisco",
	// "grandstream": "grandstream",
	// "fanvil":      "fanvil",
	// "snom":        "snom",
	// "polycom":     "polycom",
	// "htek":        "htek",
}
var supportedModels = map[string]Generator{
	// "spa": "spa",
}

func generatorFactory(vendor, model string) (Generator, error) {
	var generator Generator
	if gen, ok := supportedVendors[vendor]; ok {
		generator = gen
		if gen, ok := supportedModels[model]; ok {
			generator = gen
		}
		return generator, nil
	}
	return nil, fmt.Errorf("%w. vendor(%s), model(%s)", domain.ErrPhoneNotSupported, vendor, model)
}
