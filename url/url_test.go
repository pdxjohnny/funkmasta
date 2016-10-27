package url

import (
	"fmt"
	"net/url"
	"testing"
)

const (
	expectedTestURLEnvArray = "[name=Ava friend=Jess Sarah Zoe]"
)

func TestURLEnvArray(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	v.Add("friend", "Zoe")

	r := fmt.Sprintf("%v", EnvArray(v))
	if r != expectedTestURLEnvArray {
		t.Fatalf("Expected: %v, got: %v", expectedTestURLEnvArray, r)
	}
}
