package domain

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ipPattern = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

	macPattern = regexp.MustCompile(`^(([a-fA-F0-9]{2}[-:. ]){5}[a-fA-F0-9]{2}|[a-fA-F0-9]{12})$`)

	dnPattern = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}$`)
)

// Validate ip addresses
func validateIP(ip string) (string, error) {
	if ip == "" {
		return "", errors.New("Invalid IP address")
	}
	if s := ipPattern.FindString(ip); s != "" {
		parts := strings.Split(s, ".")
		for i, part := range parts {
			num, _ := strconv.Atoi(part)
			parts[i] = strconv.Itoa(num)
		}
		validIP := strings.Join(parts, ".")
		return validIP, nil
	}
	return "", errors.New("Invalid IP address")
}

func validateMAC(mac string) bool {
	if mac == "" || !macPattern.MatchString(mac) {
		return false
	}
	return true
}

func validateDN(dn string) bool {
	if dn == "" || !dnPattern.MatchString(dn) {
		return false
	}
	return true
}
