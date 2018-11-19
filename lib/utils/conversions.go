package utils

import "bytes"

// Returns either "yes" or "no"
// "yes" if anything except "0" or "" is given
// "no" if "0" or "" is given
func LdapBoolToStr(value string) string {
	if len(value) > 0 && value != "0" {
		return "yes"
	} else {
		return "no"
	}
}

// concatenate an array to an output string with a given separator
func LdapArrayToStr(values []string, separator string) string {
	var b bytes.Buffer
	for i := range values {
		b.WriteString(values[i])
		b.WriteString(separator + " ")
	}
	output := b.String()
	outputLen := len(output)
	if outputLen > 1 {
		return output[:outputLen-1-len(separator)]
	} else {
		return output
	}
}
