package utils

import (
	"github.com/magiconair/properties/assert"
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
