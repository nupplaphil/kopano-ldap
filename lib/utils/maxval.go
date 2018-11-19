package utils

import (
	"gopkg.in/ldap.v2"
	"log"
	"strconv"
)

func GetNextIDs(conn *ldap.Conn) (int, int) {
	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=org", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=posixAccount))",    // The filter to apply
		[]string{"uidNumber", "gidNumber"}, // A list attributes to retrieve
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	uidNumber := 0
	for _, entry := range sr.Entries {
		uidNumTemp, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))

		if err != nil {
			log.Fatal(err)
		}

		if uidNumTemp > uidNumber {
			uidNumber = uidNumTemp
		}
	}
	uidNumber++

	gidNumber := 0
	for _, entry := range sr.Entries {
		gidNumTemp, err := strconv.Atoi(entry.GetAttributeValue("gidNumber"))

		if err != nil {
			log.Fatal(err)
		}

		if gidNumTemp > gidNumber {
			gidNumber = gidNumTemp
		}
	}
	gidNumber++

	return uidNumber, gidNumber
}
