package identifier

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Yealink: 001565aabbcc.cfg  или  001565aabbcc.y000000000000.cfg
	reYealinkFile = regexp.MustCompile(`^([0-9a-fA-F]{12})(?:\.y[0-9A-F]+)?\.(?:cfg|xml)$`)
	// Cisco SPA: spa001122AABBCC.xml
	reCiscoSPAFile = regexp.MustCompile(`^spa([0-9a-fA-F]{12})\.xml$`)
	// Cisco 79xx: SEP001122AABBCC.cnf.xml
	reCisco79File = regexp.MustCompile(`^SEP([0-9a-fA-F]{12})\.cnf\.xml$`)
	// Grandstream: cfg000B82AABBCC.xml  или  fp000B82AABBCC.xml
	reGrandstreamFile = regexp.MustCompile(`^(?:cfg|fp)([0-9a-fA-F]{12})\.xml$`)
	macPattern        = regexp.MustCompile(`(([a-fA-F0-9]{2}[-:. ]){5}[a-fA-F0-9]{2}|[a-fA-F0-9]{12})`)
)

func normalizeMacAddr(fileName string) (string, error) {
	mac := ""
	mac = macPattern.FindString(fileName)
	if mac == "" {
		return "", fmt.Errorf("MAC address not found")
	}
	regexDelim := regexp.MustCompile(`[-.:\s]`)
	mac = regexDelim.ReplaceAllString(mac, "")
	mac = strings.ToLower(mac)

	return mac, nil
}

func matchVendor(ua, vendor string, parser func(string) (string, bool)) (string, string) {
	if model, ok := parser(ua); ok {
		return vendor, model
	}
	return "", ""
}
func getVendorAndModel(ua string) (vendor, model string) {
	ua = strings.TrimSpace(ua)
	ua = strings.ToLower(ua)

	switch {
	case strings.HasPrefix(ua, "yealink"):
		return matchVendor(ua, "yealink", parseYealink)

	case strings.HasPrefix(ua, "cisco/"):
		return matchVendor(ua, "cisco", parseCisco)

	case strings.HasPrefix(ua, "grandstream"):
		return matchVendor(ua, "grandstream", parseGrandstream)

	case strings.HasPrefix(ua, "fanvil"):
		return matchVendor(ua, "fanvil", parseFanvil)

	case strings.HasPrefix(ua, "snom"):
		return matchVendor(ua, "snom", parseSnom)
	case strings.HasPrefix(ua, "polycom"):
		return matchVendor(ua, "polycom", parsePoly)

	case strings.HasPrefix(ua, "poly"):
		return matchVendor(ua, "polycom", parsePoly)

	case strings.HasPrefix(ua, "htek"):
		return matchVendor(ua, "htek", parseHTek)

	case strings.HasPrefix(ua, "flyingvoice"):
		return matchVendor(ua, "flyingvoice", parseFlyingvoice)
	}

	return "", ""
}

// Yealink parser
var reYealink = regexp.MustCompile(`(?i)^Yealink\s+([A-Za-z0-9-]+)`)

func parseYealink(ua string) (string, bool) {
	m := reYealink.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

// Cisco SPA parser
var reCisco = regexp.MustCompile(`(?i)^Cisco/([A-Za-z0-9]+)`)

func parseCisco(ua string) (string, bool) {
	m := reCisco.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

// Grandstream parser

var reGrandstream = regexp.MustCompile(`(?i)^Grandstream\s+([A-Za-z0-9-]+)`)

func parseGrandstream(ua string) (string, bool) {
	m := reGrandstream.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

// Fanvil parser
var reFanvil = regexp.MustCompile(`(?i)^Fanvil\s+([A-Za-z0-9-]+)`)

func parseFanvil(ua string) (string, bool) {
	m := reFanvil.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

//Snom Parser

var reSnom = regexp.MustCompile(`(?i)^snom([A-Za-z0-9]+)`)

func parseSnom(ua string) (string, bool) {
	m := reSnom.FindStringSubmatch(ua)

	if len(m) < 1 {
		return "", false
	}
	return m[1], true
}

// Poly or Polycom parser

var polyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^PolycomVVX-([A-Za-z0-9_]+)-UA/`),
	regexp.MustCompile(`(?i)^PolyTrio-([A-Za-z0-9_]+)-UA/`),
	regexp.MustCompile(`(?i)^SoundPointIP-([A-Za-z0-9_]+)-UA/`),
}

func parsePoly(ua string) (string, bool) {
	for _, re := range polyPatterns {
		if m := re.FindStringSubmatch(ua); len(m) == 2 {
			return normalizeModel(m[1]), true
		}
	}

	return "", false
}

var excludedChars = regexp.MustCompile(`(?i)(:?[_\ ])`)

func normalizeModel(m string) string {

	m = excludedChars.ReplaceAllLiteralString(m, "")
	return m

}

// Flyingvoice parser

var reFlying = regexp.MustCompile(`(?i)^Flyingvoice\s+([A-Za-z0-9-]+)`)

func parseFlyingvoice(ua string) (string, bool) {
	fmt.Println(ua)
	m := reFlying.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

// HTek parser
var reHTek = regexp.MustCompile(`(?i)^Htek\s+([A-Za-z0-9-]+)`)

func parseHTek(ua string) (string, bool) {
	m := reHTek.FindStringSubmatch(ua)
	if len(m) != 2 {
		return "", false
	}
	return m[1], true
}

func normalizeVendor(v string) string {
	nv := strings.Fields(strings.ToLower(v))
	if len(nv) == 1 {
		return ""
	}
	return nv[0]
}
