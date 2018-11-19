package utils

import "testing"

func TestLdapBoolToStr(t *testing.T) {
	assertTrue := "yes"
	assertFalse := "no"

	test := LdapBoolToStr("1")
	if test != assertTrue {
		t.Errorf("Expected %q got %q", assertTrue, test)
	}

	test = LdapBoolToStr("0")
	if test != assertFalse {
		t.Errorf("Expected %q got %q", assertFalse, test)
	}

	test = LdapBoolToStr("")
	if test != assertFalse {
		t.Errorf("Expected %q got %q", assertFalse, test)
	}

	test = LdapBoolToStr("10")
	if test != assertTrue {
		t.Errorf("Expected %q got %q", assertTrue, test)
	}

	test = LdapBoolToStr("test")
	if test != assertTrue {
		t.Errorf("Expected %q got %q", assertTrue, test)
	}
}

func TestLdapArrayToStr(t *testing.T) {
	assertStr1 := "test1, test2, test3"
	str1 := LdapArrayToStr([]string{"test1", "test2", "test3"}, ",")
	if str1 != assertStr1 {
		t.Errorf("Expected %q got %q", assertStr1, str1)
	}

	assertStr2 := "test1. test2. test3"
	str2 := LdapArrayToStr([]string{"test1", "test2", "test3"}, ".")
	if str2 != assertStr2 {
		t.Errorf("Expected %q got %q", assertStr2, str2)
	}

	assertStr3 := "test1 test2 test3"
	str3 := LdapArrayToStr([]string{"test1", "test2", "test3"}, "")
	if str3 != assertStr3 {
		t.Errorf("Expected %q got %q", assertStr3, str3)
	}
}
