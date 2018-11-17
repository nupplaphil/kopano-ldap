package kopano

import (
	"fmt"
	"github.com/nupplaphil/kopano-ldap/ldap/utils"
	"gopkg.in/ldap.v2"
	"log"
)

type UserSettings struct {
	UID      string
	Name     string
	SName    string
	Mail     string
	Aliase   []string
	Password string
}

func List(l *ldap.Conn) {
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=org", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))", // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
}

func NewUserSettings(uid string) *UserSettings {
	return &UserSettings{
		UID: uid,
	}
}

func Add(l *ldap.Conn, settings *UserSettings) {
	defer l.Close()

	uidNumber, gidNumber := utils.GetNextIDs(l)

	addRequest := ldap.NewAddRequest(fmt.Sprintf("uid=%s,dc=example,dc=org", settings.UID))

	addRequest.Attribute("objectClass", []string{"posixAccount", "top", "kopano-user", "inetOrgPerson"})
	addRequest.Attribute("homeDirectory", []string{fmt.Sprintf("/home/%s", settings.UID)})
	addRequest.Attribute("mail", []string{fmt.Sprintf("%s", settings.Mail)})
	addRequest.Attribute("kopanoAccount", []string{"1"})
	addRequest.Attribute("kopanoAdmin", []string{"0"})
	addRequest.Attribute("userPassword", []string{fmt.Sprintf("%s", settings.Password)})
	addRequest.Attribute("kopanoUserServer", []string{"node1"})
	addRequest.Attribute("cn", []string{fmt.Sprintf("%s", settings.Name)})
	addRequest.Attribute("sn", []string{fmt.Sprintf("%s", settings.SName)})
	addRequest.Attribute("uid", []string{fmt.Sprintf("%s", settings.UID)})
	addRequest.Attribute("uidNumber", []string{fmt.Sprintf("%d", uidNumber)})
	addRequest.Attribute("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})

	err := l.Add(addRequest)
	if err != nil {
		log.Fatal(err)
	}
}
