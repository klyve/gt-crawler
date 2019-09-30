package main

import (
	"testing"
)

func TestGetFileName(t *testing.T) {
	test := "https://www.xcontest.org/track.php?t=1542098368.16_5bea8dc027be5.igc"
	expected := "xcontest.org_track.php?t=1542098368.16_5bea8dc027be5.igc"

	actual := GetFileName(test)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestGetDomain(t *testing.T) {
	test := "https://www.xcontest.org/track.php?t=1542098368.16_5bea8dc027be5.igc"
	expected := "https://www.xcontest.org"

	actual := GetDomain(test)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

// Unused code
//type mockAuth struct{}

//func (mockAuth) GetToken(cPath string, uid string) (token string, err error) {
//token = "test"
//return
//}
