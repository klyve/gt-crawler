package main

import "testing"

func TestGetFileName(t *testing.T) {
	test := "https://www.xcontest.org/track.php?t=1542098368.16_5bea8dc027be5.igc"
	expected := "xcontest.org_track.php?t=1542098368.16_5bea8dc027be5.igc"


	actual := GetFileName(test)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestUpload_UploadLinks(t *testing.T) {
	upload := &Upload{
		Auth: mockAuth{},
	}

	testURLs := []string{"https://www.xcontest.org/track.php?t=1542098368.16_5bea8dc027be5.igc"}

	success := make(chan bool, 1)

	upload.UploadLinks(testURLs, success, "", "")

	<- success
}

type mockAuth struct {}

func (mockAuth) GetToken(cPath string, uid string) (token string, err error) {
	token = "test"
	return
}
