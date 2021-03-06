package kopano

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"sort"
)

const (
	// IMAP feature of Kopano
	IMAP = "imap"
	// POP3 feature of Kopano
	POP3 = "pop3"
	// MOBILE feature of Kopano
	MOBILE = "mobile"
)

// AddUserFeatures adds new features to an specific user.
// It will automatically remove these features from the disabled features
func AddUserFeatures(client ldap.Client, baseDn, user string, features []string) error {
	defer client.Close()

	err := checkFeatures(features)
	if err != nil {
		return err
	}

	enabledFeatures, disabledFeatures, err := GetUserFeatures(client, baseDn, user)
	if err != nil {
		return err
	}

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

	err = client.Modify(modify)
	if err != nil {
		return err
	}

	return nil
}

// RemoveUserFeatures removes features from an specific user.
// It will automatically remove these features from the enabled features
func RemoveUserFeatures(client ldap.Client, baseDn, user string, features []string) error {
	defer client.Close()

	err := checkFeatures(features)
	if err != nil {
		return err
	}

	enabledFeatures, disabledFeatures, err := GetUserFeatures(client, baseDn, user)
	if err != nil {
		return err
	}

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

	err = client.Modify(modify)
	if err != nil {
		return err
	}

	return nil
}

// checkFeatures checks a list of features if they are valid
func checkFeatures(features []string) error {
	for i := range features {
		if !isValid(features[i]) {
			return fmt.Errorf("adding feature '%q' is not valid", features[i])
		}
	}

	return nil
}

// isValid checks, if a given feature is valid
func isValid(feature string) bool {
	return feature == IMAP ||
		feature == POP3 ||
		feature == MOBILE
}

// GetUserFeatures returns all enabled and disabled features of an user
func GetUserFeatures(client ldap.Client, baseDn, user string) ([]string, []string, error) {

	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid="+user+")", // The filter to apply
		[]string{"kopanoEnabledFeatures", "kopanoDisabledFeatures"}, // A list attributes to retrieve
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		return nil, nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, nil, fmt.Errorf("no user with uid %q found", user)
	}

	entry := sr.Entries[0]

	return entry.GetAttributeValues("kopanoEnabledFeatures"), entry.GetAttributeValues("kopanoDisabledFeatures"), nil
}

// findFeatureInFeatures checks if the given feature is in a list of features
func findFeatureInFeatures(feature string, features []string) bool {
	sort.Strings(features)

	i := sort.SearchStrings(features, feature)
	if i < len(features) && features[i] == feature {
		return true
	}

	return false
}
