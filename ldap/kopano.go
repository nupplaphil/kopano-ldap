package ldap

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
)

func Connect() *ldap.Conn {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.99.100", 389))
	if err != nil {
		log.Fatal(err)
	}

	err = l.Bind("cn=admin,dc=example,dc=org", "admin")
	if err != nil {
		log.Fatal(err)
	}

	return l
}
