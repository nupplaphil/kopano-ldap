package kopano

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestAddGroup(t *testing.T) {
	tests := map[string]struct {
		baseDn    string
		group     string
		active    bool
		security  bool
		hidden    bool
		nextIdErr error
		addErr    error
		error     error
	}{
		"testNormal": {
			"dc=example,dc=org",
			"testGroup",
			true,
			true,
			false,
			nil,
			nil,
			nil,
		},
		"testNegNextId": {
			"dc=example,dc=org",
			"testGroup",
			true,
			true,
			false,
			fmt.Errorf("a next id error"),
			nil,
			fmt.Errorf("a next id error"),
		},

		"testNegAdd": {
			"dc=example,dc=org",
			"testGroup",
			true,
			true,
			false,
			nil,
			fmt.Errorf("an add error"),
			fmt.Errorf("an add error"),
		},
	}

	for _, test := range tests {
		settings := NewGroupSettings(test.group)
		settings.Active = test.active
		settings.Security = test.security
		settings.Hidden = test.hidden

		addRequest := ldap.AddRequest{
			DN: fmt.Sprintf("cn=%s,ou=Groups,%s", test.group, test.baseDn),
		}
		addRequest.Attribute("objectClass", []string{"posixGroup", "top", "kopano-group"})
		actBool := "0"
		if settings.Active {
			actBool = "1"
		}
		addRequest.Attribute("kopanoAccount", []string{actBool})
		addRequest.Attribute("cn", []string{fmt.Sprintf("%s", test.group)})
		addRequest.Attribute("gidNumber", []string{fmt.Sprintf("%d", 1)})
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

		searchRequest := ldap.SearchRequest{
			BaseDN:       test.baseDn,
			Scope:        ldap.ScopeWholeSubtree,
			DerefAliases: ldap.NeverDerefAliases,
			SizeLimit:    0,
			TimeLimit:    0,
			TypesOnly:    false,
			Filter:       "(&(objectClass=posixGroup))",
			Attributes:   []string{"gidNumber"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry(test.baseDn, map[string][]string{
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

		err := AddGroup(client, test.baseDn, settings)
		assert.Equal(t, err, test.error)
	}
}

func TestDelGroup(t *testing.T) {
	tests := map[string]struct {
		baseDn   string
		group    string
		delError error
	}{
		"testNormal": {
			"dc=example,dc=org",
			"testGroup",
			nil,
		},
		"testNegNextId": {
			"dc=example,dc=org",
			"testGroup",
			fmt.Errorf("a del error"),
		},
	}

	for _, test := range tests {
		delRequest := ldap.DelRequest{
			DN: fmt.Sprintf("cn=%s,ou=Groups,%s", test.group, test.baseDn),
		}

		client := &mocks.Client{}
		client.On("Del", &delRequest).Return(test.delError).Once()
		client.On("Close").Once()

		err := DelGroup(client, test.baseDn, test.group)
		assert.Equal(t, err, test.delError)
	}
}
