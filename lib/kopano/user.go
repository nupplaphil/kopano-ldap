package kopano

import (
	"bytes"
	"fmt"
	"github.com/nupplaphil/kopano-ldap/lib/utils"
	"gopkg.in/ldap.v2"
	"log"
	"os"
	"text/tabwriter"
)

type UserSettings struct {
	User     string
	Fullname string
	Email    string
	Aliase   []string
	Password string
	Active   bool
}

func ListAll(l *ldap.Conn, baseDn string) {
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=kopano-user))",                                  // The filter to apply
		[]string{"uid", "kopanoAccount", "cn", "mail", "kopanoAliases"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(w, "User\t Active\t Full Name\t E-Mail\t Aliase")
	fmt.Fprintln(w, "----\t ------\t ---------\t ------\t ------")
	for _, entry := range sr.Entries {
		var b bytes.Buffer

		b.WriteString(entry.GetAttributeValue("uid"))
		b.WriteString("\t ")
		b.WriteString(utils.LdapBoolToStr(entry.GetAttributeValue("kopanoAccount")))
		b.WriteString("\t ")
		b.WriteString(entry.GetAttributeValue("cn"))
		b.WriteString("\t ")
		b.WriteString(entry.GetAttributeValue("mail"))
		b.WriteString("\t ")
		b.WriteString(utils.LdapArrayToStr(entry.GetAttributeValues("kopanoAliases"), ","))
		fmt.Fprintln(w, b.String())
	}
	w.Flush()
}

func ListUser(l *ldap.Conn, baseDn, user string) {
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid="+user+")", // The filter to apply
		nil,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		log.Fatalf("No user with uid '" + user + "'")
		os.Exit(1)
	}

	entry := sr.Entries[0]

	w := tabwriter.NewWriter(os.Stdout, 25, 0, 1, ' ', 0)
	fmt.Fprintln(w, fmt.Sprintf("Name:\t %s", entry.GetAttributeValue("uid")))
	fmt.Fprintln(w, fmt.Sprintf("Full name:\t %s", entry.GetAttributeValue("cn")))
	fmt.Fprintln(w, fmt.Sprintf("Email address:\t %s", entry.GetAttributeValue("mail")))
	fmt.Fprintln(w, fmt.Sprintf("Active:\t %s", utils.LdapBoolToStr(entry.GetAttributeValue("kopanoAccount"))))
	fmt.Fprintln(w, fmt.Sprintf("Administrator:\t %s", utils.LdapBoolToStr(entry.GetAttributeValue("kopanoAdmin"))))
	fmt.Fprintln(w, fmt.Sprintf("Features Enabled:\t %s", utils.LdapArrayToStr(entry.GetAttributeValues("kopanoEnabledFeatures"), ";")))
	fmt.Fprintln(w, fmt.Sprintf("Features Disabled:\t %s", utils.LdapArrayToStr(entry.GetAttributeValues("kopanoDisabledFeatures"), ";")))
	w.Flush()
}

func NewUserSettings(user string) *UserSettings {
	return &UserSettings{
		User: user,
	}
}

func Add(l *ldap.Conn, baseDn string, settings *UserSettings) {
	defer l.Close()

	uidNumber, gidNumber := utils.GetNextIDs(l)

	addRequest := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", settings.User, baseDn))

	addRequest.Attribute("objectClass", []string{"posixAccount", "top", "kopano-user", "inetOrgPerson"})
	addRequest.Attribute("homeDirectory", []string{fmt.Sprintf("/home/%s", settings.User)})
	addRequest.Attribute("mail", []string{fmt.Sprintf("%s", settings.Email)})
	actBool := "0"
	if settings.Active {
		actBool = "1"
	}
	addRequest.Attribute("kopanoAccount", []string{actBool})
	addRequest.Attribute("kopanoAdmin", []string{"0"})
	addRequest.Attribute("userPassword", []string{fmt.Sprintf("%s", settings.Password)})
	addRequest.Attribute("kopanoUserServer", []string{"node1"})
	addRequest.Attribute("cn", []string{fmt.Sprintf("%s", settings.Fullname)})
	addRequest.Attribute("sn", []string{fmt.Sprintf("%s", settings.Fullname)})
	addRequest.Attribute("uid", []string{fmt.Sprintf("%s", settings.User)})
	addRequest.Attribute("uidNumber", []string{fmt.Sprintf("%d", uidNumber)})
	addRequest.Attribute("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})
	addRequest.Attribute("kopanoAliases", settings.Aliase)
	addRequest.Attribute("kopanoEnabledFeatures", []string{MOBILE})
	addRequest.Attribute("kopanoDisabledFeatures", []string{IMAP, POP3})

	err := l.Add(addRequest)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func Del(l *ldap.Conn, baseDn, user string) {
	defer l.Close()

	delRequest := ldap.NewDelRequest(fmt.Sprintf("uid=%s,%s", user, baseDn), nil)

	err := l.Del(delRequest)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
