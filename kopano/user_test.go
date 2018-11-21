package kopano

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestListAll(t *testing.T) {
	tests := map[string]struct {
		baseDn  string
		entries map[string][]string
		error   error
	}{
		"testNormal": {
			"example.org",
			map[string][]string{
				"uid":           {"philipp7"},
				"kopanoAccount": {"1"},
				"cn":            {"Philipp Holzer"},
				"mail":          {"philipp@dieholzers.at"},
				"kopanoAliases": {"admin@philipp.info", "test@example.org"},
			},
			nil,
		},
		"testStrageNormal": {
			"example.org",
			map[string][]string{
				"uid":           {""},
				"kopanoAccount": {"0"},
				"kopanoAdmin  ": {"1"},
				"cn":            {""},
				"mail":          {"@dieholzers.at"},
				"kopanoAliases": {"admin@philipp.info"},
			},
			nil,
		},
		"testError": {
			"example.org",
			map[string][]string{
				"uid":           {""},
				"kopanoAccount": {"0"},
				"kopanoAdmin  ": {"1"},
				"cn":            {""},
				"mail":          {"@dieholzers.at"},
				"kopanoAliases": {"admin@philipp.info"},
			},
			fmt.Errorf("an expected error"),
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
		client.On("Search", &searchRequest).Return(&searchResult, test.error).Once()
		client.On("Close").Once()

		err := ListAll(client, "example.org")
		assert.Equal(t, err, test.error)
	}
}

func TestListUser(t *testing.T) {
	tests := map[string]struct {
		baseDn    string
		user      string
		result    ldap.SearchResult
		searchErr error
		error     error
	}{
		"testNormal": {
			"example.org",
			"philipp7",
			ldap.SearchResult{
				Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
					"uid":                    {"philipp7"},
					"cn":                     {"Philipp Holzer"},
					"mail":                   {"philipp@dieholzers.at"},
					"kopanoAccount":          {"1"},
					"kopanoAdmin":            {"1"},
					"kopanoEnabledFeatures":  {"imap; mobile"},
					"kopanoDisabledFeatures": {"pop3"},
				},
				)},
				Referrals: nil,
				Controls:  nil,
			},
			nil,
			nil,
		},
		"testSearchError": {
			"example.org",
			"philipp7",
			ldap.SearchResult{
				Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
					"uid":                    {"philipp7"},
					"cn":                     {"Philipp Holzer"},
					"mail":                   {"philipp@dieholzers.at"},
					"kopanoAccount":          {"1"},
					"kopanoAdmin":            {"1"},
					"kopanoEnabledFeatures":  {"imap; mobile"},
					"kopanoDisabledFeatures": {"pop3"},
				})},
				Referrals: nil,
				Controls:  nil,
			},
			fmt.Errorf("an expected error"),
			fmt.Errorf("an expected error"),
		},
		"testNoEntries": {
			"example.org",
			"philipp6",
			ldap.SearchResult{
				Entries: nil,
			},
			nil,
			fmt.Errorf("no user with uid 'philipp6'"),
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
			Filter:       "(uid=" + test.user + ")",
			Attributes:   nil,
			Controls:     nil,
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&test.result, test.searchErr).Once()
		client.On("Close").Once()

		err := ListUser(client, "example.org", test.user)
		assert.Equal(t, err, test.error)
	}
}
