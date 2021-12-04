package common

import "testing"

func TestGetBestParsableLocale(t *testing.T) {

	got, err := GetBestParsableLocale()
	if err != nil {
		t.Fatal(err)
	}
	t.Errorf("got locale: %s", got)
}
