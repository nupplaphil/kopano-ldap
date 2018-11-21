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

func TestGetNextIDs(t *testing.T) {
	tests := map[string]struct {
		baseDn    string
		uid       string
		gid       string
		searchErr error
		error     error
	}{
		"testNormal": {
			"example.org",
			"1",
			"2",
			nil,
			nil,
		},
		"testTest0": {
			"example.org",
			"0",
			"1",
			nil,
			nil,
		},
		"testNegSearch": {
			"example.",
			"-1",
			"-2",
			fmt.Errorf("an expected error"),
			fmt.Errorf("an expected error"),
		},
		"testNegCastUid": {
			"example.",
			"bla",
			"-3",
			nil,
			&strconv.NumError{
				Func: "Atoi",
				Num:  "bla",
				Err:  errors.New("invalid syntax"),
			},
		},
		"testNegCastGid": {
			"example.",
			"-4",
			"blub",
			nil,
			&strconv.NumError{
				Func: "Atoi",
				Num:  "blub",
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
		client.On("Search", &searchRequest).Return(&searchResult, test.searchErr).Once()
		client.On("Close").Once()

		uid, gid, err := GetNextIDs(client, test.baseDn)
		assert.Equal(t, err, test.error)

		if err == nil {
			assertUid, _ := strconv.Atoi(test.uid)
			assert.Equal(t, uid, assertUid+1)
			assertGid, _ := strconv.Atoi(test.gid)
			assert.Equal(t, gid, assertGid+1)
		}
	}
}
