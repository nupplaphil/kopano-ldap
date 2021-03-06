package kopano

import (
	"bytes"
	"gopkg.in/ldap.v2"
	"strconv"
	"strings"
)

// GetNextUserID returns the next valid UID for the given base DN
func GetNextUserID(client ldap.Client, baseDn string) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=posixAccount))",
		[]string{"uidNumber"},
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		return -1, err
	}

	uidNumber := 0
	for _, entry := range sr.Entries {
		uidNumTemp, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))

		if err != nil {
			return -1, err
		}

		if uidNumTemp > uidNumber {
			uidNumber = uidNumTemp
		}
	}
	uidNumber++

	return uidNumber, nil
}

// GetNextGroupID returns the next valid GID for the given base DN
func GetNextGroupID(client ldap.Client, baseDn string) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=posixGroup))",
		[]string{"gidNumber"},
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		return -1, err
	}

	gidNumber := 0
	for _, entry := range sr.Entries {
		gidNumTemp, err := strconv.Atoi(entry.GetAttributeValue("gidNumber"))

		if err != nil {
			return -1, err
		}

		if gidNumTemp > gidNumber {
			gidNumber = gidNumTemp
		}
	}
	gidNumber++

	return gidNumber, nil
}

// LdapBoolToStr returns either "yes" or "no"
// "yes" if anything except "0" or "" is given
// "no" if "0" or "" is given
func LdapBoolToStr(value string) string {
	if len(value) > 0 && value != "0" {
		return "yes"
	}

	return "no"
}

// LdapArrayToStr concatenates an array to an output string with a given separator
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
	}

	return output
}

// GetBaseDN returns the DN string based on a fully qualified domain name
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
	}

	return "<nil>"
}
