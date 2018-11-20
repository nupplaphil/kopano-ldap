package kopano

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"os"
	"sort"
)

const (
	IMAP   = "imap"
	POP3   = "pop3"
	MOBILE = "mobile"
)

func AddUserFeatures(client ldap.Client, baseDn, user string, features []string) {
	defer client.Close()

	checkFeatures(features)
	enabledFeatures, disabledFeatures := GetUserFeatures(client, baseDn, user)

	var modifyAddEnabled []string
	var modifyRemDisabled []string

	for i := range features {
		feature := features[i]
		if findFeatureInFeatures(feature, enabledFeatures) {
			// already enabled
			continue
		} else {
			if findFeatureInFeatures(feature, disabledFeatures) {
				// remove from disabled list
				modifyRemDisabled = append(modifyRemDisabled, feature)
			}
			// add to enabled list
			modifyAddEnabled = append(modifyAddEnabled, feature)
		}
	}

	modify := ldap.NewModifyRequest(fmt.Sprintf("uid=%s,%s", user, baseDn))
	if len(modifyRemDisabled) > 0 {
		modify.Delete("kopanoDisabledFeatures", modifyRemDisabled)
	}
	if len(modifyAddEnabled) > 0 {
		modify.Add("kopanoEnabledFeatures", modifyAddEnabled)
	}

	err := client.Modify(modify)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func RemoveUserFeatures(client ldap.Client, baseDn, user string, features []string) {
	defer client.Close()

	checkFeatures(features)
	enabledFeatures, disabledFeatures := GetUserFeatures(client, baseDn, user)

	var modifyAddDisabled []string
	var modifyRemEnabled []string

	for i := range features {
		feature := features[i]
		if findFeatureInFeatures(feature, disabledFeatures) {
			// already enabled
			continue
		} else {
			if findFeatureInFeatures(feature, enabledFeatures) {
				// remove from disabled list
				modifyRemEnabled = append(modifyRemEnabled, feature)
			}
			// add to enabled list
			modifyAddDisabled = append(modifyAddDisabled, feature)
		}
	}

	modify := ldap.NewModifyRequest(fmt.Sprintf("uid=%s,%s", user, baseDn))
	if len(modifyRemEnabled) > 0 {
		modify.Delete("kopanoEnabledFeatures", modifyRemEnabled)
	}
	if len(modifyAddDisabled) > 0 {
		modify.Add("kopanoDisabledFeatures", modifyAddDisabled)
	}

	err := client.Modify(modify)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func checkFeatures(features []string) {
	for i := range features {
		if !isValid(features[i]) {
			log.Fatal("Adding Feature '" + features[i] + "' is not valid")
			os.Exit(1)
		}
	}
}

func isValid(feature string) bool {
	return feature == IMAP ||
		feature == POP3 ||
		feature == MOBILE
}

func GetUserFeatures(client ldap.Client, baseDn, user string) ([]string, []string) {

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid="+user+")", // The filter to apply
		[]string{"kopanoEnabledFeatures", "kopanoDisabledFeatures"}, // A list attributes to retrieve
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		log.Fatalf("No user with uid '" + user + "'")
		os.Exit(1)
	}

	entry := sr.Entries[0]

	return entry.GetAttributeValues("kopanoEnabledFeatures"), entry.GetAttributeValues("kopanoDisabledFeatures")
}

func findFeatureInFeatures(feature string, features []string) bool {
	sort.Strings(features)

	i := sort.SearchStrings(features, feature)
	if i < len(features) && features[i] == feature {
		return true
	} else {
		return false
	}
}
