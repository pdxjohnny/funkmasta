package url

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

const (
	expectedTestURLEnvArray1 = "name=Ava"
	expectedTestURLEnvArray2 = "friend=Jess Sarah Zoe"
)

func TestURLEnvArray(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	v.Add("friend", "Zoe")

	r := fmt.Sprintf("%v", EnvArray(v))
	if !strings.Contains(r, expectedTestURLEnvArray1) {
		t.Fatalf("Could not find %q in %v", expectedTestURLEnvArray1, r)
	} else if !strings.Contains(r, expectedTestURLEnvArray2) {
		t.Fatalf("Could not find %q in %v", expectedTestURLEnvArray2, r)
	}
}
