package ouidata

import (
	"testing"
)

func TestEmbeddedDBInitialization(t *testing.T) {
	_, err := NewDB()
	if err != nil {
		t.Fatalf("can't load embedded database: %s", err)
	}
}
