package kopano

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestModifyFeatures(t *testing.T) {
	tests := map[string]struct {
		baseDn       string
		user         string
		features     []string
		enabledFeat  []string
		disabledFeat []string
		removeFeat   []string
		enableFeat   []string
		errUserFeat  error
		searchFailed bool
		modError     error
		error        error
	}{
		"testNormal": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"imap", "mobile"},
		},
		"testAdding": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{"imap"},
			removeFeat:   []string{"imap"},
			enableFeat:   []string{"imap", "mobile"},
		},
		"testAdding2": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{"mobile"},
			disabledFeat: []string{"imap"},
			removeFeat:   []string{"imap"},
			enableFeat:   []string{"imap"},
		},
		"testAdding3": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{"mobile", "imap"},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{},
		},
		"testAdding4": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile", "pop3"},
			enabledFeat:  []string{"mobile", "imap"},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"pop3"},
		},
		"testInvalidFeat": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"test", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"test", "mobile"},
			error:        fmt.Errorf("adding feature '%q' is not valid", "test"),
		},
		"testNegUserFeat": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"imap", "mobile"},
			errUserFeat:  fmt.Errorf("an expected search error"),
			error:        fmt.Errorf("an expected search error"),
		},
		"testEmptyResult": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"imap", "mobile"},
			searchFailed: true,
			error:        fmt.Errorf("no user with uid %q found", "johndoe"),
		},
		"testNegModify": {
			baseDn:       "dc=example,dc=org",
			user:         "johndoe",
			features:     []string{"imap", "mobile"},
			enabledFeat:  []string{},
			disabledFeat: []string{},
			removeFeat:   []string{},
			enableFeat:   []string{"imap", "mobile"},
			modError:     fmt.Errorf("an expected modify error"),
			error:        fmt.Errorf("an expected modify error"),
		},
	}

	for _, test := range tests {

		searchRequest := ldap.SearchRequest{
			BaseDN:       test.baseDn,
			Scope:        ldap.ScopeWholeSubtree,
			DerefAliases: ldap.NeverDerefAliases,
			SizeLimit:    0,
			TimeLimit:    0,
			TypesOnly:    false,
			Filter:       fmt.Sprintf("(uid=%s)", test.user),
			Attributes:   []string{"kopanoEnabledFeatures", "kopanoDisabledFeatures"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry(test.baseDn, map[string][]string{
				"kopanoEnabledFeatures":  test.enabledFeat,
				"kopanoDisabledFeatures": test.disabledFeat,
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		if test.searchFailed {
			searchResult.Entries = nil
		}

		modifyRequest := ldap.ModifyRequest{
			DN: fmt.Sprintf("uid=%s,%s", test.user, test.baseDn),
		}
		if len(test.removeFeat) > 0 {
			modifyRequest.Delete("kopanoDisabledFeatures", test.removeFeat)
		}
		if len(test.enableFeat) > 0 {
			modifyRequest.Add("kopanoEnabledFeatures", test.enableFeat)
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, test.errUserFeat).Once()
		client.On("Modify", &modifyRequest).Return(test.modError).Once()
		client.On("Close").Once()

		err := AddUserFeatures(client, test.baseDn, test.user, test.features)
		assert.Equal(t, err, test.error)

		searchResult = ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry(test.baseDn, map[string][]string{
				"kopanoEnabledFeatures":  test.disabledFeat,
				"kopanoDisabledFeatures": test.enabledFeat,
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		if test.searchFailed {
			searchResult.Entries = nil
		}

		modifyRequest = ldap.ModifyRequest{
			DN: fmt.Sprintf("uid=%s,%s", test.user, test.baseDn),
		}
		if len(test.disabledFeat) > 0 {
			modifyRequest.Delete("kopanoEnabledFeatures", test.disabledFeat)
		}
		if len(test.enableFeat) > 0 {
			modifyRequest.Add("kopanoDisabledFeatures", test.enableFeat)
		}

		client = &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, test.errUserFeat).Once()
		client.On("Modify", &modifyRequest).Return(test.modError).Once()
		client.On("Close").Once()

		err = RemoveUserFeatures(client, test.baseDn, test.user, test.features)
		assert.Equal(t, err, test.error)
	}
}
