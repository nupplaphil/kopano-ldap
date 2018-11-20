package kopano

import (
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestListUser(t *testing.T) {
	tests := map[string]struct {
		baseDn  string
		entries map[string][]string
	}{
		"test1": {
			"example.org",
			map[string][]string{
				"uid":           {"philipp7"},
				"kopanoAccount": {"1"},
				"cn":            {"Philipp Holzer"},
				"mail":          {"philipp@dieholzers.at"},
				"kopanoAliases": {"admin@philipp.info", "test@example.org"},
			},
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
			Filter:       "(&(objectClass=kopano-user))",
			Attributes:   []string{"uid", "kopanoAccount", "cn", "mail", "kopanoAliases"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries:   []*ldap.Entry{ldap.NewEntry("test", test.entries)},
			Referrals: nil,
			Controls:  nil,
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, nil).Once()
		client.On("Close").Once()

		ListAll(client, "example.org")
	}
}
