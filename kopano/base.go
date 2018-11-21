package kopano

import (
	"bytes"
	"fmt"
	"gopkg.in/ldap.v2"
)

// Connect creates a new LDAP client
func Connect(host string, port int, fqdn, user, password string) (ldap.Client, error) {

	baseDn := GetBaseDN(fqdn)

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.WriteString("cn=")
	b.WriteString(user)
	b.WriteString(",")
	b.WriteString(baseDn)

	err = conn.Bind(b.String(), password)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
