package utils

import (
	"bytes"
	"strings"
)

func GetBaseDN(fqdn string) string {
	parts := strings.Split(fqdn, ".")

	var b bytes.Buffer

	for i := range parts {
		b.WriteString("dc=")
		b.WriteString(parts[i])
		b.WriteString(",")
	}

	baseDn := b.String()
	baseDnLen := len(baseDn)

	if baseDnLen > 0 {
		return baseDn[:baseDnLen-1]
	} else {
		return "<nil>"
	}
}
