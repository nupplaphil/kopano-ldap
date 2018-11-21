package kopano

import (
	"github.com/magiconair/properties/assert"
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestLdapBoolToStr(t *testing.T) {
	assertTrue := "yes"
	assertFalse := "no"

	test := LdapBoolToStr("1")
	assert.Equal(t, test, assertTrue)

	test = LdapBoolToStr("0")
	assert.Equal(t, test, assertFalse)

	test = LdapBoolToStr("")
	assert.Equal(t, test, assertFalse)

	test = LdapBoolToStr("10")
	assert.Equal(t, test, assertTrue)

	test = LdapBoolToStr("test")
	assert.Equal(t, test, assertTrue)
}

func TestLdapArrayToStr(t *testing.T) {
	assertStr1 := "test1, test2, test3"
	str1 := LdapArrayToStr([]string{"test1", "test2", "test3"}, ",")
	assert.Equal(t, str1, assertStr1)

	assertStr2 := "test1. test2. test3"
	str2 := LdapArrayToStr([]string{"test1", "test2", "test3"}, ".")
	assert.Equal(t, str2, assertStr2)

	assertStr3 := "test1 test2 test3"
	str3 := LdapArrayToStr([]string{"test1", "test2", "test3"}, "")
	assert.Equal(t, str3, assertStr3)
}

func TestGetBaseDN(t *testing.T) {
	assertBaseDn1 := "dc=example,dc=org"
	baseDn1 := GetBaseDN("example.org")
	assert.Equal(t, baseDn1, assertBaseDn1)

	assertBaseDn2 := "<nil>"
	baseDn2 := GetBaseDN("")
	assert.Equal(t, baseDn2, assertBaseDn2)

	assertBaseDn3 := "dc=1,dc=2"
	baseDn3 := GetBaseDN("1.2.")
	assert.Equal(t, baseDn3, assertBaseDn3)
}

func TestGetNextIDs(t *testing.T) {
	tests := map[string]struct {
		baseDn string
		uid    string
		gid    string
	}{
		"test1": {
			"example.org",
			"1",
			"1",
		},
		"test2": {
			"example.org",
			"0",
			"0",
		},
		"test3": {
			"example.",
			"-1",
			"-1",
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
			Filter:       "(&(objectClass=posixAccount))",
			Attributes:   []string{"uidNumber", "gidNumber"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
				"uidNumber": {test.uid},
				"gidNumber": {test.gid},
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, nil).Once()
		client.On("Close").Once()

		GetNextIDs(client, test.baseDn)
	}
}
