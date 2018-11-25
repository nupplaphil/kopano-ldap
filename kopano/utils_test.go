package kopano

import (
	"errors"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/nupplaphil/kopano-ldap/mocks"
	"gopkg.in/ldap.v2"
	"strconv"
	"testing"
)

func TestLdapBoolToStr(t *testing.T) {
	assertTrue := "yes"
	assertFalse := "no"

	tests := map[string]struct {
		value     string
		assertVal string
	}{
		"test1": {"1", assertTrue},
		"test2": {"0", assertFalse},
		"test3": {"", assertFalse},
		"test4": {"10", assertTrue},
		"test5": {"test", assertTrue},
	}

	for _, test := range tests {
		testEx := LdapBoolToStr(test.value)
		assert.Equal(t, testEx, test.assertVal)
	}
}

func TestLdapArrayToStr(t *testing.T) {
	tests := map[string]struct {
		values    []string
		separator string
		assertVal string
	}{
		"test1": {
			[]string{"test1", "test2", "test3"},
			",",
			"test1, test2, test3",
		},
		"test2": {
			[]string{"test1", "test2", "test3"},
			".",
			"test1. test2. test3",
		},
		"test3": {
			[]string{"test1", "test2", "test3"},
			"",
			"test1 test2 test3",
		},
		"test4": {
			[]string{""},
			"",
			" ",
		},
	}

	for _, test := range tests {
		testStr := LdapArrayToStr(test.values, test.separator)
		assert.Equal(t, testStr, test.assertVal)
	}
}

func TestGetBaseDN(t *testing.T) {
	tests := map[string]struct {
		domain string
		baseDn string
	}{
		"test1": {"example.org", "dc=example,dc=org"},
		"test2": {"", "<nil>"},
		"test3": {"1.2.", "dc=1,dc=2"},
	}

	for _, test := range tests {
		testDn := GetBaseDN(test.domain)
		assert.Equal(t, testDn, test.baseDn)
	}
}

func TestGetNextUserAndGroupID(t *testing.T) {
	tests := map[string]struct {
		baseDn    string
		id        string
		searchErr error
		error     error
	}{
		"testNormal": {
			"example.org",
			"1",
			nil,
			nil,
		},
		"testTest0": {
			"example.org",
			"0",
			nil,
			nil,
		},
		"testNegSearch": {
			"example.",
			"-1",
			fmt.Errorf("an expected error"),
			fmt.Errorf("an expected error"),
		},
		"testNegCastUid": {
			"example.",
			"bla",
			nil,
			&strconv.NumError{
				Func: "Atoi",
				Num:  "bla",
				Err:  errors.New("invalid syntax"),
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
			Filter:       "(&(objectClass=posixAccount))",
			Attributes:   []string{"uidNumber"},
			Controls:     nil,
		}

		searchResult := ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
				"uidNumber": {test.id},
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		client := &mocks.Client{}
		client.On("Search", &searchRequest).Return(&searchResult, test.searchErr).Once()
		client.On("Close").Once()

		uid, err := GetNextUserID(client, test.baseDn)

		// TODO
		// Dirty hack because "atoi.go" fallback "ParseInt()" doesn't properly
		// set the function "Atoi" for fallback "ParseInt" function
		if nerr, ok := err.(*strconv.NumError); ok {
			nerr.Func = "Atoi"
			err = nerr
		}
		assert.Equal(t, err, test.error)

		if err == nil {
			assertUid, _ := strconv.Atoi(test.id)
			assert.Equal(t, uid, assertUid+1)
		}

		searchRequest = ldap.SearchRequest{
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

		searchResult = ldap.SearchResult{
			Entries: []*ldap.Entry{ldap.NewEntry("test", map[string][]string{
				"gidNumber": {test.id},
			}),
			},
			Referrals: nil,
			Controls:  nil,
		}

		client.On("Search", &searchRequest).Return(&searchResult, test.searchErr).Once()
		client.On("Close").Once()

		uid, err = GetNextGroupID(client, test.baseDn)

		// TODO
		// Dirty hack because "atoi.go" fallback "ParseInt()" doesn't properly
		// set the function "Atoi" for fallback "ParseInt" function
		if nerr, ok := err.(*strconv.NumError); ok {
			nerr.Func = "Atoi"
			err = nerr
		}
		assert.Equal(t, err, test.error)

		if err == nil {
			assertUid, _ := strconv.Atoi(test.id)
			assert.Equal(t, uid, assertUid+1)
		}
	}
}
