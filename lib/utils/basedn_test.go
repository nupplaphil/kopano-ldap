package utils

import "testing"

func TestGetBaseDN(t *testing.T) {
	assertBaseDn1 := "dc=example,dc=org"
	baseDn1 := GetBaseDN("example.org")
	if baseDn1 != assertBaseDn1 {
		t.Errorf("Expected %q got %q", assertBaseDn1, baseDn1)
	}

	assertBaseDn2 := "<nil>"
	baseDn2 := GetBaseDN("")
	if baseDn2 != assertBaseDn2 {
		t.Errorf("Expected %q got %q", assertBaseDn2, baseDn2)
	}

	assertBaseDn3 := "dc=1,dc=2"
	baseDn3 := GetBaseDN("1.2.")
	if baseDn3 != assertBaseDn3 {
		t.Errorf("Expected %q got %q", assertBaseDn3, baseDn3)
	}
}
