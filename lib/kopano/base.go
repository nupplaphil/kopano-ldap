package kopano

import (
	"bytes"
	"fmt"
	"github.com/nupplaphil/kopano-ldap/lib/utils"
	"gopkg.in/ldap.v2"
	"log"
)

func Connect(host string, port int, fqdn, user, password string) *ldap.Conn {

	baseDn := utils.GetBaseDN(fqdn)

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	b.WriteString("cn=")
	b.WriteString(user)
	b.WriteString(",")
	b.WriteString(baseDn)

	err = l.Bind(b.String(), password)
	if err != nil {
		log.Fatal(err)
	}

	return l
}
