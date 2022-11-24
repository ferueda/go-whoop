package whoop

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.baseURL.String(), baseURL; got != want {
		t.Errorf("NewClient baseURL is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.http == c2.http {
		t.Error("NewClient returned same http client, but they should be different")
	}
}
