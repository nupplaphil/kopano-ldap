package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

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
