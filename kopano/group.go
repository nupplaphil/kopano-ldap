package kopano

import (
	"bytes"
	"fmt"
	"gopkg.in/ldap.v2"
	"io"
	"text/tabwriter"
)

// GroupSettings for an new user or to modify an user
type GroupSettings struct {
	Name     string
	Active   bool
	Security bool
	Hidden   bool
}

// NewGroupSettings creates a new GroupSettings attribute for modifying or creating a group
func NewGroupSettings(name string) *GroupSettings {
	return &GroupSettings{
		Name:     name,
		Active:   true,
		Security: true,
		Hidden:   false,
	}
}

// ListAllGroups lists all groups for the given base DN to an output writer
func ListAllGroups(client ldap.Client, baseDn string, writer io.Writer) error {
	defer client.Close()

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=kopano-group))",
		[]string{"cn", "kopanoAccount", "kopanoHidden", "kopanpSecurityGroup"},
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(writer, 5, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(w, "Group\t Active\t Hidden\t Security Group")
	fmt.Fprintln(w, "-----\t ------\t ------\t --------------")
	for _, entry := range sr.Entries {
		var b bytes.Buffer

		b.WriteString(entry.GetAttributeValue("cn"))
		b.WriteString("\t ")
		b.WriteString(LdapBoolToStr(entry.GetAttributeValue("kopanoAccount")))
		b.WriteString("\t ")
		b.WriteString(LdapBoolToStr(entry.GetAttributeValue("kopanoHidden")))
		b.WriteString("\t ")
		b.WriteString(LdapBoolToStr(entry.GetAttributeValue("kopanoSecurityGroup")))
		fmt.Fprintln(w, b.String())
	}
	w.Flush()

	return nil
}

// ListGroup lists details of the given group for a given base DN to an output writer
func ListGroup(client ldap.Client, baseDn, name string, writer io.Writer) error {
	defer client.Close()

	groupUser, err := getUserOfGroup(client, baseDn, name)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(writer, 5, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(w, "Users (%d):\n", len(groupUser))
	fmt.Fprintln(w, "\t Username\t Fullname")
	fmt.Fprintln(w, "\t --------\t --------")
	for _, user := range groupUser {
		var b bytes.Buffer

		b.WriteString("\t ")
		b.WriteString(user.User)
		b.WriteString("\t ")
		b.WriteString(user.Fullname)
		fmt.Fprintln(w, b.String())
	}
	w.Flush()

	return nil
}

// AddGroup Adds a new group
func AddGroup(client ldap.Client, baseDn string, settings *GroupSettings) error {
	defer client.Close()

	gidNumber, err := GetNextGroupID(client, baseDn)
	if err != nil {
		return err
	}

	addRequest := ldap.NewAddRequest(fmt.Sprintf("cn=%s,ou=Groups,%s", settings.Name, baseDn))

	addRequest.Attribute("objectClass", []string{"posixGroup", "top", "kopano-group"})
	actBool := "0"
	if settings.Active {
		actBool = "1"
	}
	addRequest.Attribute("kopanoAccount", []string{actBool})
	addRequest.Attribute("cn", []string{fmt.Sprintf("%s", settings.Name)})
	addRequest.Attribute("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})
	secBool := "0"
	if settings.Security {
		secBool = "1"
	}
	addRequest.Attribute("kopanoSecurityGroup", []string{secBool})
	hiddenBool := "0"
	if settings.Hidden {
		hiddenBool = "1"
	}
	addRequest.Attribute("kopanoHidden", []string{hiddenBool})

	if err = client.Add(addRequest); err != nil {
		return err
	}

	return nil
}

func getUserOfGroup(client ldap.Client, baseDn, name string) ([]UserSettings, error) {

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s,ou=Groups,%s)", name, baseDn),
		[]string{"memberUid"},
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	settings := []UserSettings{}
	for _, entry := range sr.Entries {
		setting := UserSettings{
			User: entry.GetAttributeValue("memberUid"),
		}

		settings = append(settings, setting)
	}

	return settings, nil
}

// AddUserToGroup adds a list of users to a given group
func AddUserToGroup(client ldap.Client, baseDn, group string, userList []string) error {
	defer client.Close()

	modifyRequest := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=Groups,%s", group, baseDn))

	addUser := []string{}
	for _, user := range userList {
		addUser = append(addUser, fmt.Sprintf("cn=%s,%s", user, baseDn))
	}

	modifyRequest.Add("memberUid", addUser)

	if err := client.Modify(modifyRequest); err != nil {
		return err
	}

	return nil
}

// DelUserFromGroup deletes a list of users to a given group
func DelUserFromGroup(client ldap.Client, baseDn, group string, userList []string) error {
	defer client.Close()

	modifyRequest := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=Groups,%s", group, baseDn))

	addUser := []string{}
	for _, user := range userList {
		addUser = append(addUser, fmt.Sprintf("cn=%s,%s", user, baseDn))
	}

	modifyRequest.Delete("memberUid", addUser)

	if err := client.Modify(modifyRequest); err != nil {
		return err
	}

	return nil
}

// DelGroup deletes a group
func DelGroup(client ldap.Client, baseDn, user string) error {
	defer client.Close()

	delRequest := ldap.NewDelRequest(fmt.Sprintf("cn=%s,ou=Groups,%s", user, baseDn), nil)

	if err := client.Del(delRequest); err != nil {
		return err
	}

	return nil
}
