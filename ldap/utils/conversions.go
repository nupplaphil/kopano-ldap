package utils

import "bytes"

func LdapBool(value string) string {
	if value == "1" {
		return "yes"
	} else {
		return "no"
	}
}

func LdapArrayToStr(values []string, seperator string) string {
	var b bytes.Buffer
	for i := range values {
		b.WriteString(values[i])
		b.WriteString(seperator + " ")
	}
	output := b.String()
	outputLen := len(output)
	if outputLen > 1 {
		return output[:outputLen-2]
	} else {
		return output
	}
}
