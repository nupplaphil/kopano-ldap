package utils

import (
	"bytes"
	"strings"
)

// Returns the DN string based on a fully qualified domain name
func GetBaseDN(fqdn string) string {
	parts := strings.Split(fqdn, ".")

	var b bytes.Buffer

	for i := range parts {
		if len(parts[i]) == 0 {
			continue
		}
		b.WriteString("dc=")
		b.WriteString(parts[i])
		b.WriteString(",")
	}

	baseDn := b.String()
	baseDnLen := len(baseDn)

	if baseDnLen > 1 {
		return baseDn[:baseDnLen-1]
	} else {
		return "<nil>"
	}
}
