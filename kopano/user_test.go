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
				"uid":           {"johnedoe7"},
				"kopanoAccount": {"1"},
				"cn":            {"John Doe"},
				"mail":          {"john@doe.com"},
				"kopanoAliases": {"john.d@doe.com", "test@example.org"},
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
				"mail":          {"@doe.com"},
				"kopanoAliases": {"john@doe.com"},
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
				"mail":          {"@doe.com"},
				"kopanoAliases": {"john@doe.com"},
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

		err := ListAll(client, test.baseDn)
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
			"dc=example,dc=org",
			"johndoe7",
			ldap.SearchResult{
				Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
					"uid":                    {"johndoe7"},
					"cn":                     {"John Doe"},
					"mail":                   {"john@doe.com"},
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
			"dc=example,dc=org",
			"johndoe7",
			ldap.SearchResult{
				Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
					"uid":                    {"johndoe7"},
					"cn":                     {"John Doe"},
					"mail":                   {"john@doe.com"},
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
			"dc=example,dc=org",
			"johndoe6",
			ldap.SearchResult{
				Entries: nil,
			},
			nil,
			fmt.Errorf("no user with uid \"johndoe6\""),
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

		err := ListUser(client, test.baseDn, test.user)
		assert.Equal(t, err, test.error)
	}
}

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		baseDn    string
		user      string
		email     string
		active    bool
		password  string
		fullname  string
		aliases   []string
		nextIdErr error
		addErr    error
		error     error
	}{
		"testNormal": {
			"dc=example,dc=org",
			"johndoe",
			"johne@doe.com",
			true,
			"123456",
			"John Doe",
			[]string{"john2@doe.com"},
			nil,
			nil,
			nil,
		},
		"testNegNextId": {
			"dc=example,dc=org",
			"johndoe",
			"johne@doe.com",
			true,
			"123456",
			"John Doe",
			[]string{"john2@doe.com"},
			fmt.Errorf("a next id error"),
			nil,
			fmt.Errorf("a next id error"),
		},

		"testNegAdd": {
			"dc=example,dc=org",
			"johndoe",
			"johne@doe.com",
			true,
			"123456",
			"John Doe",
			[]string{"john2@doe.com"},
			nil,
			fmt.Errorf("an add error"),
			fmt.Errorf("an add error"),
		},
	}

	for _, test := range tests {
		settings := NewUserSettings(test.user)
		settings.Active = test.active
		settings.Aliases = test.aliases
		settings.Email = test.email
		settings.Fullname = test.fullname
		settings.Password = test.password
		settings.User = test.user

		addRequest := ldap.AddRequest{
			DN: fmt.Sprintf("uid=%s,%s", test.user, test.baseDn),
		}
		addRequest.Attribute("objectClass", []string{"posixAccount", "top", "kopano-user", "inetOrgPerson"})
		addRequest.Attribute("homeDirectory", []string{fmt.Sprintf("/home/%s", test.user)})
		addRequest.Attribute("mail", []string{fmt.Sprintf("%s", test.email)})
		actBool := "0"
		if test.active {
			actBool = "1"
		}
		addRequest.Attribute("kopanoAccount", []string{actBool})
		addRequest.Attribute("kopanoAdmin", []string{"0"})
		addRequest.Attribute("userPassword", []string{fmt.Sprintf("%s", test.password)})
		addRequest.Attribute("kopanoUserServer", []string{"node1"})
		addRequest.Attribute("cn", []string{fmt.Sprintf("%s", test.fullname)})
		addRequest.Attribute("sn", []string{fmt.Sprintf("%s", test.fullname)})
		addRequest.Attribute("uid", []string{fmt.Sprintf("%s", test.user)})
		addRequest.Attribute("uidNumber", []string{fmt.Sprintf("%d", 1)})
		addRequest.Attribute("gidNumber", []string{fmt.Sprintf("%d", 1)})
		addRequest.Attribute("kopanoAliases", test.aliases)
		addRequest.Attribute("kopanoEnabledFeatures", []string{MOBILE})
		addRequest.Attribute("kopanoDisabledFeatures", []string{IMAP, POP3})

		searchRequest := ldap.SearchRequest{
			BaseDN:       test.baseDn,
			Scope:        ldap.ScopeWholeSubtree,
			DerefAliases: ldap.NeverDerefAliases,
			SizeLimit:    0,
			TimeLimit:    0,
			TypesOnly:    false,
			Filter:       "(&(objectClass=posixAccount))",
			Attributes:   []string{"uidNumber", "gidNumber"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry(test.baseDn, map[string][]string{
				"uidNumber": {"0"},
				"gidNumber": {"0"},
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, test.nextIdErr).Once()
		client.On("Add", &addRequest).Return(test.addErr).Once()
		client.On("Close").Once()

		err := Add(client, test.baseDn, settings)
		assert.Equal(t, err, test.error)
	}
}

func TestDel(t *testing.T) {
	tests := map[string]struct {
		baseDn   string
		user     string
		delError error
	}{
		"testNormal": {
			"dc=example,dc=org",
			"johndoe",
			nil,
		},
		"testNegNextId": {
			"dc=example,dc=org",
			"johndoe",
			fmt.Errorf("a del error"),
		},
	}

	for _, test := range tests {
		delRequest := ldap.DelRequest{
			DN: fmt.Sprintf("uid=%s,%s", test.user, test.baseDn),
		}

		client := &mocks.Client{}
		client.On("Del", &delRequest).Return(test.delError).Once()
		client.On("Close").Once()

		err := Del(client, test.baseDn, test.user)
		assert.Equal(t, err, test.delError)
	}
}
