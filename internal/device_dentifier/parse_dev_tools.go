package identifier

import (
	"regexp"
	"strings"
)

var (
	// Yealink: 001565aabbcc.cfg  или  001565aabbcc.y000000000000.cfg
	reYealink = regexp.MustCompile(`^([0-9a-fA-F]{12})(?:\.y[0-9A-F]+)?\.(?:cfg|xml)$`)
	// Cisco SPA: spa001122AABBCC.xml
	reCiscoSPA = regexp.MustCompile(`^spa([0-9a-fA-F]{12})\.xml$`)
	// Cisco 79xx: SEP001122AABBCC.cnf.xml
	reCisco79 = regexp.MustCompile(`^SEP([0-9a-fA-F]{12})\.cnf\.xml$`)
	// Grandstream: cfg000B82AABBCC.xml  или  fp000B82AABBCC.xml
	reGrandstream = regexp.MustCompile(`^(?:cfg|fp)([0-9a-fA-F]{12})\.xml$`)
	macPattern    = regexp.MustCompile(`(([a-fA-F0-9]{2}[-:. ]){5}[a-fA-F0-9]{2}|[a-fA-F0-9]{12})`)
)

func getMacAddr(fileName string) (string, error) {
	mac := ""
	mac = macPattern.FindString(fileName)
	if mac == "" {
		return "", nil
	}
	regexDelim := regexp.MustCompile(`[-.:\s]`)
	mac = regexDelim.ReplaceAllString(mac, "")
	mac = strings.ToLower(mac)

	return mac, nil
}

func getVendorModel(userAgent, fileName string) (string, string) {
	vendor := ""
	model := ""

	return vendor, model
}
